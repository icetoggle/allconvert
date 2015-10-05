package main

import (
	"github.com/tealeg/xlsx"
	"strings"
)

const (
	VALUE_TYPE = iota
	ARRAY_TYPE
	ARRAY_ITEM
	HASH_TYPE
	ATTR_TYPE
)

type XlsObjNode struct {
	key      string
	value    interface{}
	nodeType int // 0 值类型， 1 数组类型 2 数组子类型 3 哈希类型 4 属性类型

	isId bool
}

func toXlsObject(file *xlsx.File) *XlsObjNode {
	for _, sheet := range file.Sheets {
		rowNum := len(sheet.Rows) - 2

		colNum := len(sheet.Rows[1].Cells)

		headerList := make([][]string, colNum)
		for i := 0; i < colNum; i++ {
			headerList[i] = strings.Split(sheet.Rows[1].Cells[i].String(), ".")
		}

		var root = new(XlsObjNode)

		root.key = "Root"

		root.nodeType = HASH_TYPE

		for i := 0; i < rowNum; i++ {
			hasRow := false

			for j := 0; j < colNum; j++ {
				if strings.TrimSpace(sheet.Rows[i+2].Cells[j].String()) != "" {
					hasRow = true
				}
			}

			if !hasRow {
				rowNum = i
				break
			}
		}

		if rowNum == 1 {
			root.value = make([]*XlsObjNode, 0)
			root.nodeType = HASH_TYPE
			for i := 0; i < colNum; i++ {
				if strings.TrimSpace(sheet.Rows[2].Cells[i].String()) != "" {
					xlsInsertValue(root, headerList[i], sheet.Rows[2].Cells[i], 0)
				}
			}

		} else {
			root.value = make([]*XlsObjNode, 1)

			item := new(XlsObjNode)
			item.key = headerList[0][0]
			item.value = make([]*XlsObjNode, rowNum)
			item.nodeType = ARRAY_TYPE

			root.value.([]*XlsObjNode)[0] = item
			valueList := item.value.([]*XlsObjNode)

			for i := 0; i < rowNum; i++ {
				valueList[i] = new(XlsObjNode)
				child := valueList[i]
				child.key = ""
				child.nodeType = HASH_TYPE
				child.value = make([]*XlsObjNode, 0)
				for j := 0; j < colNum; j++ {
					if strings.TrimSpace(sheet.Rows[i+2].Cells[j].String()) != "" {
						xlsInsertValue(child, headerList[j], sheet.Rows[i+2].Cells[j], 1)
					}
				}
			}
		}

		return root

	}

	return nil

}

func xlsInsertValue(tree *XlsObjNode, keys []string, value *xlsx.Cell, offset int) {

	if tree.nodeType == ARRAY_TYPE {
		var valueList = tree.value.([]*XlsObjNode)
		var valueLen = len(valueList)
		var node *XlsObjNode

		// fmt.Println(keys[offset])
		// for _, child := range(tree.value.([]XlsObjNode)){
		// 	if findKey(tre)
		// }
		// needNew := findKey(tree, keys, offset)
		// dfs(tree)

		if valueLen == 0 {
			node = new(XlsObjNode)
			if keys[offset] == "id" {
				node.key = keys[offset]
			}
			node.nodeType = HASH_TYPE
			node.value = make([]*XlsObjNode, 0)
			tree.value = append(valueList, node)
		} else {
			node = valueList[valueLen-1]
		}

		needNew := xlsFindKey(node, keys, offset)

		if needNew {
			node = new(XlsObjNode)
			node.nodeType = HASH_TYPE
			node.value = make([]*XlsObjNode, 0)
			tree.value = append(valueList, node)
		}
		xlsInsertValue(node, keys, value, offset)
		return
	}

	if len(keys) == offset+1 {
		var child = new(XlsObjNode)

		if strings.HasSuffix(keys[offset], "*") {
			child.key = ""
			child.nodeType = VALUE_TYPE
			child.value = value.String()

			var valueList = tree.value.([]*XlsObjNode)
			for _, key := range valueList {
				key2 := keys[offset]
				if key2[:(len(key2)-1)] == key.key {
					// xlsInsertValue(key, keys, value, offset+1)
					key.value = append(key.value.([]*XlsObjNode), child)
					return
				}
			}

			var newXlsObjNode = new(XlsObjNode)
			tree.value = append(valueList, newXlsObjNode)
			newXlsObjNode.nodeType = ARRAY_TYPE
			newXlsObjNode.key = keys[offset][:len(keys[offset])-1]
			newXlsObjNode.value = make([]*XlsObjNode, 1)
			newXlsObjNode.value.([]*XlsObjNode)[0] = child
			return
		}
		child.key = keys[offset]

		if strings.HasSuffix(keys[offset], "$") { //属性类型
			child.nodeType = ATTR_TYPE
			child.key = child.key[:len(child.key)-1]
		} else {
			child.nodeType = VALUE_TYPE
		}
		child.value = value.String()
		tree.value = append(tree.value.([]*XlsObjNode), child)
		if child.key == "id" {
			tree.key = child.value.(string)
			tree.isId = true
		}

		return
	}

	var valueList = tree.value.([]*XlsObjNode)
	for _, key := range valueList {
		key2 := keys[offset]
		if key.key == keys[offset] || (strings.HasSuffix(key2, "*") && key2[:(len(key2)-1)] == key.key) {
			xlsInsertValue(key, keys, value, offset+1)
			return
		}
	}

	var newXlsObjNode = new(XlsObjNode)
	tree.value = append(valueList, newXlsObjNode)
	newXlsObjNode.key = keys[offset]
	newXlsObjNode.value = make([]*XlsObjNode, 0)

	if strings.HasSuffix(keys[offset], "*") { //数组类型
		newXlsObjNode.nodeType = ARRAY_TYPE
		newXlsObjNode.key = keys[offset][:len(keys[offset])-1]
	} else if strings.HasPrefix(keys[offset+1], "id") && len(keys) == 2 {
		newXlsObjNode.nodeType = ARRAY_TYPE
	} else {
		newXlsObjNode.nodeType = HASH_TYPE
	}

	xlsInsertValue(newXlsObjNode, keys, value, offset+1)

}

func xlsFindKey(tree *XlsObjNode, keys []string, offset int) bool {
	if offset == len(keys) {
		return true
	}
	if tree.nodeType == ARRAY_TYPE || tree.nodeType == HASH_TYPE {
		for _, child := range tree.value.([]*XlsObjNode) {
			if child.key == keys[offset] || (strings.HasSuffix(keys[offset], "$") && keys[offset][:len(keys[offset])-1] == child.key) {
				return xlsFindKey(child, keys, offset+1)
			}
		}
	}
	return false
}

type BaseXlsParser struct {
	Parser
	app  *Application
	root *XlsObjNode
}

func (this *BaseXlsParser) initApp(app *Application) {
	this.app = app
}

func (this *BaseXlsParser) traversal(path string) error {
	xlFile, err := xlsx.OpenFile(path)
	if nil != err {
		return err
	}
	this.root = toXlsObject(xlFile)
	return nil
}

func (this *BaseXlsParser) Cleanup() {
	this.root = nil
}
