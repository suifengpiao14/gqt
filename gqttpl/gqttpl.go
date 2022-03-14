package gqttpl

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"goa.design/goa/v3/codegen"
	"gorm.io/gorm/logger"
)

var TPlSuffix = ".tpl"
var SQLNamespaceSuffix = "sql"
var DDLNamespaceSuffix = "ddl"
var ConfigNamespaceSuffix = "config"
var MetaNamespaceSuffix = "meta"
var CURLNamespaceSuffix = "curl"

// RepositoryTemplate stores  templates.
type RepositoryTemplate struct {
	templates map[string]*template.Template // namespace: template
}

// NewRepositoryTemplate create a new RepositoryTemplate.
func NewRepositoryTemplate() *RepositoryTemplate {
	return &RepositoryTemplate{
		templates: make(map[string]*template.Template),
	}
}

func AddTemplateByDir(root string, namespaceSuffix string, funcMap template.FuncMap, leftDelim string, rightDelim string) (templateMap map[string]*template.Template, err error) {
	templateMap = make(map[string]*template.Template)
	// List the directories
	allFileList, err := GetTplFilesByDir(root, namespaceSuffix)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	for _, filename := range allFileList {
		namespace := FileName2Namespace(filename, root)
		t, err := template.New(namespace).Funcs(funcMap).Delims(leftDelim, rightDelim).ParseFiles(filename)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		templateMap[namespace] = t
	}
	return
}

func AddTemplateByStr(namespace string, content string, funcMap template.FuncMap, leftDelim string, rightDelim string) (t *template.Template, err error) {
	t, err = template.New(namespace).Funcs(funcMap).Delims(leftDelim, rightDelim).Parse(content)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return
}

type DefineResult struct {
	Name      string
	Tag       string // 标记类型 sql、ddl、sqlMeta、curl、curlMeta ...
	Namespace string
	Output    string
	Input     interface{}
}

func (d *DefineResult) Fullname() (fullname string) {
	fullname = fmt.Sprintf("%s.%s", d.Namespace, d.Name)
	return
}
func (d *DefineResult) FullnameCamel() (fullnameCamel string) {
	fullname := fmt.Sprintf("%s_%s", strings.ReplaceAll(d.Namespace, ".", "_"), d.Name)
	fullnameCamel = ToCamel(fullname)

	return
}

func SplitFullname(fullname string) (namespace string, name string) {
	lastIndex := strings.LastIndex(fullname, ".")
	if lastIndex < 0 {
		panic("illegal fullname ,want fullname format namespace.name")
	}
	namespace = fullname[:lastIndex]
	name = fullname[lastIndex+1:]
	return
}

// ExecuteNamespaceTemplate execute all template under namespace
func ExecuteNamespaceTemplate(templateMap map[string]*template.Template, namespace string, data interface{}) (defineResultList []*DefineResult, err error) {
	t, ok := templateMap[namespace]
	if !ok {
		err = errors.Errorf("not found namespace:%s", namespace)
		return
	}
	if err != nil {
		return nil, err
	}
	defineResultList = make([]*DefineResult, 0)
	templates := t.Templates()
	for _, tpl := range templates {
		defineResult, err := execTpl(tpl, namespace, data)
		if err != nil {
			return nil, err
		}
		defineResultList = append(defineResultList, defineResult)
	}
	return
}

func execTpl(tpl *template.Template, namespace string, data interface{}) (defineResult *DefineResult, err error) {
	var b bytes.Buffer
	err = tpl.Execute(&b, &data) // may adding more args with template func
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	out := b.String()
	defineResult = &DefineResult{
		Name:      tpl.Name(),
		Namespace: namespace,
		Output:    out,
		Input:     data,
	}
	return
}

// ExecuteNamespaceTemplate execute all template under namespace
func ExecuteTemplate(templateMap map[string]*template.Template, fullname string, data interface{}) (defineResult *DefineResult, err error) {
	namespace, name := SplitFullname(fullname)
	t, ok := templateMap[namespace]
	if !ok {
		err = errors.Errorf("not found namespace:%s", namespace)
		return nil, err
	}
	tpl := t.Lookup(name)
	if tpl == nil {
		err = errors.Errorf("ExecuteTemplate: no template %q associated with template %q", name, t.Name())
		return nil, err
	}
	defineResult, err = execTpl(tpl, namespace, data)
	if err != nil {
		return nil, err
	}

	return
}

// GetTplFilesByDir get current and reverse dir tpl file
func GetTplFilesByDir(dir string, namespaceSuffix string) (allFileList []string, err error) {
	pattern := fmt.Sprintf("%s/*%s%s", strings.TrimRight(dir, "/"), namespaceSuffix, TPlSuffix)
	allFileList, err = filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	pattern = fmt.Sprintf("%s/**/*%s%s", strings.TrimRight(dir, "/"), namespaceSuffix, TPlSuffix)
	subDirFileList, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	allFileList = append(allFileList, subDirFileList...)
	return
}

func FileName2Namespace(filename string, root string) (namespace string) {
	prefix := strings.ReplaceAll(root, "\\", ".")
	prefix = strings.ReplaceAll(prefix, "/", ".")
	namespace = strings.TrimSuffix(filename, TPlSuffix)
	namespace = strings.ReplaceAll(namespace, "\\", ".")
	namespace = strings.ReplaceAll(namespace, "/", ".")
	namespace = strings.TrimPrefix(namespace, prefix)
	namespace = strings.Trim(namespace, ".")
	return
}

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func Statement2SQL(sqlStatement string, vars []interface{}) (sqlStr string) {
	sqlStr = logger.ExplainSQL(sqlStatement, nil, `'`, vars...)
	return
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
