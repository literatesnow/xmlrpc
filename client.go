package xmlrpc

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"strconv"
)

type Client struct {
}

func (cl *Client) CreateRequest(methodName string, values []Value) (document []byte) {
	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)
	writer.WriteString(xml.Header)

	encoder := xml.NewEncoder(writer)

	appendStart(encoder, "methodCall")
	appendStart(encoder, "methodName")
	appendCharData(encoder, methodName)
	appendEnd(encoder, "methodName")
	appendParams(encoder, values)
	appendEnd(encoder, "methodCall")

	encoder.Flush()

	return buf.Bytes()
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
			value.Number, _ = strconv.ParseInt(str, 10, 64)
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

func appendStart(encoder *xml.Encoder, name string) {
	encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: name}})
}

func appendEnd(encoder *xml.Encoder, name string) {
	encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: name}})
}

func appendCharData(encoder *xml.Encoder, str string) {
	if str != "" {
		encoder.EncodeToken(xml.CharData(str))
	}
}

func appendParams(encoder *xml.Encoder, values []Value) {
	if len(values) == 0 {
		return
	}

	appendStart(encoder, "params")

	for _, val := range values {
		appendStart(encoder, "param")
		val.appendXml(encoder)
		appendEnd(encoder, "param")
	}

	appendEnd(encoder, "params")
}
