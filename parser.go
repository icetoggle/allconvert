package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"os"
	// "strconv"
)

type Parser interface {
	initApp(app *Application)
	SaveToFile(pkg, path string)
	Cleanup()
	traversal(path string) error
}

type XmlParser interface {
	Parser
	HandleStartElement(e xml.StartElement)
	HandleCharData(e xml.CharData)
	HandleComment(e xml.Comment)
	HandleDirective(e xml.Directive)
	HandleProcInst(e xml.ProcInst)
	HandleEndElement(e xml.EndElement)
}

type xmlNode struct {
	name     string
	attr     []xml.Attr
	value    string
	children [][]*xmlNode
	childkey map[string]int
}

type BaseParser struct {
	XmlParser
	app *Application

	root  *xmlNode
	node  *xmlNode
	stack []*xmlNode
}

func (parser *BaseParser) traversal(path string) error {
	var err error
	var token xml.Token
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.Strict = true
	token, err = decoder.Token()
	if err != nil {
		log.Fatalln(err)
	}

	for token != nil {
		switch v := token.(type) {
		case xml.StartElement:
			parser.HandleStartElement(token.(xml.StartElement))
		case xml.EndElement:
			parser.HandleEndElement(token.(xml.EndElement))
		case xml.Comment:
			data := bytes.TrimSpace(token.(xml.Comment))
			if 0 != len(data) {
				parser.HandleComment(data)
			}
		case xml.CharData:
			data := bytes.TrimSpace(token.(xml.CharData))
			if 0 != len(data) {
				parser.HandleCharData(data)
			}
		case xml.Directive:
			data := bytes.TrimSpace(token.(xml.Directive))
			if 0 != len(data) {
				parser.HandleCharData(data)
			}
		case xml.ProcInst:
			parser.HandleProcInst(token.(xml.ProcInst))
		default:
			log.Fatalf("%T:\n%v\n", v, token)
		}
		token, err = decoder.Token()

		if err != nil && err != io.EOF {
			return err
		}
	}

	return nil
}

func (p *BaseParser) HandleStartElement(e xml.StartElement) {
	var node *xmlNode = &xmlNode{e.Name.Local, e.Attr, "", make([][]*xmlNode, 0), make(map[string]int)}
	if nil == p.root {
		p.root = node
	}
	if nil != p.node {
		p.stack = append(p.stack, p.node)
		id, ok := p.node.childkey[node.name]
		if !ok {
			list := make([]*xmlNode, 1)
			p.node.children = append(p.node.children, list)
			list[0] = node
			p.node.childkey[node.name] = len(p.node.children) - 1
		} else {
			p.node.children[id] = append(p.node.children[id], node)
		}
	}
	p.node = node
}

func (p *BaseParser) HandleCharData(e xml.CharData) {
	p.node.value = string(e)
}

func (p *BaseParser) HandleComment(e xml.Comment) {
}

func (p *BaseParser) HandleDirective(e xml.Directive) {
}

func (p *BaseParser) HandleProcInst(e xml.ProcInst) {
}

func (p *BaseParser) HandleEndElement(e xml.EndElement) {
	l := len(p.stack)
	if 0 == l {
		return
	}
	p.node = p.stack[l-1]
	p.stack = p.stack[0 : l-1]
}

func (p *BaseParser) initApp(app *Application) {
	p.app = app
}

func (p *BaseParser) Cleanup() {
	p.root = nil
	p.node = nil
	p.stack = nil
}

func checkPrototype(children []*xmlNode, prototype string) bool {
	for _, child := range children {
		for _, attr := range child.attr {
			if attr.Name.Local == prototype {
				return true
			}
		}
	}
	return false
}

func getAttr(node *xmlNode, prototype string) string {
	for _, attr := range node.attr {
		if attr.Name.Local == prototype {
			return attr.Value
		}
	}
	return ""
}

type KeyValue struct {
	key    string
	value  interface{}
	isAttr bool
}

func toObject(node *xmlNode) interface{} {
	// m := make(map[string]interface{})

	m := make([]*KeyValue, 0)

	for _, attr := range node.attr {
		// m[attr.Name.Local] = attr.Value

		m = append(m, &KeyValue{attr.Name.Local, attr.Value, true})
	}

	var size int
	var list []interface{}
	var set []*KeyValue
	var obj interface{}
	var useSet bool

	for _, children := range node.children {
		useSet = checkPrototype(children, "id")
		size = len(children)

		if size > 1 {
			if useSet {
				set = make([]*KeyValue, 0)
				var childName string
				for _, child := range children {
					childName = child.name
					obj = toObject(child)
					if nil != obj {
						// set[getAttr(child, "id")] = obj
						set = append(set, &KeyValue{getAttr(child, "id"), obj, false})
					}

				}
				m = append(m, &KeyValue{childName, set, false})
			} else {
				// set = make([]*KeyValue, 0)
				// var childName string
				// listLen := 0
				// for _, child := range children {
				// 	childName = child.name
				// 	obj = toObject(child)
				// 	if nil != obj {
				// 		set = append(set, &KeyValue{strconv.Itoa(listLen), obj, false})
				// 		listLen = listLen + 1
				// 	}
				// }
				list = make([]interface{}, 0)
				var childName string

				for _, child := range children {
					childName = child.name
					obj = toObject(child)
					if nil != obj {
						list = append(list, obj)
					}
				}
				// m[name] = list
				m = append(m, &KeyValue{childName, list, false})
			}
		} else {
			obj = toObject(children[0])
			name := children[0].name
			if nil != obj {
				if useSet {
					set = make([]*KeyValue, 0)
					// set[getAttr(children[0], "id")] = obj
					set = append(set, &KeyValue{getAttr(children[0], "id"), obj, false})
					m = append(m, &KeyValue{children[0].name, set, false})
					// m[name] = set
				} else if name == "item" {
					list = append(make([]interface{}, 0, size), obj)
					// m[name] = list
					m = append(m, &KeyValue{children[0].name, list, false})
				} else {
					// m[name] = obj
					m = append(m, &KeyValue{children[0].name, obj, false})
				}
			}
		}
	}

	if node.value != "" {
		if 0 == len(m) {
			return node.value
		} else {
			// m["value"] = node.value
			m = append(m, &KeyValue{"value", node.value, false})
			return m
		}
	}

	if 0 == len(m) {
		return nil
	}

	return m
}
