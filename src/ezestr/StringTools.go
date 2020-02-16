package ezestr

import (
	"regexp"
	"strconv"
	"strings"
)

func ReplaceAll(str string, old string, new string) string {
	if strings.Index(str, old) < 0 {
		return str
	} else {
		return ReplaceAll(ReplaceAll(str, old, new), old, new)
	}
}

// 这里可能有 bug，需要递归去除
func Remove(str string, removeStr []string) string {
	for _, v := range removeStr {
		str = strings.Replace(str, v, "", -1)
	}
	return str
}

func Continues(str string, continueStr string) bool {
	return strings.Contains(str, continueStr)
}

func StartWith(str string, startStr string) bool {
	return strings.Index(str, startStr) == 0
}

func GetAllStrByRegexp(str string, arrRegexp *regexp.Regexp) []string {
	return arrRegexp.FindAllString(str, -1)
}

func ConverToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

//func ConverToDouble()  {
//
//}
