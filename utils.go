package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strconv"
)

func readerConvertMd5(r io.Reader) string {
	md5h := md5.New()
	io.Copy(md5h, r)
	return hex.EncodeToString(md5h.Sum(nil))
}

func fileConvertMd5(path string) string {
	file, err := os.Open(path)
	if nil != err {
		log.Fatalln(err)
	}
	defer file.Close()
	md5h := md5.New()
	io.Copy(md5h, file)
	return hex.EncodeToString(md5h.Sum(nil))
}

func bytesConvertMd5(bytes []byte) string {
	h := md5.New()
	h.Write(bytes)
	return hex.EncodeToString(h.Sum(nil))
}

func golog(fmt string, a ...interface{}) {
	log.Println(fmt, a)
}

func string2value(s string) interface{} {
	var result interface{} = nil
	if s == "" {
		return nil
	} else if r, err := strconv.ParseInt(s, 10, 64); err == nil {
		result = r
	} else if r, err := strconv.ParseFloat(s, 64); err == nil {
		result = r
	} else {
		result = s
	}
	return result
}

func toKey(key string) string {
	obj := string2value(key)
	switch obj.(type) {
	case string:
		return key
	case int64:
		return "[" + key + "]"
	case float64:
		return "[" + key + "]"
	}
	return ""
}

func toValue(value string) string {
	obj := string2value(value)
	switch obj.(type) {
	case string:
		return "\"" + value + "\""
	case int64, float64:
		return value
	}
	return ""
}
