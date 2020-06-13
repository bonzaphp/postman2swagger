package main

import (
	"flag"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"postman2swagger/src/lib"
)

var (
	inputFile = flag.String("source", "/Users/xxxx/Desktop/crm2.postman_collection.json", "Postman 导出的json文件")
	//inputFile  = flag.String("source", "/Users/zabon/Desktop/code/元典接口/yang/crm-yang.20191231.json", "Postman 导出的json文件")
	outputFile = flag.String("output", "/Users/xxxx/code/swagger-php/result.php", "产生的注释文件")

	host        = flag.String("host", "tpa.test", "项目地址,例如：api.xx.com")
	basePath    = flag.String("base_path", "/", "项目地址,例如 /")
	title       = flag.String("title", "xxxAPI", "项目名称")
	description = flag.String("description", "易点通crm", "项目描述")
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
	//defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	lib.FindRequest(string(fileContent), "")
	lib.ErrorPut(err)
	requestNum := len(lib.AllRequest)
	if requestNum == 0 {
		color.Red("没有找到request数据")
		os.Exit(0)
	}
	i := 0
	commentString := "<?php\n/**"
	comment := lib.MakeTile(*host, *basePath, *version, *title, *description, *contact)
	for _, c := range comment {
		commentString = joinComment(commentString, " *"+c)
	}
	commentString = joinComment(commentString, " */\n\n")
	/*	for _,v := range lib.AllRequest {
		fmt.Println(v.Path)
	}*/
	//panic(9)
	for i < requestNum {
		str := lib.MakeComment(lib.AllRequest[i])
		commentString = joinComment(commentString, "/**")
		for _, c := range str {
			commentString = joinComment(commentString, " *"+c)
		}

		commentString = joinComment(commentString, " */")
		commentString = joinComment(commentString, "\n\n")
		i = i + 1
	}

	//panic(3)
	f, err := os.OpenFile(*outputFile, os.O_RDWR|os.O_CREATE, 0755)
	defer func() {
		if err := f.Close(); err != nil {
			//log.Println(`file no found`)
			panic(err)
			// log etc
		}
	}()
	//defer f.Close()
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
