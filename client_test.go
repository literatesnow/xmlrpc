package xmlrpc

import (
	"bytes"
	"encoding/xml"
	"strconv"
	"strings"
	"testing"
	"time"
)

var NZ = time.FixedZone("NZ", 12*3600)
var ZZ time.Time

// Both

func createCompareRequest(methodName string, params []Value, expected string, t *testing.T) {
	actual := string(CreateRequest(methodName, params))

	if actual != expected {
		t.Fatalf("Expected document: %s\ngot: %s\n", expected, actual)
	}
}

func runParseXmlResponse(items []string, expecteds []Value, isParam bool, t *testing.T) {
	if len(items) != len(expecteds) {
		t.Fatalf("Expected items count %v, got %v", len(items), len(expecteds))
	}

	prefix := "\n<?xml version=\"1.0\" encoding=\"UTF-8\"?> <methodResponse> <params> "
	suffix := "\n </params> </methodResponse>"

	for i, item := range items {
		expected := expecteds[i]

		xml := prefix
		if isParam {
			xml += "<param><value>" + item + "</value></param>"
		} else {
			xml += item
		}
		xml += suffix

		t.Logf("Expected (XML): %s", xml)

		buf := bytes.NewBufferString(xml)
		actual, err := ParseResponse(buf)

		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		t.Logf("===========================\n")
		t.Logf("Expected: %s\n", printValue(&expected))
		t.Logf("---------------------------\n")
		t.Logf("Actual: %s\n", printValue(actual))
		t.Logf("===========================\n")

		compareValue(&expected, actual, t)
	}
}

func runParseJsonRequest(params []string, expecteds []Value, t *testing.T) {
	prefix := "\n" + `{"methodName":"helloMethod","params":[`
	suffix := "\n" + `]}`

	json := prefix
	json += strings.Join(params, ",")
	json += suffix

	t.Logf("Expected (JSON): %s", json)

	buf := bytes.NewReader([]byte(json))

	methodName, actuals, err := ParseJsonRequest(buf)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}

	if methodName != "helloMethod" {
		t.Fatalf("Expected helloMethod, got %s\n", methodName)
	}

	if len(actuals) != len(expecteds) {
		t.Fatalf("Expected count %v values, got %v", len(expecteds), len(actuals))
	}

	t.Logf("===========================\n")
	t.Logf("Expected: %s\n", printValueArray(expecteds))
	t.Logf("---------------------------\n")
	t.Logf("Actual: %s\n", printValueArray(actuals))
	t.Logf("===========================\n")

	for i, expected := range expecteds {
		actual := actuals[i]

		compareValue(&expected, &actual, t)
	}
}

func compareMember(expected *Member, actual *Member, t *testing.T) {
	if expected == nil {
		t.Fatalf("Expected was nil")
	}
	if actual == nil {
		t.Fatalf("Actual was nil")
	}

	if expected.Name != actual.Name {
		t.Fatalf("Expecting name %s, got %s", expected.Name, actual.Name)
	}

	compareValue(&expected.Value, &actual.Value, t)
}

