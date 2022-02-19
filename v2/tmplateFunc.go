package gqt

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"
)

var TemplatefuncMap = template.FuncMap{
	"zeroTime":      ZeroTime,
	"currentTime":   CurrentTime,
	"permanentTime": PermanentTime,
	"contains":      strings.Contains,
	"in":            In,
}

func ZeroTime(dataVolume *map[string]interface{}) string {
	named := "ZeroTime"
	placeholder := ":" + named
	(*dataVolume)[named] = "0000-00-00 00:00:00"
	return placeholder
}

func CurrentTime(dataVolume *map[string]interface{}) string {
	named := "CurrentTime"
	placeholder := ":" + named
	(*dataVolume)[named] = time.Now().Format("2006-01-02 15:04:05")
	return placeholder
}

func PermanentTime(dataVolume *map[string]interface{}) string {
	named := "PermanentTime"
	placeholder := ":" + named
	(*dataVolume)[named] = "3000-12-31 23:59:59"
	return placeholder
}

func In(dataVolume *map[string]interface{}, key string) (str string, err error) {
	placeholders := make([]string, 0)
	data, ok := (*dataVolume)[key]
	if !ok {
		err = fmt.Errorf("not found %s in %#v", key, dataVolume)
		return "", err
	}
	v := reflect.Indirect(reflect.ValueOf(data))

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		num := v.Len()
		for i := 0; i < num; i++ {
			named := fmt.Sprintf("%s_%d", key, i)
			placeholder := ":" + named
			placeholders = append(placeholders, placeholder)
			(*dataVolume)[named] = v.Index(i).Interface()
		}

	case reflect.String:
		arr := strings.Split(v.String(), ",")
		num := len(arr)
		for i := 0; i < num; i++ {
			named := fmt.Sprintf("%s_%d", key, i)
			placeholder := ":" + named
			placeholders = append(placeholders, placeholder)
			(*dataVolume)[named] = arr[i]
		}
	default:
		err = fmt.Errorf("want slice/array/string ,have %s", v.Kind().String())
	}

	if err != nil {
		return "", err
	}
	str = strings.Join(placeholders, ",")
	return str, nil

}
