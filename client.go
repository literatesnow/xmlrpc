package xmlrpc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"

	"bitbucket.org/unrulyknight/xmlrpc/util"
)

const (
	curlyLeft   = json.Delim('{')
	curlyRight  = json.Delim('}')
	squareLeft  = json.Delim('[')
	squareRight = json.Delim(']')
)

func CreateRequest(methodName string, params []Value) (document []byte) {
	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)
	writer.WriteString(xml.Header)

	encoder := xml.NewEncoder(writer)

	util.Start(encoder, "methodCall")
	util.Start(encoder, "methodName")
	util.CharData(encoder, methodName)
	util.End(encoder, "methodName")
	xmlParams(encoder, params)
	util.End(encoder, "methodCall")

	encoder.Flush()

	return buf.Bytes()
}

func ParseResponse(response *bytes.Buffer) (value *Value, err error) {
	decoder := xml.NewDecoder(response)

	if !nextElem(decoder, "methodResponse") ||
		!nextElem(decoder, "params") ||
		!nextElem(decoder, "param") {
		return nil, errors.New("Invalid response")
	}

	value, err = parseValue(decoder)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func parseValue(decoder *xml.Decoder) (value *Value, err error) {
	value = nil

	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		if token == nil {
			break
		}

		switch elem := token.(type) {
		case xml.CharData:
			parseCharData(value, string(elem))
		case xml.StartElement:
			if err = parseStartElement(decoder, &value, elem.Name.Local); err != nil {
				return nil, err
			}
		case xml.EndElement:
			if elem.Name.Local == "value" || value == nil {
				return value, nil
			}
		}
	}

	return nil, nil
}

func parseCharData(value *Value, str string) {
	if value == nil {
		return
	}

	value.FromString(str)
}

func parseStartElement(decoder *xml.Decoder, valuePtr **Value, elemName string) (err error) {
	if elemName == "value" {
		*valuePtr = &Value{}
		return nil
	}

	if *valuePtr == nil {
		return nil
	}

	value := *valuePtr
	value.FromRpc(elemName)

	if value.Array != nil {
		if err = parseValueArray(decoder, value); err != nil {
			return err
		}
	}

	return nil
}

func parseValueArray(decoder *xml.Decoder, value *Value) (err error) {
	if !nextElem(decoder, "data") {
		return errors.New("Unexpected element")
	}

	for {
		val, err := parseValue(decoder)
		if err != nil {
			return err
		}
		if val != nil {
			value.Array = append(value.Array, *val)
		} else {
			break
		}
	}

	return nil
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

func xmlParams(encoder *xml.Encoder, values []Value) {
	if len(values) == 0 {
		return
	}

	util.Start(encoder, "params")

	for _, val := range values {
		util.Start(encoder, "param")
		val.asXml(encoder)
		util.End(encoder, "param")
	}

	util.End(encoder, "params")
}

//JSON

func ParseJsonRequest(body io.Reader) (methodName string, params []Value, err error) {
	decoder := json.NewDecoder(body)

	if err := nextJsonDelim(decoder, curlyLeft); err != nil {
		return "", nil, err
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", nil, err
		}

		switch token {
		case "methodName":
			if methodName, err = nextJsonString(decoder); err != nil {
				return "", nil, err
			}
		case "params":
			if params, err = parseJsonValues(decoder); err != nil {
				return "", nil, err
			}
		}
	}

	return methodName, params, nil
}

func nextJsonString(decoder *json.Decoder) (delim string, err error) {
	token, err := decoder.Token()
	if err != nil {
		return "", err
	}

	if str, ok := token.(string); ok {
		return str, nil
	}

	return "", errors.New("Not a string")
}

func nextJsonDelim(decoder *json.Decoder, delim json.Delim) (err error) {
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	if token != delim {
		return errors.New("Unexpected next token")
	}

	return nil
}

func parseJsonValues(decoder *json.Decoder) (values []Value, err error) {
	if err := nextJsonDelim(decoder, squareLeft); err != nil {
		return nil, err
	}

	values = make([]Value, 0)
	hasType := false

	var val *Value

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if token == squareRight {
			break
		}

		if token == curlyLeft {
			val = &Value{}

		} else if token == curlyRight {
			if val != nil {
				values = append(values, *val)
				val = nil
				hasType = false
			}

		} else if val != nil {
			if !hasType {
				if str, ok := token.(string); ok {
					val.FromRpc(str)
					hasType = true

					if val.Array != nil {
						val.Array, err = parseJsonValues(decoder)
						if err != nil {
							return nil, err
						}
					}

				} else {
					return nil, errors.New("Invalid token")
				}
			} else {
				switch p := token.(type) {
				case string:
					val.FromString(p)
				case float64:
					val.FromNumber(p)
				case bool:
					val.FromBoolean(p)
				case nil:
				default:
					return nil, errors.New("Unexpected token")
				}
			}
		}
	}

	return values, nil
}