func compareValue(expected *Value, actual *Value, t *testing.T) {
	if expected == nil {
		t.Fatalf("Expected was nil")
	}
	if actual == nil {
		t.Fatalf("Actual was nil")
	}

	matched := false

	if expected.Int != nil {
		if actual.Int == nil {
			t.Errorf("Expected %#v, got nil", *expected.Int)
		} else if *expected.Int != *actual.Int {
			t.Errorf("Expected %#v, got %#v", *expected.Int, *actual.Int)
		}
		matched = true
	}

	if expected.Boolean != nil {
		if actual.Boolean == nil {
			t.Errorf("Expected %#v, got nil", *expected.Boolean)
		} else if *expected.Boolean != *actual.Boolean {
			t.Errorf("Expected %#v, got %#v", *expected.Boolean, *actual.Boolean)
		}
		matched = true
	}

	if expected.String != nil {
		if actual.String == nil {
			t.Errorf("Expected %#v, got nil", *expected.String)
		} else if *expected.String != *actual.String {
			t.Errorf("Expected %#v, got %#v", *expected.String, *actual.String)
		}
		matched = true
	}

	if expected.Double != nil {
		if actual.Double == nil {
			t.Errorf("Expected %#v, got nil", *expected.Double)
		} else if *expected.Double != *actual.Double {
			t.Errorf("Expected %#v, got %#v", *expected.Double, *actual.Double)
		}
		matched = true
	}

	if expected.DateTime != nil {
		if actual.DateTime == nil {
			t.Errorf("Expected %#v, got nil", *expected.DateTime)
		} else if !(*expected.DateTime).Equal(*actual.DateTime) {
			t.Errorf("Expected %#v, got %#v", *expected.DateTime, *actual.DateTime)
		}
		matched = true
	}

	if expected.Base64 != nil {
		if actual.Base64 == nil {
			t.Errorf("Expected %#v, got nil", *expected.Base64)
		} else if *expected.Base64 != *actual.Base64 {
			t.Errorf("Expected %#v, got %#v", *expected.Base64, *actual.Base64)
		}
		matched = true
	}

	if expected.Array != nil {
		if len(expected.Array) != len(actual.Array) {
			t.Fatalf("Expected values length %d, got %d", len(expected.Array), len(actual.Array))
		}

		for i, expectedValue := range expected.Array {
			actualValue := actual.Array[i]
			compareValue(&expectedValue, &actualValue, t)
		}

		matched = true
	}

	if expected.Struct != nil {
		if len(expected.Struct) != len(actual.Struct) {
			t.Fatalf("Expected struct length %d, got %d", len(expected.Struct), len(actual.Struct))
		}

		for i, expectedMember := range expected.Struct {
			actualMember := actual.Struct[i]
			compareMember(&expectedMember, &actualMember, t)
		}

		matched = true
	}

	//Extensions

	if expected.Nil != nil {
		if actual.Nil == nil {
			t.Errorf("Expected %#v, got nil", *expected.Nil)
		} else if *expected.Nil != *actual.Nil {
			t.Errorf("Expected %#v, got %#v", *expected.Nil, *actual.Nil)
		}
		matched = true
	}

	if expected.Byte != nil {
		if actual.Byte == nil {
			t.Errorf("Expected %#v, got nil", *expected.Byte)
		} else if *expected.Byte != *actual.Byte {
			t.Errorf("Expected %#v, got %#v", *expected.Byte, *actual.Byte)
		}
		matched = true
	}

	if expected.Float != nil {
		if actual.Float == nil {
			t.Errorf("Expected %#v, got nil", *expected.Float)
		} else if *expected.Float != *actual.Float {
			t.Errorf("Expected %#v, got %#v", *expected.Float, *actual.Float)
		}
		matched = true
	}

	if expected.Long != nil {
		if actual.Long == nil {
			t.Errorf("Expected %#v, got nil", *expected.Long)
		} else if *expected.Long != *actual.Long {
			t.Errorf("Expected %#v, got %#v", *expected.Long, *actual.Long)
		}
		matched = true
	}

	if expected.Short != nil {
		if actual.Short == nil {
			t.Errorf("Expected %#v, got nil", *expected.Short)
		} else if *expected.Short != *actual.Short {
			t.Errorf("Expected %#v, got %#v", *expected.Short, *actual.Short)
		}
		matched = true
	}

	if !matched {
		t.Fatalf("Nothing matched while comparing\n")
	}

	//Check if multiple value types have been set
	expectedScore := scoreValue(expected)
	actualScore := scoreValue(actual)

	if expectedScore != actualScore {
		t.Fatalf("Expecting score of %d, got %d\n", expectedScore, actualScore)
	}
}

func scoreValue(v *Value) (score int) {
	score = 0
	if v.Int != nil {
		score |= 1
	}
	if v.Boolean != nil {
		score |= 2
	}
	if v.String != nil {
		score |= 4
	}
	if v.Double != nil {
		score |= 8
	}
	if v.DateTime != nil {
		score |= 16
	}
	if v.Base64 != nil {
		score |= 32
	}
	if v.Array != nil {
		score |= 64
	}
	if v.Struct != nil {
		score |= 128
	}
	if v.Nil != nil {
		score |= 256
	}
	if v.Byte != nil {
		score |= 512
	}
	if v.Float != nil {
		score |= 1024
	}
	if v.Long != nil {
		score |= 2048
	}
	if v.Short != nil {
		score |= 4096
	}

	return score
}

