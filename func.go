package gqt

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"strings"
	"text/template"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/suifengpiao14/gqt/v2/gqttpl"
	"goa.design/goa/codegen"
	"golang.org/x/sync/singleflight"
)

func GetMD5LOWER(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Model2Entity copy model to entity ,some times input used to insert and update ,in this case input mybe model, copy model value to insertEntity and updateEntity
func Model2TplEntity(from interface{}, to gqttpl.TplEntityInterface) {
	err := copier.Copy(to, from)
	if err != nil {
		panic(err)
	}
}

var g = singleflight.Group{}

func Flight(sqlStr string, output interface{}, fn func() (interface{}, error)) (err error) {
	if sqlStr == "" {
		err = errors.New("sql must not be empty")
		return
	}
	value, err, _ := g.Do(GetMD5LOWER(sqlStr), fn)
	if err != nil {
		err = errors.WithStack(err)
	}
	rv := reflect.Indirect(reflect.ValueOf(output))
	if rv.CanSet() {
		valueRv := reflect.Indirect(reflect.ValueOf(value))
		rv.Set(valueRv)
	}
	return
}

// ConvertStruct 转换结构体
func ConvertStruct(from interface{}, to interface{}) {
	err := copier.Copy(to, from)
	if err != nil {
		panic(err)
	}
}

// 封装 goa.design/goa/v3/codegen 方便后续可定制
func ToCamel(name string) string {
	return codegen.CamelCase(name, true, true)
}

func ToLowerCamel(name string) string {
	return codegen.CamelCase(name, false, true)
}

func SnakeCase(name string) string {
	return codegen.SnakeCase(name)
}

// TplOutput 模板中执行模板，获取数据时使用 gqttool 生成的entity 会调用该方法，实现 TplEntityInterface 接口
func ExecTpl(dataVolume TplEntityInterface, fullname string) (output string, err error) {
	templateMapI, ok := dataVolume.GetValue(TEMPLATE_MAP_KEY)
	if !ok {
		err = errors.Errorf("not found key %s in %#v", TEMPLATE_MAP_KEY, fullname)
		return
	}
	var templateMap map[string]*template.Template
	templateMapRef, ok := templateMapI.(*map[string]*template.Template)
	if ok {
		templateMap = *templateMapRef
	} else {
		templateMap, ok = templateMapI.(map[string]*template.Template)
		if !ok {
			err = errors.Errorf(" key %s value want  %#v,got %#v", TEMPLATE_MAP_KEY, templateMap, fullname)
			return
		}
	}
	tplDefine, err := ExecuteTemplate(templateMap, fullname, dataVolume)
	if err != nil {
		return
	}
	output = tplDefine.Output
	return
}

func ToEOF(s string) string {
	out := strings.ReplaceAll(s, WINDOW_EOF, EOF) // 统一换行符
	return out
}
