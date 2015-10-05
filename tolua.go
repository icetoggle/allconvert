package main

import (
	// "encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func toLuaObject(obj interface{}) string {
	var result string = ""
	switch v := obj.(type) {
	case []interface{}:
		result += "{"
		list := obj.([]interface{})
		prefix := "\n"
		for _, child := range list {
			result += prefix + toLuaObject(child)
			prefix = ",\n"
		}
		result += "\n}"
	case []*KeyValue:
		result += "{"
		set := obj.([]*KeyValue)
		prefix := "\n"
		for _, child := range set {
			result += prefix + "[" + toLuaObject(child.key) + "]=" + toLuaObject(child.value)
			prefix = ",\n"
		}
		result += "\n}"
	case string:
		var s string = obj.(string)
		if s == "" {
			result = "nil"
		} else if _, err := strconv.ParseInt(s, 10, 64); err == nil {
			result = s
		} else if _, err := strconv.ParseFloat(s, 64); err == nil {
			result = s
		} else {
			result = "\"" + obj.(string) + "\""
		}
	case float64:
		var f float64 = obj.(float64)
		result = strconv.FormatFloat(f, byte('f'), -1, 64)
	case bool:
		var b bool = obj.(bool)
		result = strconv.FormatBool(b)
	default:
		result = "nil"
		log.Printf("%T:\n%v\n", v, obj)
	}

	return result
}

type LuaParser struct {
	BaseParser
}

func (p *LuaParser) SaveToFile(pkg, path string) {
	var result string = ""
	var obj interface{} = toObject(p.root)
	switch obj.(type) {
	case []*KeyValue:
		set := obj.([]*KeyValue)
		prefix := "\n"
		for _, child := range set {
			result += prefix + child.key + "=" + toLuaObject(child.value)
		}
	case string:
		result = p.root.name + "=" + obj.(string)
	}

	file, err := os.Create(path)
	if nil != err {
		log.Fatalln(err)
	}
	defer file.Close()
	pkg = strings.Replace(pkg, string(filepath.Separator), ".", -1)
	if p.app.prefix != "" {
		pkg = strings.Join([]string{p.app.prefix, pkg}, ".")
	}

	fmt.Fprintf(file, "module(\"%s\", package.seeall)\n", pkg)

	fileBytes := []byte(result)
	file.Write(fileBytes)
}
