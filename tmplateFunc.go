package gqt

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/suifengpiao14/gqt/v2/gqttpl"
)

var TemplatefuncMap = template.FuncMap{
	"zeroTime":      ZeroTime,
	"currentTime":   CurrentTime,
	"permanentTime": PermanentTime,
	"contains":      strings.Contains,
	"newPreComma":   NewPreComma,
	"in":            In,
	"toCamel":       gqttpl.ToCamel,
	"toLowerCamel":  gqttpl.ToLowerCamel,
	"snakeCase":     gqttpl.SnakeCase,
}

func SetMapData(dataVolume interface{}, k string, value interface{}) {
	t := reflect.TypeOf(dataVolume)
	if t.Kind() != reflect.Ptr {
		msg := fmt.Errorf("template data must be *map[string]interface{} when use tmplate func ,have %#v ,not address", dataVolume)
		panic(msg)
	}
	v := reflect.Indirect(reflect.ValueOf(dataVolume))
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Map {
		msg := fmt.Errorf("template data must be *map[string]interface{} when use tmplate func ,have %#v, not map", v.Kind())
		panic(msg)
	}
	v.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(value))

}
func GetMapData(dataVolume interface{}, k string) (value interface{}) {
	t := reflect.TypeOf(dataVolume)
	if t.Kind() != reflect.Ptr {
		msg := fmt.Errorf("template data must be *map[string]interface{} when use tmplate func ,have %#v ,not address", dataVolume)
		panic(msg)
	}
	v := reflect.Indirect(reflect.ValueOf(dataVolume))
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Map {
		msg := fmt.Errorf("template data must be *map[string]interface{} when use tmplate func ,have %#v, not map", v.Kind())
		panic(msg)
	}
	fv := v.MapIndex(reflect.ValueOf(k))
	if !fv.IsValid() {
		return nil
	}
	value = fv.Interface()

	return
}

func ZeroTime(dataVolume interface{}) string {
	named := "ZeroTime"
	placeholder := ":" + named
	value := "0000-00-00 00:00:00"
	SetMapData(dataVolume, named, value)
	return placeholder
}

func CurrentTime(dataVolume interface{}) string {
	named := "CurrentTime"
	placeholder := ":" + named
	value := time.Now().Format("2006-01-02 15:04:05")
	SetMapData(dataVolume, named, value)
	return placeholder
}

func PermanentTime(dataVolume interface{}) string {
	named := "PermanentTime"
	placeholder := ":" + named
	value := "3000-12-31 23:59:59"
	SetMapData(dataVolume, named, value)
	return placeholder
}

type preComma struct {
	comma string
}

func NewPreComma() *preComma {
	return &preComma{}
}

func (c *preComma) PreComma() string {
	out := c.comma
	c.comma = ","
	return out
}

func In(dataVolume interface{}, data interface{}) (str string, err error) {
	placeholders := make([]string, 0)
	inIndexKey := "InIndex_"
	inIndex := 0
	inIndexInterface := GetMapData(dataVolume, inIndexKey)
	if inIndexInterface != nil {
		inIndexInt, ok := inIndexInterface.(int)
		if ok {
			inIndex = inIndexInt
		}
	}

	v := reflect.Indirect(reflect.ValueOf(data))

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		num := v.Len()
		for i := 0; i < num; i++ {
			inIndex++
			named := fmt.Sprintf("in_%d", inIndex)
			placeholder := ":" + named
			placeholders = append(placeholders, placeholder)
			SetMapData(dataVolume, named, v.Index(i).Interface())
		}

	case reflect.String:
		arr := strings.Split(v.String(), ",")
		num := len(arr)
		for i := 0; i < num; i++ {
			inIndex++
			named := fmt.Sprintf("in_%d", inIndex)
			placeholder := ":" + named
			placeholders = append(placeholders, placeholder)
			SetMapData(dataVolume, named, arr[i])
		}
	default:
		err = fmt.Errorf("want slice/array/string ,have %s", v.Kind().String())
	}

	SetMapData(dataVolume, inIndexKey, inIndex) // 更新InIndex_

	if err != nil {
		return "", err
	}
	str = strings.Join(placeholders, ",")
	return str, nil

}
