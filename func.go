package gqt

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var TemplatefuncMap = template.FuncMap{
	"zeroTime":      ZeroTime,
	"currentTime":   CurrentTime,
	"permanentTime": PermanentTime,
	"contains":      strings.Contains,
	"inIntSet":      InIntSet,
	"inStrSet":      InStrSet,
}

func ZeroTime() string {
	return "0000-00-00 00:00:00"
}

func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func PermanentTime() string {
	return "3000-12-31 23:59:59"
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