func printValue(value *Value) (str string) {
	if value == nil {
		return "(nil)"
	}

	return (*value).Print()
}

func printValueArray(values []Value) (str string) {
	if values == nil {
		return "(nil)"
	}

	for i, val := range values {
		str += "\n[" + strconv.Itoa(i) + "]: " + val.Print()
	}

	return str
}

func formatParamValues(param []string) (paramsXml string) {
	paramsXml = ""
	for _, p := range param {
		paramsXml += "<param><value>" + p + "</value></param>"
	}
	return paramsXml
}

// XML

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
			NewInt(2147483647),
			NewInt(-2147483648),
			NewInt(0),
			NewInt(-1),
			NewInt(2147483647),
			NewInt(-2147483648),
			NewInt(0)}
}

func integerValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<int>2147483647</int>",  //max int32
			"<int>-2147483648</int>", //min int32
			"<int>0</int>",
			"<int>-1</int>"},
		[]Value{
			NewInt(2147483647),
			NewInt(-2147483648),
			NewInt(0),
			NewInt(-1)}
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
			NewBoolean(true),
			NewBoolean(true),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(true),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false)}
}

func booleanValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<boolean>1</boolean>",
			"<boolean>0</boolean>"},
		[]Value{
			NewBoolean(true),
			NewBoolean(false)}
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
			NewString("This is a string."),
			NewString("This is a string with & < > \" ' characters."),
			NewString("Unicode: \" \"."),
			NewString("New line: \n"),
			NewString("yes"),
			NewString("no"),
			NewString("bananas"),
			NewString("")}
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
			NewString("This is a string."),
			NewString("This is a string with & < > \" ' characters."),
			NewString("Unicode: \" \"."),
			NewString("New line: \n"),
			NewString("yes"),
			NewString("no"),
			NewString("bananas"),
			NewString("")}
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
			NewDouble(0.123456789012345678901234567890),
			NewDouble(-0.123456789012345678901234567890),
			NewDouble(0),
			NewDouble(0.0),
			NewDouble(-1),
			NewDouble(9223372036854775807),
			NewDouble(-9223372036854775808),
			NewDouble(0)}
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
			NewDouble(0.123456789012345678901234567890),
			NewDouble(-0.123456789012345678901234567890),
			NewDouble(0),
			NewDouble(-1),
			NewDouble(9223372036854776000),
			NewDouble(-9223372036854776000)}
}

func dateTimeData() (xmlDoc []string, values []Value) {
	return []string{
			"<dateTime.iso8601>2016-04-07T21:13:58+1200</dateTime.iso8601>",
			"<dateTime8601>2016-04-07T21:13:58+1200</dateTime8601>",
			"<dateTime>2016-04-07T21:13:58+1200</dateTime>",
			"<dateTime.iso8601></dateTime.iso8601>"},
		[]Value{
			NewDateTime(time.Date(2016, 4, 7, 21, 13, 58, 0, NZ)),
			NewDateTime(time.Date(2016, 4, 7, 21, 13, 58, 0, NZ)),
			NewDateTime(time.Date(2016, 4, 7, 21, 13, 58, 0, NZ)),
			NewDateTime(ZZ)}
}

func dateTimeValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<dateTime.iso8601>2016-04-07T21:13:58+1200</dateTime.iso8601>"},
		[]Value{
			NewDateTime(time.Date(2016, 4, 7, 21, 13, 58, 0, NZ))}
}

func base64Data() (xmlDoc []string, values []Value) {
	return []string{
			"<base64>TODO</base64>",
			"<base64></base64>"},
		[]Value{
			NewBase64("TODO"),
			NewBase64("")}
}

func base64ValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<base64>TODO</base64>"},
		[]Value{
			NewBase64("TODO")}
}

