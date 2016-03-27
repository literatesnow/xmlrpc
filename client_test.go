package xmlrpc

import (
	"bytes"
	"encoding/xml"
	"testing"
)

func integerData() (xmlDoc []string, values []Value) {
	return []string{
			"<int>2147483647</int>", //max int32
			"<i4>-2147483648</i4>",  //min int32
			"<i4>0</i4>",
			"<int>-1</int>",
			"<int>9223372036854775807</int>", //overflow
			"<i4>-9223372036854775808</i4>",  //-overflow
			"<i4></i4>"},                     //invalid
		[]Value{
			NewValueInt(2147483647),
			NewValueInt(-2147483648),
			NewValueInt(0),
			NewValueInt(-1),
			NewValueInt(2147483647),
			NewValueInt(-2147483648),
			NewValueInt(0)}
}

func integerValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<int>2147483647</int>",  //max int32
			"<int>-2147483648</int>", //min int32
			"<int>0</int>",
			"<int>-1</int>"},
		[]Value{
			NewValueInt(2147483647),
			NewValueInt(-2147483648),
			NewValueInt(0),
			NewValueInt(-1)}
}

func booleanData() (xmlDoc []string, values []Value) {
	return []string{
			"<boolean>true</boolean>",
			"<boolean>TRUE</boolean>",
			"<boolean>tRuE</boolean>",
			"<boolean>false</boolean>",
			"<boolean>FALSE</boolean>",
			"<boolean>fAlSe</boolean>",
			"<boolean>1</boolean>",
			"<boolean>0</boolean>",
			"<boolean>yes</boolean>",
			"<boolean>no</boolean>",
			"<boolean>bananas</boolean>",
			"<boolean></boolean>"}, //invalid
		[]Value{
			NewValueBoolean(true),
			NewValueBoolean(true),
			NewValueBoolean(false),
			NewValueBoolean(false),
			NewValueBoolean(false),
			NewValueBoolean(false),
			NewValueBoolean(true),
			NewValueBoolean(false),
			NewValueBoolean(false),
			NewValueBoolean(false),
			NewValueBoolean(false),
			NewValueBoolean(false)}
}

func booleanValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<boolean>true</boolean>",
			"<boolean>false</boolean>"},
		[]Value{
			NewValueBoolean(true),
			NewValueBoolean(false)}
}

func stringData() (xmlDoc []string, values []Value) {
	return []string{
			"<string>This is a string.</string>",
			"<string>This is a string with &amp; &lt; &gt; &quot; &apos; characters.</string>",
			"<string>Unicode: &#0034; &#x0022;.</string>",
			"<string>New line: \n</string>",
			"<string>yes</string>",
			"<string>no</string>",
			"<string>bananas</string>",
			"<string></string>"},
		[]Value{
			NewValueString("This is a string."),
			NewValueString("This is a string with & < > \" ' characters."),
			NewValueString("Unicode: \" \"."),
			NewValueString("New line: \n"),
			NewValueString("yes"),
			NewValueString("no"),
			NewValueString("bananas"),
			NewValueString("")}
}

func stringValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<string>This is a string.</string>",
			"<string>This is a string with &amp; &lt; &gt; &#34; &#39; characters.</string>",
			"<string>Unicode: &#34; &#34;.</string>",
			"<string>New line: \n</string>",
			"<string>yes</string>",
			"<string>no</string>",
			"<string>bananas</string>",
			"<string></string>"},
		[]Value{
			NewValueString("This is a string."),
			NewValueString("This is a string with & < > \" ' characters."),
			NewValueString("Unicode: \" \"."),
			NewValueString("New line: \n"),
			NewValueString("yes"),
			NewValueString("no"),
			NewValueString("bananas"),
			NewValueString("")}
}

func doubleData() (xmlDoc []string, values []Value) {
	return []string{
			"<double>0.123456789012345678901234567890</double>",
			"<double>-0.123456789012345678901234567890</double>",
			"<double>0</double>",
			"<double>0.000000</double>",
			"<double>-1</double>",
			"<double>9223372036854775807</double>",
			"<double>-9223372036854775808</double>",
			"<double></double>"}, //invalid
		[]Value{
			NewValueDouble(0.123456789012345678901234567890),
			NewValueDouble(-0.123456789012345678901234567890),
			NewValueDouble(0),
			NewValueDouble(0.0),
			NewValueDouble(-1),
			NewValueDouble(9223372036854775807),
			NewValueDouble(-9223372036854775808),
			NewValueDouble(0)}
}

func doubleValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<double>0.12345678901234568</double>",
			"<double>-0.12345678901234568</double>",
			"<double>0</double>",
			"<double>-1</double>",
			"<double>9223372036854776000</double>",
			"<double>-9223372036854776000</double>"},
		[]Value{
			NewValueDouble(0.123456789012345678901234567890),
			NewValueDouble(-0.123456789012345678901234567890),
			NewValueDouble(0),
			NewValueDouble(-1),
			NewValueDouble(9223372036854776000),
			NewValueDouble(-9223372036854776000)}
}

func dateTimeData() (xmlDoc []string, values []Value) {
	return []string{
			"<dateTime.iso8601>TODO</dateTime.iso8601>",
			"<dateTime>TODO</dateTime>",
			"<dateTime.iso8601></dateTime.iso8601>"},
		[]Value{
			NewValueDateTime("TODO"),
			NewValueDateTime("TODO"),
			NewValueDateTime("")}
}

func dateTimeValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<dateTime.iso8601>TODO</dateTime.iso8601>"},
		[]Value{
			NewValueDateTime("TODO")}
}

func base64Data() (xmlDoc []string, values []Value) {
	return []string{
			"<base64>TODO</base64>",
			"<base64></base64>"},
		[]Value{
			NewValueBase64("TODO"),
			NewValueBase64("")}
}

func base64ValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<base64>TODO</base64>"},
		[]Value{
			NewValueBase64("TODO")}
}

func arrayData() (xmlDoc []string, values []Value) {
	return []string{`<array><data>
              <value><int>59392</int></value>
              <value><i4>-49528</i4></value>
              <value><array><data><value><string>wheee</string></value></data></array></value>
              <value><i8>1294959993</i8></value>
              </data></array>`},
		[]Value{{Array: []Value{
			NewValueInt(59392),
			NewValueInt(-49528),
			NewValueArray([]Value{NewValueString("wheee")}),
			NewValueLong(1294959993)}}}
}

func arrayValidData() (xmlDoc []string, values []Value) {
	return []string{"<array><data>" +
			"<value><int>59392</int></value>" +
			"<value><int>-49528</int></value>" +
			"<value><array><data><value><string>wheee</string></value></data></array></value>" +
			"<value><ex:i8>1294959993</ex:i8></value>" +
			"</data></array>"},
		[]Value{{Array: []Value{
			NewValueInt(59392),
			NewValueInt(-49528),
			NewValueArray([]Value{NewValueString("wheee")}),
			NewValueLong(1294959993)}}}
}

func nilData() (xmlDoc []string, values []Value) {
	return []string{
			"<nil/>",
			"<ex:nil/>",
			"<nil></nil>",
			"<ex:nil></ex:nil>"},
		[]Value{
			NewValueNil(),
			NewValueNil(),
			NewValueNil(),
			NewValueNil()}
}

func nilValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<ex:nil></ex:nil>"},
		[]Value{
			NewValueNil()}
}

func byteData() (xmlDoc []string, values []Value) {
	return []string{
			"<byte>255</byte>", //max byte
			"<i1>0</i1>",       //min byte
			"<ex:i1>0</ex:i1>",
			"<byte>-1</byte>",
			"<i1>2147483647</i1>",        //overflow
			"<ex:i1>-2147483648</ex:i1>", //-overflow
			"<byte></byte>"},             //invalid
		[]Value{
			NewValueByte(255),
			NewValueByte(0),
			NewValueByte(0),
			NewValueByte(255),
			NewValueByte(255),
			NewValueByte(0),
			NewValueByte(0)}
}

func byteValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<ex:i1>255</ex:i1>", //max byte
			"<ex:i1>0</ex:i1>"},  //min byte
		[]Value{
			NewValueByte(255),
			NewValueByte(0)}
}

