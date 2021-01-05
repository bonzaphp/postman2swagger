package swagger_json

import (
	"github.com/tidwall/gjson"
	"postman2swagger/src/lib"
	"regexp"
	"strings"
)

type LineComment struct {
	Content   string
	IndentNum int
}

type AllComment []SingleComment

type SingleComment struct {
	Body []string
}

func MakeTile(host string, basePath string, version string, title string, description string, contact string) []string {
	comment := make([]string, 0)
	blankIndex := 0

	comment = append(comment, blankRepeat(blankIndex)+" \"swagger\": \"2.0\",")
	blankIndex = blankIndex + 1

	comment = append(comment, blankRepeat(blankIndex)+"schemes={\"http\",\"https\"},")

	comment = append(comment, blankRepeat(blankIndex)+"host=\""+host+"\",")
	comment = append(comment, blankRepeat(blankIndex)+"basePath=\""+basePath+"\",")
	comment = append(comment, blankRepeat(blankIndex)+"\"Info\": {")

	blankIndex = blankIndex + 1

	comment = append(comment, blankRepeat(blankIndex)+"version=\""+version+"\",")
	comment = append(comment, blankRepeat(blankIndex)+"title=\""+title+"\",")
	comment = append(comment, blankRepeat(blankIndex)+"description=\""+description+"\",")
	comment = append(comment, blankRepeat(blankIndex)+"termsOfService=\"\",")
	comment = append(comment, blankRepeat(blankIndex)+"\"Contact\": {")
	blankIndex = blankIndex + 1
	comment = append(comment, blankRepeat(blankIndex)+"email=\""+contact+"\"")

	blankIndex = blankIndex - 1
	comment = append(comment, blankRepeat(blankIndex)+"}")

	blankIndex = blankIndex - 1
	comment = append(comment, blankRepeat(blankIndex)+"}")

	return comment
}

func MakeComment(singeRequest lib.Request) []string {
	comment := make([]string, 0)
	blankIndex := 1
	comment = append(comment, blankRepeat(blankIndex)+"\"paths\": {")
	blankIndex = blankIndex + 1

	//path
	comment = append(comment, blankRepeat(blankIndex)+"\""+singeRequest.Path+"\": {")
	blankIndex = blankIndex + 1

	//请求方式
	comment = append(comment, "\""+singeRequest.Method+"\": {")
	blankIndex = blankIndex + 1

	//tags
	tags := strings.Split(singeRequest.Name, "/")[0]
	comment = append(comment, blankRepeat(blankIndex)+"\"tags\": [")
	blankIndex = blankIndex + 1

	comment = append(comment, blankRepeat(blankIndex)+tags)

	blankIndex = blankIndex - 1
	comment = append(comment, blankRepeat(blankIndex)+"]")

	//summary
	summary := strings.Replace(singeRequest.Name, "/", "-", -1)
	comment = append(comment, blankRepeat(blankIndex)+"\"summary\":\""+summary+"\"")

	//deprecated
	deprecated := strings.Contains(singeRequest.Name, "无效")
	if deprecated {
		comment = append(comment, blankRepeat(blankIndex)+"deprecated=\""+"true"+"\",")
		//fmt.Println(singeRequest.Name)
		//panic(6)
	}

	var singeParameter string

	r, _ := regexp.Compile(`{\w+}`)
	s := r.FindString(singeRequest.Path)
	//println(singeRequest.Path)
	if s != "" {
		var n string
		n = strings.Replace(s, "{", "", 1)
		n = strings.Replace(n, "}", "", 1)
		singeParameter = "@SWG\\Parameter(name =\"" + n + "\", type=\"" + "integer" + "\", required=true, in=\"path\",description=\"" + "路径参数" + "\",format=\"" + "int64" + "\"),"
	}

	comment = append(comment, blankRepeat(blankIndex)+singeParameter)

	//Parameter
	parameterIndex := 0
	queryNum := len(singeRequest.Query)

	for parameterIndex < queryNum {
		if singeRequest.Query[parameterIndex].Disabled == "true" {
			singeParameter = "@SWG\\Parameter(name =\"" + singeRequest.Query[parameterIndex].Key + "\", type=\"" + singeRequest.Query[parameterIndex].Type + "\", required=false, in=\"query\",description=\"" + singeRequest.Query[parameterIndex].Description + "\"),"
		} else {
			singeParameter = "@SWG\\Parameter(name =\"" + singeRequest.Query[parameterIndex].Key + "\", type=\"" + singeRequest.Query[parameterIndex].Type + "\", required=true, in=\"query\",description=\"" + singeRequest.Query[parameterIndex].Description + "\"),"
		}
		comment = append(comment, blankRepeat(blankIndex)+singeParameter)
		parameterIndex = parameterIndex + 1
	}

	//Body
	if singeRequest.Body.Mode == "raw" {
		comment = append(comment, blankRepeat(blankIndex)+"@SWG\\Schema(")

		bodyComment := make([]LineComment, 0)
		bodyComment = Comment(gjson.Parse(singeRequest.Body.Content.(string)), blankIndex, bodyComment)

		for _, singleBodyComment := range bodyComment {
			comment = append(comment, blankRepeat(blankIndex+singleBodyComment.IndentNum)+singleBodyComment.Content)
		}
		comment = append(comment, blankRepeat(blankIndex)+"),")

	} else if singeRequest.Body.Mode == "formdata" {
		for _, singleBodyParameter := range singeRequest.Body.Content.([]lib.Parameter) {
			//singeBodyParameter := "@SWG\\Parameter(name =\"" + singleBodyParameter.Key + "\", type=\"" + singleBodyParameter.Type + "\", required=true, in=\"body\",description=\"" + singleBodyParameter.Description + "\"),"
			singeBodyParameter := "@SWG\\Parameter(name =\"" + singleBodyParameter.Key + "\", type=\"" + singleBodyParameter.Type + "\", required=true, in=\"formdata\",description=\"" + singleBodyParameter.Description + "\"),"
			comment = append(comment, blankRepeat(blankIndex)+singeBodyParameter)
		}
	} else if singeRequest.Body.Mode == "urlencoded" {
		for _, singleBodyParameter := range singeRequest.Body.Content.([]lib.Parameter) {
			singeBodyParameter := "@SWG\\Parameter(name =\"" + singleBodyParameter.Key + "\", type=\"" + singleBodyParameter.Type + "\", required=true, in=\"formdata\",description=\"" + singleBodyParameter.Description + "\"),"
			comment = append(comment, blankRepeat(blankIndex)+singeBodyParameter)
		}
	}

	//Response
	comment = append(comment, blankRepeat(blankIndex)+"@SWG\\Response(")
	blankIndex = blankIndex + 1
	comment = append(comment, blankRepeat(blankIndex)+"response=\"200\",")
	comment = append(comment, blankRepeat(blankIndex)+"description=\"接口响应\",")
	comment = append(comment, blankRepeat(blankIndex)+"@SWG\\Schema(")

	responseComment := make([]LineComment, 0)
	responseComment = Comment(gjson.Parse(singeRequest.Response), blankIndex, responseComment)

	for _, singleResponse := range responseComment {
		comment = append(comment, blankRepeat(blankIndex+singleResponse.IndentNum)+singleResponse.Content)
	}
	comment = append(comment, blankRepeat(blankIndex)+")")
	blankIndex = blankIndex - 1
	comment = append(comment, blankRepeat(blankIndex)+")")
	blankIndex = blankIndex - 1
	comment = append(comment, blankRepeat(blankIndex)+")")
	return comment
}