func arrayData() (xmlDoc []string, values []Value) {
	return []string{`<array><data>
              <value><int>59392</int></value>
              <value><i4>-49528</i4></value>
              <value><array><data><value><string>wheee</string></value></data></array></value>
              <value><i8>1294959993</i8></value>
              </data></array>`},
		[]Value{{Array: []Value{
			NewInt(59392),
			NewInt(-49528),
			NewArray([]Value{NewString("wheee")}),
			NewLong(1294959993)}}}
}

func arrayValidData() (xmlDoc []string, values []Value) {
	return []string{"<array><data>" +
			"<value><int>59392</int></value>" +
			"<value><int>-49528</int></value>" +
			"<value><array><data><value><string>wheee</string></value></data></array></value>" +
			"<value><ex:i8>1294959993</ex:i8></value>" +
			"</data></array>"},
		[]Value{{Array: []Value{
			NewInt(59392),
			NewInt(-49528),
			NewArray([]Value{NewString("wheee")}),
			NewLong(1294959993)}}}
}

func structData() (xmlDoc []string, values []Value) {
	return []string{
			`<value>  <struct>
         <member> <name>1st item</name> <value><i4>-100</i4> </value> </member>
         <member> <name>2nd item</name> <value><string>-200</string> </value> </member>
         <member> <name>3rd item</name> <value><i8>-300</i8> </value> </member>
         <member>
           <name>4th Array</name>
           <value>
             <array>
               <data>
                 <value> <string>4th #1</string> </value>
                 <value> <string>4th #2</string> </value>
                 <value>
                   <struct>
                     <member>
                       <name>4th #3 - #1</name>
                       <value> <int>1000</int> </value>
                     </member>
                     <member>
                       <name>4th #3 - #2</name>
                       <value> <int>2000</int> </value>
                     </member>
                   </struct>
                 </value>
               </data>
             </array>
           </value>
         </member>
       </struct> </value> `},
		[]Value{NewStruct([]Member{
			{Name: "1st item", Value: NewInt(-100)},
			{Name: "2nd item", Value: NewString("-200")},
			{Name: "3rd item", Value: NewLong(-300)},
			{Name: "4th Array", Value: NewArray([]Value{
				NewString("4th #1"),
				NewString("4th #2"),
				NewStruct([]Member{
					{Name: "4th #3 - #1", Value: NewInt(1000)},
					{Name: "4th #3 - #2", Value: NewInt(2000)}})})}})}
}

func nilData() (xmlDoc []string, values []Value) {
	return []string{
			"<nil/>",
			"<ex:nil/>",
			"<nil></nil>",
			"<ex:nil></ex:nil>"},
		[]Value{
			NewNil(),
			NewNil(),
			NewNil(),
			NewNil()}
}

func nilValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<ex:nil></ex:nil>"},
		[]Value{
			NewNil()}
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
			NewByte(255),
			NewByte(0),
			NewByte(0),
			NewByte(255),
			NewByte(255),
			NewByte(0),
			NewByte(0)}
}

func byteValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<ex:i1>255</ex:i1>", //max byte
			"<ex:i1>0</ex:i1>"},  //min byte
		[]Value{
			NewByte(255),
			NewByte(0)}
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
			NewFloat(0.123456789012345678901234567890),
			NewFloat(-0.123456789012345678901234567890),
			NewFloat(0),
			NewFloat(0.0),
			NewFloat(-1),
			NewFloat(9223372036854775807),
			NewFloat(-9223372036854775808),
			NewFloat(0)}
}

func floatValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<ex:float>0.12345679</ex:float>",
			"<ex:float>-0.12345679</ex:float>",
			"<ex:float>0</ex:float>",
			"<ex:float>-1</ex:float>"},
		[]Value{
			NewFloat(0.12345679),
			NewFloat(-0.12345679),
			NewFloat(0),
			NewFloat(-1)}
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
			NewLong(9223372036854775807),
			NewLong(-9223372036854775808),
			NewLong(0),
			NewLong(-1),
			NewLong(9223372036854775807),
			NewLong(-9223372036854775808),
			NewLong(0)}
}

func longValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<ex:i8>9223372036854775807</ex:i8>",  //max int64
			"<ex:i8>-9223372036854775808</ex:i8>", //min int64
			"<ex:i8>0</ex:i8>",
			"<ex:i8>-1</ex:i8>"},
		[]Value{
			NewLong(9223372036854775807),
			NewLong(-9223372036854775808),
			NewLong(0),
			NewLong(-1)}
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
			NewShort(32767),
			NewShort(-32768),
			NewShort(0),
			NewShort(-1),
			NewShort(32767),
			NewShort(-32768),
			NewShort(0)}
}

func shortValidData() (xmlDoc []string, values []Value) {
	return []string{
			"<ex:i2>32767</ex:i2>",  //max int16
			"<ex:i2>-32768</ex:i2>", //min int16
			"<ex:i2>0</ex:i2>",
			"<ex:i2>-1</ex:i2>"},
		[]Value{
			NewShort(32767),
			NewShort(-32768),
			NewShort(0),
			NewShort(-1)}
}

func mixedArrayData() (xmlDoc []string, values []Value) {
	return []string{`<array><data>
              <value>  <int>485838</int>   </value>
              <value><i4>58388</i4></value>
              <value><i8>-4829485744</i8></value>
              <value><boolean>1</boolean></value>
              <value><string>Hello World &amp; You!</string></value>
              <value>:) :D !:&lt; :&gt;</value>
              <value><struct>
                 <member>     <name>A mighty fine struct</name>
                 <value>  <string>It sure is</string>    </value>
                 </member>    </struct>    </value>
              </data></array>`},
		[]Value{{Array: []Value{
			NewInt(485838),
			NewInt(58388),
			NewLong(-4829485744),
			NewBoolean(true),
			NewString("Hello World & You!"),
			NewString(":) :D !:< :>"),
			NewStruct([]Member{
				{Name: "A mighty fine struct", Value: NewString("It sure is")}})}}}
}

func faultData() (xmlDoc []string, values []Value) {
	return []string{
			`<fault>
      <value><struct>
      <member><name>faultCode</name>
      <value><i4>-506</i4></value></member>
      <member><name>faultString</name>
      <value><string>Method 'what' not defined</string></value></member>
      </struct></value>
    </fault>`},
		[]Value{NewStruct([]Member{
			{Name: "faultCode", Value: NewInt(-506)},
			{Name: "faultString", Value: NewString("Method 'what' not defined")}})}
}

