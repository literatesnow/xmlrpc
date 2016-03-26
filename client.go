package xmlrpc

import (
	"bytes"
	"encoding/xml"
	"strconv"
)

type Client struct {
	doc bytes.Buffer
}

func NewClient() *Client {
	cl := &Client{}

	return cl
}

func (cl *Client) Begin(methodName string) {
	cl.doc.WriteString("<?xml version=\"1.0\"?><methodCall><methodName>")
	cl.doc.WriteString(methodName)
	cl.doc.WriteString("</methodName><params>")
}

func (cl *Client) End() (doc string) {
	cl.doc.WriteString("</params></methodCall>")
	return cl.doc.String()
}

func (cl *Client) ParamString(value string) {
	cl.doc.WriteString("<param><value><string>")
	cl.doc.WriteString(value)
	cl.doc.WriteString("</string></value></param>")
}

func (cl *Client) Parse(response *bytes.Buffer) (values []Value) {
	decoder := xml.NewDecoder(response)
	values = make([]Value, 0)

	if !cl.nextElem(decoder, "methodResponse") ||
		!cl.nextElem(decoder, "params") {
		return values
	}

	for {
		if !cl.nextElem(decoder, "param") {
			break
		}

		if value := cl.parseValue(decoder); value != nil {
			values = append(values, *value)
		}
	}

	return values
}

func (cl *Client) parseValue(decoder *xml.Decoder) (value *Value) {
	value = nil

	for {
		token, err := decoder.Token()
		if token == nil || err != nil {
			break
		}

		switch elem := token.(type) {
		case xml.CharData:
			cl.parseCharData(value, string(elem))
		case xml.StartElement:
			cl.parseStartElement(decoder, &value, elem.Name.Local)
		case xml.EndElement:
			if elem.Name.Local == "value" || value == nil {
				return value
			}
		}
	}

	return nil
}

func (cl *Client) parseCharData(value *Value, str string) {
	if value == nil {
		return
	}

	if value.Type == 0 {
		value.Type = ValueString
	}

	if str == "" {
		value.Type = ValueEmpty

	} else {
		switch value.Type {
		case ValueInt:
			value.Number, _ = strconv.Atoi(str)
		case ValueDouble:
			value.Double, _ = strconv.ParseFloat(str, 64)
		case ValueBoolean:
			value.Boolean, _ = strconv.ParseBool(str)
		case ValueString, ValueDate, ValueBase64:
			value.String = str
		}
	}
}

func (cl *Client) parseStartElement(decoder *xml.Decoder, valuePtr **Value, elemName string) {
	if elemName == "value" {
		*valuePtr = &Value{Type: 0}
		return
	}

	value := *valuePtr
	value.setType(elemName)

	if value.Type == 0 {
		return
	}

	if value.Type == ValueArray {
		cl.parseValueArray(decoder, value)
	}
}

func (cl *Client) parseValueArray(decoder *xml.Decoder, value *Value) {
	if !cl.nextElem(decoder, "data") {
		return //TODO return error
	}

	if value.Array == nil {
		value.Array = make([]Value, 0)
	}

	for {
		if val := cl.parseValue(decoder); val != nil {
			value.Array = append(value.Array, *val)
		} else {
			break
		}
	}
}

func (cl *Client) nextElem(decoder *xml.Decoder, name string) (found bool) {
	for {
		token, err := decoder.Token()
		if token == nil || err != nil {
			return false
		}

		switch elem := token.(type) {
		case xml.StartElement:
			if elem.Name.Local == name {
				return true
			}
		}
	}

	return false
}
