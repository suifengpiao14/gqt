package main

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

var TemplatefuncMap = template.FuncMap{
	"contains": strings.Contains,
	"inIntSet": InIntSet,
	"inStrSet": InStrSet,
}

func InIntSet(data []int) (str string) {
	var formatData = make([]string, 0)
	for _, val := range data {
		valStr := strconv.Itoa(val)
		formatData = append(formatData, valStr)
	}
	str = strings.Join(formatData, ",")
	return
}

func InStrSet(data []string) (str string) {
	str = strings.Join(data, "\",\"")
	str = fmt.Sprintf("\"%s\"", str)
	return
}