func TestIntegerParam(t *testing.T) {
	xmlDoc, values := integerData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestBooleanParam(t *testing.T) {
	xmlDoc, values := booleanData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestStringParam(t *testing.T) {
	xmlDoc, values := stringData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestDoubleParam(t *testing.T) {
	xmlDoc, values := doubleData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestDateTimeParam(t *testing.T) {
	xmlDoc, values := dateTimeData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestBase64Param(t *testing.T) {
	xmlDoc, values := base64Data()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestArrayParam(t *testing.T) {
	xmlDoc, values := arrayData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestStructParam(t *testing.T) {
	xmlDoc, values := structData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestNilParam(t *testing.T) {
	xmlDoc, values := nilData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestByteParam(t *testing.T) {
	xmlDoc, values := byteData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestFloatParam(t *testing.T) {
	xmlDoc, values := floatData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestLongParam(t *testing.T) {
	xmlDoc, values := longData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestShortParam(t *testing.T) {
	xmlDoc, values := shortData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestMixedArrayParam(t *testing.T) {
	xmlDoc, values := mixedArrayData()
	runParseXmlResponse(xmlDoc, values, true, t)
}

func TestFault(t *testing.T) {
	xmlDoc, values := faultData()
	runParseXmlResponse(xmlDoc, values, false, t)
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
	xmlValues, values := base64ValidData()

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

// JSON

func integerJsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"int":2147483647}`, //max int32
			`{"i4":-2147483648}`, //min int32
			`{"i4":0}`,
			`{"int":-1}`,
			`{"int":9223372036854775807}`, //overflow
			`{"i4":-9223372036854775808}`, //-overflow
			`{"i4":"invalid"}`},           //invalid
		[]Value{
			NewInt(2147483647),
			NewInt(-2147483648),
			NewInt(0),
			NewInt(-1),
			NewInt(-2147483648),
			NewInt(-2147483648),
			NewInt(0)}
}

func integerJsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"int":2147483647}`,  //max int32
			`{"int":-2147483648}`, //min int32
			`{"int":0}`,
			`{"int":-1}`},
		[]Value{
			NewInt(2147483647),
			NewInt(-2147483648),
			NewInt(0),
			NewInt(-1)}
}

func booleanJsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"boolean":true}`,
			`{"boolean":"TRUE"}`,
			`{"boolean":"tRuE"}`,
			`{"boolean":false}`,
			`{"boolean":"FALSE"}`,
			`{"boolean":"fAlSe"}`,
			`{"boolean":1}`,
			`{"boolean":0}`,
			`{"boolean":"yes"}`,
			`{"boolean":"no"}`,
			`{"boolean":null}`,
			`{"boolean":"invalid"}`}, //invalid
		[]Value{
			NewBoolean(true),
			NewBoolean(true),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false),
			NewBoolean(false)}
}

func booleanJsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"boolean":true}`,
			`{"boolean":false}`},
		[]Value{
			NewBoolean(true),
			NewBoolean(false)}
}

func stringJsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"string":"This is a \"string\"."}`,
			`{"string":"This is a 'string' with &amp; &lt; &gt; &quot; &apos; characters."}`,
			`{"string":"Unicode: \u0027hello\u0027 \u0022there\u0022."}`,
			`{"string":"New line: \n"}`,
			`{"string":"yes"}`,
			`{"string":"no"}`,
			`{"string":"bananas"}`,
			`{"string":null}`,
			`{"string":""}`},
		[]Value{
			NewString("This is a \"string\"."),
			NewString("This is a 'string' with &amp; &lt; &gt; &quot; &apos; characters."),
			NewString("Unicode: 'hello' \"there\"."),
			NewString("New line: \n"),
			NewString("yes"),
			NewString("no"),
			NewString("bananas"),
			NewString(""),
			NewString("")}
}

func stringJsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"string":"This is a string."}`,
			`{"string":"This is a string with &amp; &lt; &gt; &#34; &#39; characters."}`,
			`{"string":"Unicode: &#34; &#34;."}`,
			`{"string":"New line: \n"}`,
			`{"string":"yes"}`,
			`{"string":"no"}`,
			`{"string":"bananas"}`,
			`{"string":""}`},
		[]Value{
			NewString("This is a string."),
			NewString("This is a string with & < > \" ' characters."),
			NewString("Unicode: \" \"."),
			NewString("New line: \n"),
			NewString("yes"),
			NewString("no"),
			NewString("bananas"),
			NewString("")}
}

func doubleJsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"double":0.123456789012345678901234567890}`,
			`{"double":-0.123456789012345678901234567890}`,
			`{"double":0}`,
			`{"double":0.000000}`,
			`{"double":-1}`,
			`{"double":9223372036854775807}`,
			`{"double":-9223372036854775808}`,
			`{"double":"invalid"}`}, //invalid
		[]Value{
			NewDouble(0.123456789012345678901234567890),
			NewDouble(-0.123456789012345678901234567890),
			NewDouble(0),
			NewDouble(0.0),
			NewDouble(-1),
			NewDouble(9223372036854775807),
			NewDouble(-9223372036854775808),
			NewDouble(0)}
}

func doubleJsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"double":0.12345678901234568}`,
			`{"double":-0.12345678901234568}`,
			`{"double":0}`,
			`{"double":-1}`,
			`{"double":9223372036854776000}`,
			`{"double":-9223372036854776000}`},
		[]Value{
			NewDouble(0.123456789012345678901234567890),
			NewDouble(-0.123456789012345678901234567890),
			NewDouble(0),
			NewDouble(-1),
			NewDouble(9223372036854776000),
			NewDouble(-9223372036854776000)}
}

func dateTimeJsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"dateTime.iso8601":"2016-04-07T21:13:58+1200"}`,
			`{"dateTime":"2016-04-07T21:13:58+1200"}`,
			`{"dateTime.iso8601":""}`},
		[]Value{
			NewDateTime(time.Date(2016, 4, 7, 21, 13, 58, 0, NZ)),
			NewDateTime(time.Date(2016, 4, 7, 21, 13, 58, 0, NZ)),
			NewDateTime(ZZ)}
}

func dateTimeJsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"dateTime.iso8601":"2016-04-07T21:13:58+1200"}`},
		[]Value{
			NewDateTime(time.Date(2016, 4, 7, 21, 13, 58, 0, NZ))}
}