func floatData() (xmlDoc []string, values []Value) {
	return []string{
			"<float>0.123456789012345678901234567890</float>",
			"<ex:float>-0.123456789012345678901234567890</ex:float>",
			"<float>0</float>",
			"<ex:float>0.000000</ex:float>",
			"<float>-1</float>",
			"<ex:float>9223372036854775807</ex:float>",
			"<float>-9223372036854775808</float>",
			"<ex:float></ex:float>"}, //invalid
		[]Value{
			NewValueFloat(0.123456789012345678901234567890),
			NewValueFloat(-0.123456789012345678901234567890),
			NewValueFloat(0),
			NewValueFloat(0.0),
			NewValueFloat(-1),
			NewValueFloat(9223372036854775807),
			NewValueFloat(-9223372036854775808),
			NewValueFloat(0)}
}

func floatValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<ex:float>0.12345679</ex:float>",
			"<ex:float>-0.12345679</ex:float>",
			"<ex:float>0</ex:float>",
			"<ex:float>-1</ex:float>"},
		[]Value{
			NewValueFloat(0.12345679),
			NewValueFloat(-0.12345679),
			NewValueFloat(0),
			NewValueFloat(-1)}
}

func longData() (xmlDoc []string, values []Value) {
	return []string{
			"<i8>9223372036854775807</i8>",        //max int64
			"<ex:i8>-9223372036854775808</ex:i8>", //min int64
			"<long>0</long>",
			"<i8>-1</i8>",
			"<ex:i8>92233720368547758079223372036854775807</ex:i8>", //overflow
			"<long>-92233720368547758089223372036854775808</long>",  //-overflow
			"<i8></i8>"}, //invalid
		[]Value{
			NewValueLong(9223372036854775807),
			NewValueLong(-9223372036854775808),
			NewValueLong(0),
			NewValueLong(-1),
			NewValueLong(9223372036854775807),
			NewValueLong(-9223372036854775808),
			NewValueLong(0)}
}

func longValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<ex:i8>9223372036854775807</ex:i8>",  //max int64
			"<ex:i8>-9223372036854775808</ex:i8>", //min int64
			"<ex:i8>0</ex:i8>",
			"<ex:i8>-1</ex:i8>"},
		[]Value{
			NewValueLong(9223372036854775807),
			NewValueLong(-9223372036854775808),
			NewValueLong(0),
			NewValueLong(-1)}
}

func shortData() (xmlDoc []string, values []Value) {
	return []string{
			"<short>32767</short>", //max int16
			"<i2>-32768</i2>",      //min int16
			"<ex:i2>0</ex:i2>",
			"<short>-1</short>",
			"<i2>2147483647</i2>",        //overflow
			"<ex:i2>-2147483648</ex:i2>", //-overflow
			"<short></short>"},           //invalid
		[]Value{
			NewValueShort(32767),
			NewValueShort(-32768),
			NewValueShort(0),
			NewValueShort(-1),
			NewValueShort(32767),
			NewValueShort(-32768),
			NewValueShort(0)}
}

func shortValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<ex:i2>32767</ex:i2>",  //max int16
			"<ex:i2>-32768</ex:i2>", //min int16
			"<ex:i2>0</ex:i2>",
			"<ex:i2>-1</ex:i2>"},
		[]Value{
			NewValueShort(32767),
			NewValueShort(-32768),
			NewValueShort(0),
			NewValueShort(-1)}
}

func mixedArrayData() (xmlDoc []string, values []Value) {
	return []string{`<array><data>
              <value><int>485838</int></value>
              <value><i4>58388</i4></value>
              <value><i8>-4829485744</i8></value>
              <value><string>Hello World &amp; You!</string></value>
              <value>:) :D !:&lt; :&gt;</value>
              </data></array>`},
		[]Value{{Array: []Value{
			NewValueInt(485838),
			NewValueInt(58388),
			NewValueLong(-4829485744),
			NewValueString("Hello World & You!"),
			NewValueString(":) :D !:< :>")}}}
}

func TestIntegerParam(t *testing.T) {
	xmlDoc, values := integerData()
	parseRequest(xmlDoc, values, t)
}

func TestBooleanParam(t *testing.T) {
	xmlDoc, values := booleanData()
	parseRequest(xmlDoc, values, t)
}

func TestStringParam(t *testing.T) {
	xmlDoc, values := stringData()
	parseRequest(xmlDoc, values, t)
}

func TestDoubleParam(t *testing.T) {
	xmlDoc, values := doubleData()
	parseRequest(xmlDoc, values, t)
}

func TestDateTimeParam(t *testing.T) {
	xmlDoc, values := dateTimeData()
	parseRequest(xmlDoc, values, t)
}

