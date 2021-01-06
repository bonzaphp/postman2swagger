package main

import (
	"flag"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"postman2swagger/src/lib"
	"postman2swagger/src/lib/swagger-json"
)

var (
	inputFile = flag.String("source", "/Users/yons/Downloads/test01.json", "Postman 导出的json文件")
	//inputFile  = flag.String("source", "/Users/zabon/Desktop/code/元典接口/yang/crm-yang.20191231.json", "Postman 导出的json文件")
	outputFile = flag.String("output", "/Users/yons/code/swagger-php/result.json", "产生的注释文件")

	host        = flag.String("host", "192.168.2.199:8086", "项目地址,例如：ilessonpen.com")
	basePath    = flag.String("base_path", "/", "项目地址,例如 /")
	title       = flag.String("title", "校外生训v1.0.0", "校外生训")
	description = flag.String("description", "评测", "校外生训v1")
	version     = flag.String("version", "v1", "项目版本号")
	contact     = flag.String("contact", "bonzaphp@gmail.com", "联系方式")
)

func init() {
	flag.Parse()
	if *inputFile == "" {
		color.Red("缺少参数：source")
		os.Exit(0)
	}
	if *outputFile == "" {
		*outputFile = "result.php"
	}
}

func main() {
	file, err := os.OpenFile(*inputFile, os.O_RDONLY, os.ModePerm)
	lib.ErrorPut(err)
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
			// log etc
		}
	}()

	fileContent, err := ioutil.ReadAll(file)
	lib.FindRequest(string(fileContent), "")
	lib.ErrorPut(err)
	requestNum := len(lib.AllRequest)
	if requestNum == 0 {
		color.Red("没有找到request数据")
		os.Exit(0)
	}
	i := 0
	commentString := "{"
	comment := swagger_json.MakeTile(*host, *basePath, *version, *title, *description, *contact)
	for _, c := range comment {
		commentString = joinComment(commentString, c)
	}
	//commentString = joinComment(commentString, "")
	for i < requestNum {
		str := swagger_json.GeneratePaths(lib.AllRequest[i])
		//commentString = joinComment(commentString, "")
		for _, c := range str {
			commentString = joinComment(commentString, c)
		}

		//commentString = joinComment(commentString, "")
		//commentString = joinComment(commentString, "\n")
		i = i + 1
	}
	f, err := os.OpenFile(*outputFile, os.O_RDWR|os.O_CREATE, 0755)
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		color.Red(err.Error())
	} else {
		_, err := f.WriteString(commentString)
		if err != nil {
			panic("写入文件失败")
		}
		color.Green("Done!")
	}
}

func joinComment(source, newLine string) string {
	if source != "" {
		return source + "\n" + newLine
	} else {
		return newLine
	}
}