func base64JsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"base64":"TODO"}`,
			`{"base64":""}`},
		[]Value{
			NewBase64("TODO"),
			NewBase64("")}
}

func base64JsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"base64":"TODO"}`},
		[]Value{
			NewBase64("TODO")}
}

func arrayJsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"array":[
      {"int":59392},
      {"i4":-49528},
      {"array":[{"string":"wheee"}]},
      {"i8":1294959993}]}`},
		[]Value{{Array: []Value{
			NewInt(59392),
			NewInt(-49528),
			NewArray([]Value{NewString("wheee")}),
			NewLong(1294959993)}}}
}

func arrayJsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"int":59392}`,
			`{"i4":-49528}`,
			`{"array":[{"string":"wheee"}]}`,
			`{"i8":1294959993}`},
		[]Value{{Array: []Value{
			NewInt(59392),
			NewInt(-49528),
			NewArray([]Value{NewString("wheee")}),
			NewLong(1294959993)}}}
}

func nilJsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"nil":null}`,
			`{"ex:nil":null}`},
		[]Value{
			NewNil(),
			NewNil()}
}

func nilJsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"ex:nil":null}`},
		[]Value{
			NewNil()}
}

func byteJsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"byte":255}`, //max byte
			`{"i1":0}`,     //min byte
			`{"ex:i1":0}`,
			`{"byte":-1}`,
			`{"i1":2147483647}`,     //overflow
			`{"ex:i1":-2147483648}`, //-overflow
			`{"byte":"invalid"}`},   //invalid
		[]Value{
			NewByte(255),
			NewByte(0),
			NewByte(0),
			NewByte(255),
			NewByte(255),
			NewByte(0),
			NewByte(0)}
}

func byteJsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"ex:i1":255}`, //max byte
			`{"ex:i1":0}`},  //min byte
		[]Value{
			NewByte(255),
			NewByte(0)}
}

func floatJsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"float":0.123456789012345678901234567890}`,
			`{"ex:float":-0.123456789012345678901234567890}`,
			`{"float":0}`,
			`{"ex:float":0.000000}`,
			`{"float":-1}`,
			`{"ex:float":9223372036854775807}`,
			`{"float":-9223372036854775808}`,
			`{"ex:float":"invalid"}`}, //invalid
		[]Value{
			NewFloat(0.123456789012345678901234567890),
			NewFloat(-0.123456789012345678901234567890),
			NewFloat(0),
			NewFloat(0.0),
			NewFloat(-1),
			NewFloat(9223372036854775807),
			NewFloat(-9223372036854775808),
			NewFloat(0)}
}

func floatJsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"ex:float":0.12345679}`,
			`{"ex:float":-0.12345679}`,
			`{"ex:float":0}`,
			`{"ex:float":-1}`},
		[]Value{
			NewFloat(0.12345679),
			NewFloat(-0.12345679),
			NewFloat(0),
			NewFloat(-1)}
}

func longJsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"i8":9223372036854775807}`,     //max int64
			`{"ex:i8":-9223372036854775808}`, //min int64
			`{"long":0}`,
			`{"i8":-1}`,
			`{"ex:i8":92233720368547758079223372036854775807}`, //overflow
			`{"long":-92233720368547758089223372036854775808}`, //-overflow
			`{"i8":"invalid"}`},                                //invalid
		[]Value{
			NewLong(-9223372036854775808),
			NewLong(-9223372036854775808),
			NewLong(0),
			NewLong(-1),
			NewLong(-9223372036854775808),
			NewLong(-9223372036854775808),
			NewLong(0)}
}

func longJsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"ex:i8":9223372036854775807}`,  //max int64
			`{"ex:i8":-9223372036854775808}`, //min int64
			`{"ex:i8":0}`,
			`{"ex:i8":-1}`},
		[]Value{
			NewLong(9223372036854775807),
			NewLong(-9223372036854775808),
			NewLong(0),
			NewLong(-1)}
}