func TestBase64Param(t *testing.T) {
	xmlDoc, values := base64Data()
	parseRequest(xmlDoc, values, t)
}

func TestArrayParam(t *testing.T) {
	xmlDoc, values := arrayData()
	parseRequest(xmlDoc, values, t)
}

func TestNilParam(t *testing.T) {
	xmlDoc, values := nilData()
	parseRequest(xmlDoc, values, t)
}

func TestByteParam(t *testing.T) {
	xmlDoc, values := byteData()
	parseRequest(xmlDoc, values, t)
}

func TestFloatParam(t *testing.T) {
	xmlDoc, values := floatData()
	parseRequest(xmlDoc, values, t)
}

func TestLongParam(t *testing.T) {
	xmlDoc, values := longData()
	parseRequest(xmlDoc, values, t)
}

func TestShortParam(t *testing.T) {
	xmlDoc, values := shortData()
	parseRequest(xmlDoc, values, t)
}

func TestMixedArrayParam(t *testing.T) {
	xmlDoc, values := mixedArrayData()
	parseRequest(xmlDoc, values, t)
}

func formatParamValues(param []string) (paramsXml string) {
	paramsXml = ""
	for _, p := range param {
		paramsXml += "<param><value>" + p + "</value></param>"
	}
	return paramsXml
}

func TestCreateRequest(t *testing.T) {
	expected := xml.Header +
		"<methodCall><methodName>Calling</methodName></methodCall>"

	createCompareRequest("Calling", nil, expected, t)
}

func TestCreateRequestIntegerParam(t *testing.T) {
	xmlValues, values := integerValidData()

	expected := xml.Header +
		"<methodCall><methodName>Integer Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("Integer Test", values, expected, t)
}

func TestCreateRequestBooleanParam(t *testing.T) {
	xmlValues, values := booleanValidData()

	expected := xml.Header +
		"<methodCall><methodName>Boolean Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("Boolean Test", values, expected, t)
}

func TestCreateRequestStringParam(t *testing.T) {
	xmlValues, values := stringValidData()

	expected := xml.Header +
		"<methodCall><methodName>String Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("String Test", values, expected, t)
}

func TestCreateRequestDoubleParam(t *testing.T) {
	xmlValues, values := doubleValidData()

	expected := xml.Header +
		"<methodCall><methodName>Double Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("Double Test", values, expected, t)
}

func TestCreateRequestDateTimeParam(t *testing.T) {
	xmlValues, values := dateTimeValidData()

	expected := xml.Header +
		"<methodCall><methodName>Date Time Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("Date Time Test", values, expected, t)
}

func TestCreateRequestBase64Param(t *testing.T) {
	xmlValues, values := dateTimeValidData()

	expected := xml.Header +
		"<methodCall><methodName>Base64 Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("Base64 Test", values, expected, t)
}

func TestCreateRequestArrayParam(t *testing.T) {
	xmlValues, values := arrayValidData()

	expected := xml.Header +
		"<methodCall><methodName>Array Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("Array Test", values, expected, t)
}

func TestCreateRequestNilParam(t *testing.T) {
	xmlValues, values := nilValidData()

	expected := xml.Header +
		"<methodCall><methodName>Nil Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("Nil Test", values, expected, t)
}

func TestCreateRequestByteParam(t *testing.T) {
	xmlValues, values := byteValidData()

	expected := xml.Header +
		"<methodCall><methodName>Byte Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("Byte Test", values, expected, t)
}

func TestCreateRequestFloatParam(t *testing.T) {
	xmlValues, values := floatValidData()

	expected := xml.Header +
		"<methodCall><methodName>Float Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("Float Test", values, expected, t)
}

func TestCreateRequestLongParam(t *testing.T) {
	xmlValues, values := longValidData()

	expected := xml.Header +
		"<methodCall><methodName>Long Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("Long Test", values, expected, t)
}

