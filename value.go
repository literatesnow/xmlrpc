package xmlrpc

import (
	"encoding/xml"
	"strconv"
)

const (
	ValueEmpty   = 1
	ValueInt     = 2
	ValueDouble  = 3
	ValueBoolean = 4
	ValueString  = 5
	ValueDate    = 6
	ValueBase64  = 7
	ValueArray   = 8
)

type Value struct {
	Type    int     //internal
	Number  int64   //i4 i8 int
	Boolean bool    //boolean
	Double  float64 //double
	String  string  //string
	Date    string  //dateTime.iso8601
	Base64  string  //base64
	Array   []Value //array
}

func (v *Value) appendXml(encoder *xml.Encoder) {
	appendStart(encoder, "value")

	switch v.Type {
	case ValueInt:
		v.appendElementXml(encoder, "int", strconv.FormatInt(v.Number, 10))
	case ValueDouble:
		v.appendElementXml(encoder, "double", strconv.FormatFloat(v.Double, 'f', -1, 64))
	case ValueBoolean:
		v.appendElementXml(encoder, "boolean", strconv.FormatBool(v.Boolean))
	case ValueString:
		v.appendElementXml(encoder, "string", v.String)
	case ValueDate:
		v.appendElementXml(encoder, "dateTime.iso8601", v.Date)
	case ValueBase64:
		v.appendElementXml(encoder, "base64", v.Base64)
	case ValueArray:
		v.appendArrayXml(encoder, v.Array)
	}

	appendEnd(encoder, "value")
}

func (v *Value) appendElementXml(encoder *xml.Encoder, name string, value string) {
	appendStart(encoder, name)
	appendCharData(encoder, value)
	appendEnd(encoder, name)
}

func (v *Value) appendArrayXml(encoder *xml.Encoder, values []Value) {
	appendStart(encoder, "array")
	appendStart(encoder, "data")

	for _, val := range values {
		val.appendXml(encoder)
	}

	appendEnd(encoder, "data")
	appendEnd(encoder, "array")
}

func (v *Value) setType(rpcName string) {
	switch rpcName {
	case "int", "i4", "i8":
		v.Type = ValueInt
	case "double":
		v.Type = ValueDouble
	case "boolean":
		v.Type = ValueBoolean
	case "string":
		v.Type = ValueString
	case "dateTime.iso8601":
		v.Type = ValueDate
	case "base64":
		v.Type = ValueBase64
	case "array":
		v.Type = ValueArray
	default:
		v.Type = 0
	}
}