func shortJsonData() (jsonDoc []string, values []Value) {
	return []string{
			`{"short":32767}`, //max int16
			`{"i2":-32768}`,   //min int16
			`{"ex:i2":0}`,
			`{"short":-1}`,
			`{"i2":2147483647}`,     //overflow
			`{"ex:i2":-2147483648}`, //-overflow
			`{"short":"invalid"}`},  //invalid
		[]Value{
			NewShort(32767),
			NewShort(-32768),
			NewShort(0),
			NewShort(-1),
			NewShort(-1),
			NewShort(0),
			NewShort(0)}
}

func shortJsonValidData() (jsonDoc []string, values []Value) {
	return []string{
			`{"ex:i2":32767}`,  //max int16
			`{"ex:i2":-32768}`, //min int16
			`{"ex:i2":0}`,
			`{"ex:i2":-1}`},
		[]Value{
			NewShort(32767),
			NewShort(-32768),
			NewShort(0),
			NewShort(-1)}
}

func mixedArrayJsonData() (jsonDoc []string, values []Value) {
	return []string{`{"array":[
              {"int":485838},
              {"i4":58388},
              {"i8":-4829485744},
              {"boolean":true},
              {"string":"Hello World & \"You\"!"}]}`},
		[]Value{{Array: []Value{
			NewInt(485838),
			NewInt(58388),
			NewLong(-4829485744),
			NewBoolean(true),
			NewString("Hello World & \"You\"!")}}}
}

func TestIntegerJsonParam(t *testing.T) {
	jsonDoc, values := integerJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestBooleanJsonParam(t *testing.T) {
	jsonDoc, values := booleanJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestStringJsonParam(t *testing.T) {
	jsonDoc, values := stringJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestDoubleJsonParam(t *testing.T) {
	jsonDoc, values := doubleJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestDateTimeJsonParam(t *testing.T) {
	jsonDoc, values := dateTimeJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestBase64JsonParam(t *testing.T) {
	jsonDoc, values := base64JsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestArrayJsonParam(t *testing.T) {
	jsonDoc, values := arrayJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestNilJsonParam(t *testing.T) {
	jsonDoc, values := nilJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestByteJsonParam(t *testing.T) {
	jsonDoc, values := byteJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestFloatJsonParam(t *testing.T) {
	jsonDoc, values := floatJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestLongJsonParam(t *testing.T) {
	jsonDoc, values := longJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestShortJsonParam(t *testing.T) {
	jsonDoc, values := shortJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

func TestMixedArrayJsonParam(t *testing.T) {
	jsonDoc, values := mixedArrayJsonData()
	runParseJsonRequest(jsonDoc, values, t)
}

/*
func TestParseJsonRequest(t *testing.T) {
  json := `{"methodName":"system.client_version"}`
  body := bytes.NewReader([]byte(json))

  methodName, params, err := ParseJsonRequest(body)

  if (err != nil) {
    t.Fatalf("Unexpected error: %s\n", err)
  }
  if methodName != "system.client_version" {
    t.Fatalf("Expected %s, got %s\n", "system.client_version", methodName)
  }
  if params != nil {
    t.Fatalf("Expected no params\n")
  }
}
*/

/*
func TestParseJsonRequestStringParams(t *testing.T) {
  json := `{"methodName": "d.multicall",
            "params": [
              {"string": "main"},
              {"string": "d.base_filename="},
              {"string": "d.base_path="}]}`
  body := bytes.NewReader([]byte(json))

  methodName, params, err := ParseJsonRequest(body)

  if (err != nil) {
    t.Fatalf("Unexpected error: %s\n", err)
  }
  if methodName != "d.multicall" {
    t.Fatalf("Expected %s, got %s\n", "d.multicall", methodName)
  }
  if params == nil {
    t.Fatalf("Expected params\n")
  }
  if len(params) != 3 {
    t.Fatalf("Expected 3 params, got %d\n", len(params))
  }
}
*/
