package lib

//字符数组转化为字符串
func ArrToString(str []string, sep string) string {
	tmpStr := ""
	for _, s := range str {
		if tmpStr == "" {
			tmpStr += s
		} else {
			tmpStr = tmpStr + sep + "\n" + s
		}
	}
	return tmpStr
}
