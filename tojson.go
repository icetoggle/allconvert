package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type JsonParser struct {
	BaseParser
}

func toJsonObject(obj interface{}) string {
	var result string = ""
	switch v := obj.(type) {
	case []interface{}:
		result += "["
		list := obj.([]interface{})
		prefix := "\n"
		for _, child := range list {
			result += prefix + toJsonObject(child)
			prefix = ",\n"
		}
		result += "\n]"
	case []*KeyValue:
		result += "{"
		set := obj.([]*KeyValue)
		prefix := "\n"
		for _, child := range set {
			result += prefix + "\"" + child.key + "\":" + toJsonObject(child.value)
			prefix = ",\n"
		}
		result += "\n}"
	case string:
		var s string = strings.TrimSpace(obj.(string))
		if s == "" {
			result = "null"
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
		result = "null"
		log.Printf("%T:\n%v\n", v, obj)
	}

	return result
}

func (p *JsonParser) SaveToFile(pkg, path string) {
	var result string = ""
	var obj interface{} = toObject(p.root)
	result += "{"
	switch obj.(type) {
	case []*KeyValue:
		set := obj.([]*KeyValue)
		prefix := "\n"
		for _, child := range set {
			result += prefix + "\"" + child.key + "\":" + toJsonObject(child.value)
			prefix = ",\n"
		}
	case string:
		result = "\"" + p.root.name + "\":" + obj.(string)
	}
	result += "}"
	file, err := os.Create(path)
	if nil != err {
		log.Fatalln(err)
	}
	defer file.Close()
	file.Write([]byte(result))
	//return toLuaObject(obj)
}
