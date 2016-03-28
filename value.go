package xmlrpc

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type Value struct {
	Int      *int32 //i4
	Boolean  *bool
	String   *string
	Double   *float64
	DateTime *string //dateTime.iso8601
	Base64   *string
	//Struct - unsupported
	Array []Value

	//Extensions
	Nil   *bool    //nil, ex:nil
	Byte  *byte    //i1, ex:i1
	Float *float32 //float, ex:float
	Long  *int64   //i8, ex:i8
	//Dom - unsupported
	Short *int16 //i2, ex:i2
}

func NewValueInt(val int32) Value {
	return Value{Int: &val}
}
func NewValueBoolean(val bool) Value {
	return Value{Boolean: &val}
}
func NewValueString(val string) Value {
	return Value{String: &val}
}
func NewValueDouble(val float64) Value {
	return Value{Double: &val}
}
func NewValueDateTime(val string) Value {
	return Value{DateTime: &val}
}
func NewValueBase64(val string) Value {
	return Value{Base64: &val}
}
func NewValueArray(val []Value) Value {
	return Value{Array: val}
}
func NewValueNil() Value {
	b := true
	return Value{Nil: &b}
}
func NewValueByte(val byte) Value {
	return Value{Byte: &val}
}
func NewValueFloat(val float32) Value {
	return Value{Float: &val}
}
func NewValueLong(val int64) Value {
	return Value{Long: &val}
}
func NewValueShort(val int16) Value {
	return Value{Short: &val}
}

func (v *Value) FromString(str string) {
	if v.Int != nil {
		i, _ := strconv.ParseInt(str, 10, 32)
		*v.Int = int32(i)
	} else if v.Boolean != nil {
		*v.Boolean, _ = strconv.ParseBool(str)
	} else if v.String != nil {
		*v.String = str
	} else if v.Double != nil {
		*v.Double, _ = strconv.ParseFloat(str, 64)
	} else if v.DateTime != nil {
		v.DateTime = &str
	} else if v.Base64 != nil {
		v.Base64 = &str
	} else if v.Array != nil {
		//noop
	} else if v.Nil != nil {
		//noop
	} else if v.Byte != nil {
		i, _ := strconv.Atoi(str)
		*v.Byte = byte(i)
	} else if v.Float != nil {
		f, _ := strconv.ParseFloat(str, 32)
		*v.Float = float32(f)
	} else if v.Long != nil {
		*v.Long, _ = strconv.ParseInt(str, 10, 64)
	} else if v.Short != nil {
		i, _ := strconv.ParseInt(str, 10, 16)
		*v.Short = int16(i)
	} else {
		v.String = &str
	}
}

func (v *Value) FromRpc(name string) {
	switch name {
	case "int", "i4":
		var val int32
		v.Int = &val
	case "boolean":
		var val bool
		v.Boolean = &val
	case "string":
		var val string
		v.String = &val
	case "double":
		var val float64
		v.Double = &val
	case "dateTime", "dateTime.iso8601":
		var val string
		v.DateTime = &val
	case "base64":
		var val string
		v.Base64 = &val
	case "array":
		v.Array = make([]Value, 0)
	case "none", "nil", "ex:nil":
		var val bool = true
		v.Nil = &val
	case "byte", "i1", "ex:i1":
		var val byte
		v.Byte = &val
	case "float", "ex:float":
		var val float32
		v.Float = &val
	case "long", "i8", "ex:i8":
		var val int64
		v.Long = &val
	case "short", "i2", "ex:i2":
		var val int16 = 0
		v.Short = &val
	}
}

func (v *Value) asString() (dataType string, text string) {
	if v.Int != nil {
		return "int", strconv.FormatInt(int64(*v.Int), 10)
	} else if v.Boolean != nil {
		return "boolean", strconv.FormatBool(*v.Boolean)
	} else if v.String != nil {
		return "string", *v.String
	} else if v.Double != nil {
		return "double", strconv.FormatFloat(*v.Double, 'f', -1, 64)
	} else if v.DateTime != nil {
		return "dateTime.iso8601", *v.DateTime
	} else if v.Base64 != nil {
		return "base64", *v.Base64
	} else if v.Nil != nil {
		return "ex:nil", ""
	} else if v.Byte != nil {
		return "ex:i1", strconv.Itoa(int(*v.Byte))
	} else if v.Float != nil {
		return "ex:float", strconv.FormatFloat(float64(*v.Float), 'f', -1, 32)
	} else if v.Long != nil {
		return "ex:i8", strconv.FormatInt(*v.Long, 10)
	} else if v.Short != nil {
		return "ex:i2", strconv.FormatInt(int64(*v.Short), 10)
	} else if v.Array != nil {
		return "array", v.asStringArray(v.Array)
	} else {
		return "empty", ""
	}
}

func (v *Value) asStringArray(values []Value) (text string) {
	parts := make([]string, len(values))

	for i, val := range values {
		parts[i] = val.Print()
	}

	return "[" + strings.Join(parts, ", ") + "]"
}

func (v *Value) Print() (text string) {
	dataType, text := v.asString()
	return "{" + dataType + " " + text + "}"
}

func (v *Value) asXml(encoder *xml.Encoder) {
	xmlStart(encoder, "value")

	dataType, text := v.asString()

	switch dataType {
	case "ex:nil":
		xmlEmpty(encoder, dataType)
	case "array":
		v.xmlArrayValue(encoder, v.Array)
	default:
		v.xmlValue(encoder, dataType, text)
	}

	xmlEnd(encoder, "value")
}

func (v *Value) xmlValue(encoder *xml.Encoder, name string, value string) {
	xmlStart(encoder, name)
	xmlCharData(encoder, value)
	xmlEnd(encoder, name)
}

func (v *Value) xmlArrayValue(encoder *xml.Encoder, values []Value) {
	xmlStart(encoder, "array")
	xmlStart(encoder, "data")

	for _, val := range values {
		val.asXml(encoder)
	}

	xmlEnd(encoder, "data")
	xmlEnd(encoder, "array")
}