func TestCreateRequestShortParam(t *testing.T) {
	xmlValues, values := shortValidData()

	expected := xml.Header +
		"<methodCall><methodName>Short Test</methodName><params>" +
		formatParamValues(xmlValues) +
		"</params></methodCall>"

	createCompareRequest("Short Test", values, expected, t)
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

	t.Logf("Expected (XML): %s", xml)

	buf := bytes.NewBufferString(xml)
	cl := Client{}
	actuals := cl.ParseResponse(buf)

	if len(actuals) != len(expecteds) {
		t.Fatalf("Expected count %v values, got %v", len(expecteds), len(actuals))
	}

	t.Logf("===========================\n")
	t.Logf("Expected: %v\n", expecteds)
	t.Logf("---------------------------\n")
	t.Logf("Actual: %v\n", actuals)
	t.Logf("===========================\n")

	for i, expected := range expecteds {
		actual := actuals[i]

		compareValue(&expected, &actual, t)
	}
}

func compareValue(expected *Value, actual *Value, t *testing.T) {
	if expected.Int != nil {
		if actual.Int == nil {
			t.Errorf("Expected %#v, got nil", *expected.Int)
		} else if *expected.Int != *actual.Int {
			t.Errorf("Expected %#v, got %#v", *expected.Int, *actual.Int)
		}
		return
	}

	if expected.Boolean != nil {
		if actual.Boolean == nil {
			t.Errorf("Expected %#v, got nil", *expected.Boolean)
		} else if *expected.Boolean != *actual.Boolean {
			t.Errorf("Expected %#v, got %#v", *expected.Boolean, *actual.Boolean)
		}
		return
	}

	if expected.String != nil {
		if actual.String == nil {
			t.Errorf("Expected %#v, got nil", *expected.String)
		} else if *expected.String != *actual.String {
			t.Errorf("Expected %#v, got %#v", *expected.String, *actual.String)
		}
		return
	}

	if expected.Double != nil {
		if actual.Double == nil {
			t.Errorf("Expected %#v, got nil", *expected.Double)
		} else if *expected.Double != *actual.Double {
			t.Errorf("Expected %#v, got %#v", *expected.Double, *actual.Double)
		}
		return
	}

	if expected.DateTime != nil {
		if actual.DateTime == nil {
			t.Errorf("Expected %#v, got nil", *expected.DateTime)
		} else if *expected.DateTime != *actual.DateTime {
			t.Errorf("Expected %#v, got %#v", *expected.DateTime, *actual.DateTime)
		}
		return
	}

	if expected.Base64 != nil {
		if actual.Base64 == nil {
			t.Errorf("Expected %#v, got nil", *expected.Base64)
		} else if *expected.Base64 != *actual.Base64 {
			t.Errorf("Expected %#v, got %#v", *expected.Base64, *actual.Base64)
		}
		return
	}

	if expected.Array != nil {
		if len(expected.Array) != len(actual.Array) {
			t.Fatalf("Expected values length %d, got %d", len(expected.Array), len(actual.Array))
		}

		for i, expectedValue := range expected.Array {
			actualValue := actual.Array[i]
			compareValue(&expectedValue, &actualValue, t)
		}

		return
	}

	//Extensions

	if expected.Nil != nil {
		if actual.Nil == nil {
			t.Errorf("Expected %#v, got nil", *expected.Nil)
		} else if *expected.Nil != *actual.Nil {
			t.Errorf("Expected %#v, got %#v", *expected.Nil, *actual.Nil)
		}
		return
	}

	if expected.Byte != nil {
		if actual.Byte == nil {
			t.Errorf("Expected %#v, got nil", *expected.Byte)
		} else if *expected.Byte != *actual.Byte {
			t.Errorf("Expected %#v, got %#v", *expected.Byte, *actual.Byte)
		}
		return
	}

	if expected.Float != nil {
		if actual.Float == nil {
			t.Errorf("Expected %#v, got nil", *expected.Float)
		} else if *expected.Float != *actual.Float {
			t.Errorf("Expected %#v, got %#v", *expected.Float, *actual.Float)
		}
		return
	}

	if expected.Long != nil {
		if actual.Long == nil {
			t.Errorf("Expected %#v, got nil", *expected.Long)
		} else if *expected.Long != *actual.Long {
			t.Errorf("Expected %#v, got %#v", *expected.Long, *actual.Long)
		}
		return
	}

	if expected.Short != nil {
		if actual.Short == nil {
			t.Errorf("Expected %#v, got nil", *expected.Short)
		} else if *expected.Short != *actual.Short {
			t.Errorf("Expected %#v, got %#v", *expected.Short, *actual.Short)
		}
		return
	}

	t.Fatalf("Nothing matched while comparing\n")
}
