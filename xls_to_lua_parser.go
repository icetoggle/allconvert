package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func xlsToLua(root *XlsObjNode, out *bytes.Buffer) {
	if root.nodeType == ARRAY_TYPE {
		valueList := root.value.([]*XlsObjNode)
		var flag = false

		out.WriteString(toKey(root.key) + "={\n")

		for _, child := range valueList {
			if !flag {
				flag = true
			} else {
				out.WriteString(",")
			}
			xlsToLua(child, out)

		}
		out.WriteString("}\n")
	} else if root.nodeType == HASH_TYPE {
		if root.key != "" {
			out.WriteString(toKey(root.key) + "=")
		}
		out.WriteString("{\n")
		valueList := root.value.([]*XlsObjNode)
		var flag = false
		for _, child := range valueList {
			if !flag {
				flag = true
			} else {
				out.WriteString(",")
			}
			xlsToLua(child, out)
		}
		out.WriteString("}\n")
	} else if root.nodeType == VALUE_TYPE || root.nodeType == ATTR_TYPE {
		if root.key != "" {
			out.WriteString(toKey(root.key) + "=")
		}
		out.WriteString(toValue(root.value.(string)) + "\n")
	}

}

type XlsToLuaParser struct {
	BaseXlsParser
}

func (this *XlsToLuaParser) SaveToFile(pkg, path string) {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("module(\"%s\", package.seeall)\n", strings.Replace(pkg, "/", ".", -1)))

	// out.WriteString("{\n")
	flag := false
	for _, child := range this.root.value.([]*XlsObjNode) {
		if !flag {
			flag = true
		} else {
			out.WriteString("\n")
		}

		xlsToLua(child, &out)
	}
	// out.WriteString("}\n")

	fout, _ := os.Create(path)
	s := out.String()

	fout.WriteString(s)
}
