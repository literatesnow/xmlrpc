package xmlrpc

import (
	"encoding/xml"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/unrulyknight/xmlrpc/util"
)

const (
	iso8601 = "2006-01-02T15:04:05-0700"
)

type Value struct {
	Int      *int32     `json:"int,omitempty"` //i4
	Boolean  *bool      `json:"boolean,omitempty"`
	String   *string    `json:"string,omitempty"`
	Double   *float64   `json:"double,omitempty"`
	DateTime *time.Time `json:"dateTime8601,omitempty"` //dateTime.iso8601
	Base64   *string    `json:"base64,omitempty"`
	Array    []Value    `json:"array,omitempty"`
	//Struct - unsupported

	//Extensions
	Nil   *bool    `json:"nil,omitempty"`   //nil, ex:nil
	Byte  *byte    `json:"i1,omitempty"`    //i1, ex:i1
	Float *float32 `json:"float,omitempty"` //float, ex:float
	Long  *int64   `json:"i8,omitempty"`    //i8, ex:i8
	Short *int16   `json:"i2,omitempty"`    //i2, ex:i2
	//Dom - unsupported
}

func NewInt(val int32) Value {
	return Value{Int: &val}
}
func NewBoolean(val bool) Value {
	return Value{Boolean: &val}
}
func NewString(val string) Value {
	return Value{String: &val}
}
func NewDouble(val float64) Value {
	return Value{Double: &val}
}
func NewDateTime(val time.Time) Value {
	return Value{DateTime: &val}
}
func NewBase64(val string) Value {
	return Value{Base64: &val}
}
func NewArray(val []Value) Value {
	return Value{Array: val}
}
func NewNil() Value {
	b := true
	return Value{Nil: &b}
}
func NewByte(val byte) Value {
	return Value{Byte: &val}
}
func NewFloat(val float32) Value {
	return Value{Float: &val}
}
func NewLong(val int64) Value {
	return Value{Long: &val}
}
func NewShort(val int16) Value {
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
		*v.DateTime, _ = time.Parse(iso8601, str)
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

func (v *Value) FromNumber(num float64) {
	if v.Int != nil {
		*v.Int = int32(num)
	} else if v.Double != nil {
		*v.Double = num
	} else if v.Byte != nil {
		*v.Byte = byte(num)
	} else if v.Float != nil {
		*v.Float = float32(num)
	} else if v.Long != nil {
		*v.Long = int64(num)
	} else if v.Short != nil {
		*v.Short = int16(num)
	}
}

func (v *Value) FromBoolean(b bool) {
	if v.Boolean != nil {
		*v.Boolean = b
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
	case "dateTime", "dateTime8601", "dateTime.iso8601":
		var val time.Time
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
		if *v.Boolean {
			return "boolean", "1"
		} else {
			return "boolean", "0"
		}
	} else if v.String != nil {
		return "string", *v.String
	} else if v.Double != nil {
		return "double", strconv.FormatFloat(*v.Double, 'f', -1, 64)
	} else if v.DateTime != nil {
		return "dateTime.iso8601", (*v.DateTime).Format(iso8601)
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
	util.Start(encoder, "value")

	dataType, text := v.asString()

	switch dataType {
	case "ex:nil":
		util.Empty(encoder, dataType)
	case "array":
		v.xmlArrayValue(encoder, v.Array)
	default:
		v.xmlValue(encoder, dataType, text)
	}

	util.End(encoder, "value")
}

func (v *Value) xmlValue(encoder *xml.Encoder, name string, value string) {
	util.Start(encoder, name)
	util.CharData(encoder, value)
	util.End(encoder, name)
}

func (v *Value) xmlArrayValue(encoder *xml.Encoder, values []Value) {
	util.Start(encoder, "array")
	util.Start(encoder, "data")

	for _, val := range values {
		val.asXml(encoder)
	}

	util.End(encoder, "data")
	util.End(encoder, "array")
}
