package xmlrpc

import (
	"bytes"
	"encoding/xml"
	"testing"
)

func TestIntegerParam(t *testing.T) {
	parseRequest(
		[]string{"<int>482384</int>"},
		[]Value{{Type: ValueInt, Number: 482384}},
		t)

	parseRequest(
		[]string{"<i4>-958372</i4>"},
		[]Value{{Type: ValueInt, Number: -958372}},
		t)

	parseRequest(
		[]string{"<i8>38482148485</i8>"},
		[]Value{{Type: ValueInt, Number: 38482148485}},
		t)

	parseRequest(
		[]string{
			"<int>384424</int>",
			"<i4>73849</i4>",
			"<i8>-938284932367</i8>"},
		[]Value{
			{Type: ValueInt, Number: 384424},
			{Type: ValueInt, Number: 73849},
			{Type: ValueInt, Number: -938284932367}},
		t)
}

func TestArrayParam(t *testing.T) {
	parseRequest(
		[]string{`<array><data>
              <value><int>59392</int></value>
              <value><i4>-49528</i4></value>
              <value><i8>1294959993</i8></value>
              </data></array>`},
		[]Value{{Type: ValueArray, Array: []Value{
			{Type: ValueInt, Number: 59392},
			{Type: ValueInt, Number: -49528},
			{Type: ValueInt, Number: 1294959993}}}},
		t)
}

func TestMixedArrayParam(t *testing.T) {
	parseRequest(
		[]string{`<array><data>
              <value><int>485838</int></value>
              <value><i4>58388</i4></value>
              <value><i8>-4829485744</i8></value>
              <value><string>Hello World &amp; You!</string></value>
              <value>:) :D !:&lt; :&gt;</value>
              </data></array>`},
		[]Value{{Type: ValueArray, Array: []Value{
			{Type: ValueInt, Number: 485838},
			{Type: ValueInt, Number: 58388},
			{Type: ValueInt, Number: -4829485744},
			{Type: ValueString, String: "Hello World & You!"},
			{Type: ValueString, String: ":) :D !:< :>"}}}},
		t)
}

func TestCreateRequest(t *testing.T) {
	expected := xml.Header +
		"<methodCall><methodName>calling</methodName></methodCall>"

	createCompareRequest("calling", nil, expected, t)
}

func TestCreateRequestIntParam(t *testing.T) {
	expected := xml.Header +
		"<methodCall><methodName>hello</methodName>" +
		"<params><param><value><int>100</int></value></param></params></methodCall>"

	createCompareRequest("hello", []Value{{Type: ValueInt, Number: 100}}, expected, t)
}

func createCompareRequest(methodName string, values []Value, expected string, t *testing.T) {
	cl := Client{}

	actual := string(cl.CreateRequest(methodName, values))

	if actual != expected {
		t.Fatalf("Expected document: %s\ngot: %s\n", expected, actual)
	}
}

func parseRequest(params []string, expecteds []Value, t *testing.T) {
	prefix := "\n<?xml version=\"1.0\" encoding=\"UTF-8\"?><methodResponse><params>"
	suffix := "\n</params></methodResponse>"

	xml := prefix
	for _, param := range params {
		xml += "\n<param><value>" + param + "</value></param>"
	}
	xml += suffix

	t.Logf("XML: %s", xml)

	buf := bytes.NewBufferString(xml)
	cl := Client{}
	actuals := cl.Parse(buf)

	if len(actuals) != len(expecteds) {
		t.Fatalf("Expected count %v values, got %v", len(expecteds), len(actuals))
	}

	for i, expected := range expecteds {
		actual := actuals[i]

		compareValue(&actual, &expected, t)
	}
}

func compareValue(actual *Value, expected *Value, t *testing.T) {
	if actual.Type != expected.Type {
		t.Fatalf("Expected type of %d, got %d\n%#v\n%#v\n", expected.Type, actual.Type, expected, actual)
	}

	switch expected.Type {
	case ValueInt:
		if actual.Number != expected.Number {
			t.Errorf("Expected %v, got %v", expected.Number, actual.Number)
		}
	case ValueBoolean:
		if actual.Boolean != expected.Boolean {
			t.Errorf("Expected %v, got %v", expected.Boolean, actual.Boolean)
		}
	case ValueDouble:
		if actual.Double != expected.Double {
			t.Errorf("Expected %v, got %v", expected.Double, actual.Double)
		}
	case ValueString:
		if actual.String != expected.String {
			t.Errorf("Expected %v, got %v", expected.String, actual.String)
		}
	case ValueDate:
		if actual.Date != expected.Date {
			t.Errorf("Expected %v, got %v", expected.Date, actual.Date)
		}
	case ValueBase64:
		if actual.Base64 != expected.Base64 {
			t.Errorf("Expected %v, got %v", expected.Base64, actual.Base64)
		}
	case ValueArray:
		if len(actual.Array) != len(expected.Array) {
			t.Fatalf("Expected values length %v, got %v", len(actual.Array), len(expected.Array))
		}

		for i, expectedValue := range expected.Array {
			actualValue := actual.Array[i]

			compareValue(&actualValue, &expectedValue, t)
		}

	default:
		t.Fatalf("Unhandled type %v", expected.Type)
	}
}
