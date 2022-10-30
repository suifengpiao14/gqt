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
	"newPreComma":   NewPreComma,
	"in":            In,
	"toCamel":       ToCamel,
	"toLowerCamel":  ToLowerCamel,
	"snakeCase":     SnakeCase,
	"tplOutput":     TplOutput,
}

func ZeroTime(tplEntity TplEntityInterface) (string, error) {
	named := "ZeroTime"
	placeholder := ":" + named
	value := "0000-00-00 00:00:00"
	tplEntity.SetValue(named, value)
	return placeholder, nil
}

func CurrentTime(tplEntity TplEntityInterface) (string, error) {
	named := "CurrentTime"
	placeholder := ":" + named
	value := time.Now().Format("2006-01-02 15:04:05")
	tplEntity.SetValue(named, value)
	return placeholder, nil
}

func PermanentTime(tplEntity TplEntityInterface) (string, error) {
	named := "PermanentTime"
	placeholder := ":" + named
	value := "3000-12-31 23:59:59"
	tplEntity.SetValue(named, value)
	return placeholder, nil
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

func In(tplEntity TplEntityInterface, data interface{}) (str string, err error) {
	placeholders := make([]string, 0)
	inIndexKey := "InIndex_"
	inIndex := 0
	inIndexInterface, _ := tplEntity.GetValue(inIndexKey)
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
			tplEntity.SetValue(named, v.Index(i).Interface())
		}

	case reflect.String:
		arr := strings.Split(v.String(), ",")
		num := len(arr)
		for i := 0; i < num; i++ {
			inIndex++
			named := fmt.Sprintf("in_%d", inIndex)
			placeholder := ":" + named
			placeholders = append(placeholders, placeholder)
			tplEntity.SetValue(named, arr[i])
		}
	default:
		err = fmt.Errorf("want slice/array/string ,have %s", v.Kind().String())
		if err != nil {
			return "", err
		}
	}
	tplEntity.SetValue(inIndexKey, inIndex) // 更新InIndex_
	str = strings.Join(placeholders, ",")
	return str, nil

}

func TplOutput(dataVolume TplEntityInterface, tplEntity TplEntityInterface) (output string, err error) {
	return ExecTpl(dataVolume, tplEntity.TplName())
}