func blankRepeat(num int) string {
	return strRepeat("  ", num)
}
func strRepeat(str string, num int) string {
	return strings.Repeat(str, num)
}
func Comment(json gjson.Result, level int, responseComment []LineComment) []LineComment {
	json.ForEach(func(key, value gjson.Result) bool {
		switch value.Type.String() {
		case "Number":
			line := LineComment{}

			thisValue := value.String()
			thisType := ""
			if strings.Contains(thisValue, ".") == true {
				thisType = "float"
			} else {
				thisType = "int"
			}
			line.Content = "@SWG\\Property( property=\"" + key.String() + "\" , type=\"" + thisType + "\" , example=\"" + thisValue + "\",description=\"填写描述\"),"
			line.IndentNum = level
			responseComment = append(responseComment, line)
			break

		case "String":
			line := LineComment{}
			thisValue := value.String()
			thisType := ""
			if thisValue == "true" || thisValue == "false" {
				thisType = "bool"
			} else {
				thisType = "string"
			}
			line.Content = "@SWG\\Property( property=\"" + key.String() + "\" , type=\"" + thisType + "\" , example=\"" + thisValue + "\",description=\"填写描述\"),"
			line.IndentNum = level
			responseComment = append(responseComment, line)
			break

		case "JSON":
			if value.IsArray() == true { //返回数据为空时，则设置为空字符串
				len := len(value.Array())
				if len == 0 {
					line := LineComment{}
					line.Content = "@SWG\\Property( property=\"data\" , type=\"string\" , example=\"\",description=\"填写描述\"),"
					line.IndentNum = level
					responseComment = append(responseComment, line)

				} else {
					value = value.Array()[0]

					lineStart := LineComment{}
					lineStart.Content = "@SWG\\Property( property=\"" + key.String() + "\" ,type=\"object\","
					//lineStart.Content = "@SWG\\Property( property=\"" + key.String() + "\" ,type=\"array\","
					lineStart.IndentNum = level
					responseComment = append(responseComment, lineStart)
					responseComment = Comment(value, level+1, responseComment)

					lineEnd := LineComment{}
					lineEnd.Content = "),"
					lineEnd.IndentNum = level
					responseComment = append(responseComment, lineEnd)
				}
			} else {
				lineStart := LineComment{}
				lineStart.Content = "@SWG\\Property( property=\"" + key.String() + "\" ,type=\"object\","
				lineStart.IndentNum = level
				responseComment = append(responseComment, lineStart)
				responseComment = Comment(value, level+1, responseComment)

				lineEnd := LineComment{}
				lineEnd.Content = "),"
				lineEnd.IndentNum = level
				responseComment = append(responseComment, lineEnd)
			}
		case "True", "False":
			//case "False":
			line := LineComment{}
			thisValue := value.String()
			thisType := "bool"

			line.Content = "@SWG\\Property( property=\"" + key.String() + "\" , type=\"" + thisType + "\" , example=\"" + thisValue + "\",description=\"填写描述\"),"
			line.IndentNum = level
			responseComment = append(responseComment, line)
			break
		default:
			//fmt.Println(value.Type.String())
			//fmt.Println(key.String())
		}

		return true
	})
	return responseComment
}
