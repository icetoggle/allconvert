package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type JsonFile struct {
	data map[string]interface{}
}

func (this *JsonFile) load(path string) {
	var file *os.File
	var err error

	file, err = os.OpenFile(path+"/"+"data.md5", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("Read md5 data err: ", err)
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)

	var result map[string]interface{}

	if len(fileBytes) == 0 {
		result = make(map[string]interface{})
	} else if err = json.Unmarshal(fileBytes, &result); err != nil {
		log.Fatalln(err)
	}

	this.data = result
}

func (this *JsonFile) isChangeAndSave(path string) bool {

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	paths := strings.Split(path, string(filepath.Separator))

	var childPath interface{}
	var ok bool

	root := this.data

	pathLen := len(paths)
	for i := 0; i < pathLen-1; i++ {
		childPath, ok = root[paths[i]]
		if !ok {
			childPath = make(map[string]interface{})
			root[paths[i]] = childPath
		}
		root = childPath.(map[string]interface{})
	}

	md5key, ok := root[paths[pathLen-1]]

	newMd5Key := readerConvertMd5(file)

	if !ok || md5key != newMd5Key {
		root[paths[pathLen-1]] = newMd5Key
		return true
	}

	return false

}

func (this *JsonFile) saveMd5(path string) {
	paths := strings.Split(path, string(filepath.Separator))
	var childPath interface{}
	var ok bool
	root := this.data

	pathLen := len(paths)
	for i := 0; i < pathLen-1; i++ {
		childPath, ok = root[paths[i]]
		if !ok {
			childPath = make(map[string]interface{})
			root[paths[i]] = childPath
		}
		root = childPath.(map[string]interface{})
	}

	root[paths[pathLen-1]] = fileConvertMd5(path)

}

func (this *JsonFile) save(path string) {
	var err error
	var fileBytes []byte
	var file *os.File

	fileBytes, err = json.Marshal(this.data)
	file, err = os.Create(path + "/data.md5")
	if nil != err {
		log.Fatalln(err)
	}
	defer file.Close()
	file.Write(fileBytes)
}
