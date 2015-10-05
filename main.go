package main

import (
	// "crypto/md5"
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
)

type Application struct {
	mapChan chan string

	group sync.WaitGroup

	inputFile string
	ouputFile string
	prefix    string

	inType  string
	outType string
}

func (this *Application) Run() error {

	for i := 0; i < 10; i++ {
		this.group.Add(1)
		go this.routine()
	}

	err := filepath.Walk(this.inputFile, this.fileCallback)
	if err != nil {
		log.Println(err)
		return nil
	}

	close(this.mapChan)

	this.group.Wait()
	return nil

}

func (this *Application) fileCallback(path string, info os.FileInfo, err error) error {
	if nil != err {
		if os.IsNotExist(err) {
			log.Println(err, path)
			return nil
		}
		log.Fatalln(path, err)
	}

	if info.IsDir() {
		return nil
	}

	ext := strings.ToLower(filepath.Ext(path))[1:]
	if ext != this.inType {
		return nil
	}

	this.mapChan <- path
	return nil

}

func (this *Application) routine() {

	var parser Parser
	if this.inType == "xml" && this.outType == "lua" {
		parser = &LuaParser{}
	} else if this.inType == "xml" && this.outType == "json" {
		parser = &JsonParser{}
	} else if this.inType == "xml" && this.outType == "xls" {
		parser = newExcelParser()
	} else if this.inType == "xls" && this.outType == "xml" {
		parser = &XlsToXmlParser{}
	} else if this.inType == "xls" && this.outType == "lua" {
		parser = &XlsToLuaParser{}
	} else if this.inType == "xls" && this.outType == "json" {
		// fun = this.xls2json
	}
	parser.initApp(this)

	for path := range this.mapChan {
		// log.Println("find path", path)
		// log.Println(fileConvertMd5(path))

		// file, err := os.Open(path)
		this.handler(path, parser)
		parser.Cleanup()
	}
	this.group.Done()
}

func (this *Application) handler(path string, parser Parser) {
	var err error
	// var file *os.File

	var name string
	var pkg string

	// file, err = os.Open(path)
	// if nil != err {
	// 	log.Fatalln("Open: ", err)
	// }
	// defer file.Close()

	name, err = filepath.Rel(this.inputFile, path)
	if nil != err {
		log.Fatalln(err)
	}

	pkg = strings.TrimSuffix(name, "."+this.inType)
	outputFileName := filepath.Join(this.ouputFile, pkg+"."+this.outType)

	err = os.MkdirAll(filepath.Dir(outputFileName), os.ModePerm)
	if nil != err {
		log.Fatalln("MkdirAll: ", err)
	}

	err = parser.traversal(path)

	if nil != err {
		log.Fatalln("parser error", path, err)
	}
	parser.SaveToFile(pkg, outputFileName)
}

// func (this *Application) xml2lua(pkg, outputFile string) error {

// 	luaParser := &LuaParser{}
// 	luaParser.app = this

// 	return nil
// }

// func (this *Application) xml2json(pkg, outputFile string) error {

// }

// func (this *Application) xml2xls(pkg, outputFile string) error {

// }

// func (this *Application) xls2lua(pkg, outputFile string) error {

// }

// func (this *Application) xls2json(pkg, outputFile string) error {

// }

// func (this *Application) xls2xml(pkg, outputFile string) error {

// }

func NewApplication(inputFile, ouputFile, inType, outType, prefix string) *Application {
	app := Application{}
	app.inputFile = inputFile
	app.ouputFile = ouputFile
	app.inType = inType
	app.outType = outType
	app.prefix = prefix

	app.group = sync.WaitGroup{}
	app.mapChan = make(chan string)

	return &app
}

func main() {

	pp, _ := os.Create("profile_file")
	pprof.StartCPUProfile(pp)
	defer pprof.StopCPUProfile()

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Lshortfile)

	input := flag.String("i", ".", "输入目录")
	output := flag.String("o", ".", "输出目录")

	inType := flag.String("it", "xml", "输入类型")
	outType := flag.String("ot", "xls", "输出类型")

	prefix := flag.String("p", "", "前置包名")

	flag.Parse()

	app := NewApplication(*input, *output, *inType, *outType, *prefix)

	app.Run()

}
