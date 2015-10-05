package main

import (
	"bytes"
	"encoding/xml"
	"log"
	"os"
)

func makeStartXml(root *XlsObjNode, enc *xml.Encoder) {
	valueList := root.value.([]*XlsObjNode)
	token := xml.StartElement{xml.Name{"", root.key}, make([]xml.Attr, 0)}

	for _, key := range valueList {
		if key.nodeType == ATTR_TYPE {
			attr := xml.Attr{xml.Name{"", key.key}, key.value.(string)}
			token.Attr = append(token.Attr, attr)
		}
	}
	enc.EncodeToken(token)
}

func makeEndXml(root *XlsObjNode, enc *xml.Encoder) {
	enc.EncodeToken(xml.EndElement{xml.Name{"", root.key}})
}

func xlsToXml(root *XlsObjNode, enc *xml.Encoder) {
	valueList := root.value.([]*XlsObjNode)

	makeStartXml(root, enc)
	for _, key := range valueList {
		switch key.nodeType {
		case VALUE_TYPE:
			enc.EncodeToken(xml.StartElement{xml.Name{"", key.key}, make([]xml.Attr, 0)})
			enc.EncodeToken(xml.CharData(key.value.(string)))
			enc.EncodeToken(xml.EndElement{xml.Name{"", key.key}})

		case HASH_TYPE:
			hasId := false

			valueList2 := key.value.([]*XlsObjNode)
			for _, child := range valueList2 {
				if child.isId {
					hasId = true
					break
				}
			}

			if hasId {
				for _, child := range valueList2 {
					child.key = key.key
					xlsToXml(child, enc)
				}

			} else {
				// hasId := false
				// valueList2 := key.value.([]*XlsObjNode)
				// for _, child := range valueList2 {
				// 	if child.isId {
				// 		hasId = true
				// 		break
				// 	}
				// }

				// if hasId {
				// 	for _, child := range valueList2 {
				// 		child.key = key.key
				// 		xlsToXml(key, enc)
				// 	}

				// } else {
				xlsToXml(key, enc)
			}
			// }
		case ARRAY_TYPE:
			childList := key.value.([]*XlsObjNode)
			for _, childKey := range childList {
				if childKey.nodeType == VALUE_TYPE {
					enc.EncodeToken(xml.StartElement{xml.Name{"", key.key}, make([]xml.Attr, 0)})
					enc.EncodeToken(xml.CharData(childKey.value.(string)))
					enc.EncodeToken(xml.EndElement{xml.Name{"", key.key}})
				} else {
					childKey.key = key.key
					xlsToXml(childKey, enc)
				}
			}
		}
	}
	makeEndXml(root, enc)
}

type XlsToXmlParser struct {
	BaseXlsParser
}

func (this *XlsToXmlParser) SaveToFile(pkg, path string) {
	var out bytes.Buffer
	enc := xml.NewEncoder(&out)
	enc.Indent("", "\t")
	xlsToXml(this.root, enc)
	enc.Flush()
	fout, err := os.Create(path)
	if nil != err {
		log.Fatalln(err)
	}
	s := out.String()

	fout.WriteString(s)
}
