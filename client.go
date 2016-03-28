package xmlrpc

import (
	"bufio"
	"bytes"
	"encoding/xml"
)

func CreateRequest(methodName string, values []Value) (document []byte) {
	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)
	writer.WriteString(xml.Header)

	encoder := xml.NewEncoder(writer)

	xmlStart(encoder, "methodCall")
	xmlStart(encoder, "methodName")
	xmlCharData(encoder, methodName)
	xmlEnd(encoder, "methodName")
	xmlParams(encoder, values)
	xmlEnd(encoder, "methodCall")

	encoder.Flush()

	return buf.Bytes()
}

func ParseResponse(response *bytes.Buffer) (values []Value) {
	decoder := xml.NewDecoder(response)
	values = make([]Value, 0)

	if !nextElem(decoder, "methodResponse") ||
		!nextElem(decoder, "params") {
		return values
	}

	for {
		if !nextElem(decoder, "param") {
			break
		}

		if value := parseValue(decoder); value != nil {
			values = append(values, *value)
		}
	}

	return values
}

func parseValue(decoder *xml.Decoder) (value *Value) {
	value = nil

	for {
		token, err := decoder.Token()
		if token == nil || err != nil {
			break
		}

		switch elem := token.(type) {
		case xml.CharData:
			parseCharData(value, string(elem))
		case xml.StartElement:
			parseStartElement(decoder, &value, elem.Name.Local)
		case xml.EndElement:
			if elem.Name.Local == "value" || value == nil {
				return value
			}
		}
	}

	return nil
}

func parseCharData(value *Value, str string) {
	if value == nil {
		return
	}

	value.FromString(str)
}

func parseStartElement(decoder *xml.Decoder, valuePtr **Value, elemName string) {
	if elemName == "value" {
		*valuePtr = &Value{}
		return
	}

	if *valuePtr == nil {
		return
	}

	value := *valuePtr
	value.FromRpc(elemName)

	if value.Array != nil {
		parseValueArray(decoder, value)
	}
}

func parseValueArray(decoder *xml.Decoder, value *Value) {
	if !nextElem(decoder, "data") {
		return //TODO return error
	}

	for {
		if val := parseValue(decoder); val != nil {
			value.Array = append(value.Array, *val)
		} else {
			break
		}
	}
}

func nextElem(decoder *xml.Decoder, name string) (found bool) {
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

func xmlStart(encoder *xml.Encoder, name string) {
	encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: name}})
}

func xmlEnd(encoder *xml.Encoder, name string) {
	encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: name}})
}

func xmlEmpty(encoder *xml.Encoder, name string) {
	xmlStart(encoder, name)
	xmlEnd(encoder, name)
}

func xmlCharData(encoder *xml.Encoder, str string) {
	if str != "" {
		encoder.EncodeToken(xml.CharData(str))
	}
}

func xmlParams(encoder *xml.Encoder, values []Value) {
	if len(values) == 0 {
		return
	}

	xmlStart(encoder, "params")

	for _, val := range values {
		xmlStart(encoder, "param")
		val.asXml(encoder)
		xmlEnd(encoder, "param")
	}

	xmlEnd(encoder, "params")
}
