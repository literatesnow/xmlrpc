package util

import (
	xml "encoding/xml"
)

func Start(encoder *xml.Encoder, name string) {
	encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: name}})
}

func End(encoder *xml.Encoder, name string) {
	encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: name}})
}

func Empty(encoder *xml.Encoder, name string) {
	Start(encoder, name)
	End(encoder, name)
}

func CharData(encoder *xml.Encoder, str string) {
	if str != "" {
		encoder.EncodeToken(xml.CharData(str))
	}
}
