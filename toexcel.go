package main

import (
	"github.com/tealeg/xlsx"
	"log"
	"strconv"
)

type ExcelParser struct {
	BaseParser

	key2col   map[string]int
	valueList map[int](map[int]interface{})

	headerList []string

	colNum int

	rowNum int

	path string

	trie *Trie
}

func newExcelParser() *ExcelParser {
	rt := new(ExcelParser)
	rt.key2col = make(map[string]int)
	rt.valueList = make(map[int](map[int]interface{}))
	rt.headerList = make([]string, 0)
	rt.colNum = 0
	rt.rowNum = 0

	return rt
}

func (p *ExcelParser) Cleanup() {
	p.BaseParser.Cleanup()
	p.key2col = make(map[string]int)
	p.valueList = make(map[int](map[int]interface{}))
	p.headerList = make([]string, 0)
	p.colNum = 0
	p.rowNum = 0
}

func (p *ExcelParser) toListObject(row int, prefix string, headpre string, obj interface{}) {
	// var result string = ""
	var result interface{} = nil

	var isValue = false

	switch v := obj.(type) {
	case []interface{}:
		list := obj.([]interface{})
		for i, child := range list {
			newPrefix := prefix + "." + strconv.Itoa(i)
			headpre := headpre + "*"
			p.toListObject(row, newPrefix, headpre, child)
		}
	case []*KeyValue:
		set := obj.([]*KeyValue)
		for _, keyValue := range set {
			// key := keyValue.key
			child := keyValue.value

			key := string2value(keyValue.key)

			newheadpre := headpre
			switch key.(type) {
			case string:
				newheadpre = newheadpre + "." + key.(string)
			default:
				newheadpre = newheadpre + "*"
			}

			newPrefix := prefix + "." + keyValue.key

			if keyValue.isAttr {
				// newPrefix = newPrefix + "$"
				newheadpre = newheadpre + "$"
			}
			p.toListObject(row, newPrefix, newheadpre, child)
		}
	case string:
		var s string = obj.(string)

		if s == "" {
			return
		} else if _, err := strconv.ParseInt(s, 10, 64); err == nil {
			result = s
		} else if _, err := strconv.ParseFloat(s, 64); err == nil {
			result = s
		} else {
			result = obj.(string)
		}
		isValue = true

	case float64:
		isValue = true
		var f float64 = obj.(float64)
		result = strconv.FormatFloat(f, byte('f'), -1, 64)
	case bool:
		isValue = true
		var b bool = obj.(bool)
		result = strconv.FormatBool(b)
	default:
		result = "nil"
		log.Printf("%T:\n%v\n", v, obj)
	}
	if isValue {
		var id = p.trie.findKey(prefix)
		p.key2col[prefix] = id
		p.headerList[id] = headpre
		p.valueList[row][id] = result
	}
}

func (p *ExcelParser) ToExcel() *xlsx.File {
	// cellGrid := make([][]*xlsx.Cell, p.rowNum)
	file := xlsx.NewFile()
	sheet := file.AddSheet("Sheet1")

	sheet.AddRow()

	row := sheet.AddRow()

	for i := 0; i < p.colNum; i++ {
		cell := row.AddCell()
		cell.Value = p.headerList[i]
	}

	for i := 1; i <= p.rowNum; i++ {
		row := sheet.AddRow()
		// cellGrid[i-1] = make([]*xlsx.Cell, p.colNum)

		for j := 0; j < p.colNum; j++ {
			cell := row.AddCell()
			value, present := p.valueList[i][j]
			if present {
				switch value.(type) {
				case bool:
					cell.SetBool(value.(bool))
				case float64:
					cell.SetFloat(value.(float64))
				case string:
					cell.SetString(value.(string))
				}
			}
		}
	}

	return file

}

func (p *ExcelParser) ToObject() {
	var obj interface{} = toObject(p.root)

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("parse ", p.path, r)
	// 	}
	// }()

	trie := newTrie()
	p.trie = trie

	switch obj.(type) {
	case []*KeyValue:
		set := obj.([]*KeyValue)
		p.rowNum = 1

		if len(set) == 1 {

			for _, keyValue := range set {
				prefix := keyValue.key
				child := keyValue.value

				set := child.([]*KeyValue)

				trieNode := new(TrieNode)
				trieNode.key = prefix
				trieNode.childNode = make([]*TrieNode, 0)
				trie.root.childNode = append(trie.root.childNode, trieNode)

				for _, value := range set {

					// p.valueList[p.rowNum] = make(map[int]interface{})
					// p.toListObject(p.rowNum, prefix, prefix, value.value)
					trie.pushXmlNode(trieNode, value.value)
					// p.rowNum = p.rowNum + 1
				}

			}
			trie.recordId(trie.root)
			// trie.dfs()

			p.headerList = make([]string, p.trie.id)
			p.colNum = p.trie.id

			for _, keyValue := range set {
				prefix := keyValue.key
				child := keyValue.value

				set := child.([]*KeyValue)

				for _, value := range set {
					p.valueList[p.rowNum] = make(map[int]interface{})
					p.toListObject(p.rowNum, prefix, prefix, value.value)
					p.rowNum = p.rowNum + 1
				}
			}

		} else {
			p.valueList[p.rowNum] = make(map[int]interface{})
			trie.pushXmlNode(trie.root, set)

			trie.recordId(trie.root)
			// trie.dfs()
			p.colNum = trie.id
			p.headerList = make([]string, p.trie.id)

			for _, keyValue := range set {
				child := keyValue.value

				newPrefix := keyValue.key
				p.toListObject(1, newPrefix, newPrefix, child)
			}
		}
	case string:
		// result = p.root.name + "=" + obj.(string)
		p.key2col[p.root.name] = 0
		p.valueList[1][0] = obj.(string)
		p.colNum = p.colNum + 1
	}
	//return toLuaObject(obj)
}

func (p *ExcelParser) SaveToFile(pkg, path string) {
	p.ToObject()
	file := p.ToExcel()
	file.Save(path)
}
