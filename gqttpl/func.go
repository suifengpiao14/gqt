package gqttpl

import (
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"goa.design/goa/v3/codegen"
)

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

//TplOutput 模板中执行模板，获取数据时使用 gqttool 生成的entity 会调用该方法，实现 TplEntityInterface 接口
func TplOutput(dataVolume TplEntityInterface, tplEntity TplEntityInterface) (output string, err error) {
	templateMapI, ok := dataVolume.GetValue(TEMPLATE_MAP_KEY)
	if !ok {
		err = errors.Errorf("not found key %s in %#v", TEMPLATE_MAP_KEY, tplEntity)
		return
	}
	var templateMap map[string]*template.Template
	templateMapRef, ok := templateMapI.(*map[string]*template.Template)
	if ok {
		templateMap = *templateMapRef
	} else {
		templateMap, ok = templateMapI.(map[string]*template.Template)
		if !ok {
			err = errors.Errorf(" key %s value want  %#v,got %#v", TEMPLATE_MAP_KEY, templateMap, tplEntity)
			return
		}
	}
	tplDefine, err := ExecuteTemplate(templateMap, tplEntity.TplName(), dataVolume)
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
