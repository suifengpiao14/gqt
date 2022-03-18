package gqttpl

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
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

// DataVolumeInterface 如果模板中需要新增、修改数据，作为 sql prepared statement 值，则传递给模板的数据变量必须实现DataVolumeInterface接口
type DataVolumeInterface interface {
	SetValue(key string, value interface{})
	GetValue(key string) (value interface{}, ok bool)
	GetDynamicValus() (values map[string]interface{})
}

// TplEntityInterface 模板参数对象，由于sql、curl经常需要在模板中增加数据，所以直接在模板输入实体接口融合DataVolumeInterface 接口功能，实体包含隐藏字段类型DataVolumeMap，即可实现DataVolumeInterface 功能
type TplEntityInterface interface {
	TplName() string
	SetValue(key string, value interface{})
	GetValue(key string) (value interface{}, ok bool)
	GetDynamicValus() (values map[string]interface{})
}

type DataVolumeMap map[string]interface{}

func (v *DataVolumeMap) init() {
	if *v == nil {
		*v = make(map[string]interface{})
	}
}

func (v DataVolumeMap) SetValue(key string, value interface{}) {
	v.init()
	(v)[key] = value // todo 并发lock
}

func (v DataVolumeMap) GetValue(key string) (value interface{}, ok bool) {
	v.init()
	value, ok = (v)[key]
	return
}

func (v DataVolumeMap) GetDynamicValus() (values map[string]interface{}) {
	v.init()
	return v
}

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

func AddTemplateByFS(fsys fs.FS, root string, namespaceSuffix string, funcMap template.FuncMap, leftDelim string, rightDelim string) (templateMap map[string]*template.Template, err error) {
	templateMap = make(map[string]*template.Template)
	// List the directories
	allFileList, err := GetTplFilesByFS(fsys, root, namespaceSuffix)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	for _, filename := range allFileList {
		namespace := FileName2Namespace(filename, root)
		b, err := fs.ReadFile(fsys, filename)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		content := string(b)
		t, err := AddTemplateByStr(namespace, content, funcMap, leftDelim, rightDelim)
		if err != nil {
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

type TPLDefine struct {
	Name      string
	Namespace string
	Output    string
	Input     interface{}
}

func (d *TPLDefine) Fullname() (fullname string) {
	fullname = fmt.Sprintf("%s.%s", d.Namespace, d.Name)
	return
}
func (d *TPLDefine) FullnameCamel() (fullnameCamel string) {
	fullname := fmt.Sprintf("%s_%s", strings.ReplaceAll(d.Namespace, ".", "_"), d.Name)
	fullnameCamel = ToCamel(fullname)

	return
}
func (d *TPLDefine) Tag() (tag string) {
	lastIndex := strings.Index(d.Namespace, ".")
	tag = d.Namespace
	if lastIndex > -1 {
		tag = d.Namespace[lastIndex+1:]
	}
	return
}

func SplitFullname(fullname string) (namespace string, name string) {
	lastIndex := strings.LastIndex(fullname, ".")
	if lastIndex < 0 {
		namespace = ""
		name = fullname
		return
	}
	namespace = fullname[:lastIndex]
	name = fullname[lastIndex+1:]
	return
}
func SpellFullname(namespace string, name string) (fullname string) {
	fullname = fmt.Sprintf("%s.%s", namespace, name)
	return
}

// ExecuteNamespaceTemplate execute all template under namespace
func ExecuteNamespaceTemplate(templateMap map[string]*template.Template, namespace string, data interface{}) (tplDefineList []*TPLDefine, err error) {
	t, ok := templateMap[namespace]
	if !ok {
		err = errors.Errorf("not found namespace:%s", namespace)
		return
	}
	if err != nil {
		return nil, err
	}
	tplDefineList = make([]*TPLDefine, 0)
	templates := t.Templates()
	for _, tpl := range templates {
		tplDefine, err := execTpl(tpl, namespace, data)
		if err != nil {
			return nil, err
		}
		tplDefineList = append(tplDefineList, tplDefine)
	}
	return
}

func execTpl(tpl *template.Template, namespace string, data interface{}) (tplDefine *TPLDefine, err error) {
	var b bytes.Buffer
	err = tpl.Execute(&b, &data) // may adding more args with template func
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	out := b.String()
	tplDefine = &TPLDefine{
		Name:      tpl.Name(),
		Namespace: namespace,
		Output:    out,
		Input:     data,
	}
	return
}

// ExecuteNamespaceTemplate execute all template under namespace
func ExecuteTemplate(templateMap map[string]*template.Template, fullname string, data interface{}) (tplDefine *TPLDefine, err error) {
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
	tplDefine, err = execTpl(tpl, namespace, data)
	if err != nil {
		return nil, err
	}
	return
}

//ExecuteTemplateTry 找不到模板的时候，返回null，不报错(curl 模板需要先执行xxxBody_模板)
func ExecuteTemplateTry(templateMap map[string]*template.Template, fullname string, data interface{}) (tplDefine *TPLDefine, err error) {
	namespace, name := SplitFullname(fullname)
	t, ok := templateMap[namespace]
	if !ok {
		err = errors.Errorf("not found namespace:%s", namespace)
		return nil, err
	}
	tpl := t.Lookup(name)
	if tpl == nil {
		return
	}
	tplDefine, err = execTpl(tpl, namespace, data)
	if err != nil {
		return nil, err
	}
	return
}

func GetTplFilesByFS(fsys fs.FS, dir string, namespaceSuffix string) (allFileList []string, err error) {
	dir = strings.TrimRight(dir, "/")
	pattern := fmt.Sprintf("%s/*%s%s", dir, namespaceSuffix, TPlSuffix)
	directFileList, err := fs.Glob(fsys, pattern)
	if err != nil {
		return nil, err
	}
	allFileList = append(allFileList, directFileList...)
	pattern = fmt.Sprintf("%s/**/*%s%s", dir, namespaceSuffix, TPlSuffix)
	subDirFileList, err := Glob(fsys, pattern)
	if err != nil {
		return nil, err
	}
	allFileList = append(allFileList, subDirFileList...)
	return
}

// GetTplFilesByDir get current and reverse dir tpl file
func GetTplFilesByDir(dir string, namespaceSuffix string) (allFileList []string, err error) {
	dir = strings.TrimRight(dir, "/")
	pattern := fmt.Sprintf("%s/*%s%s", dir, namespaceSuffix, TPlSuffix)
	directFileList, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	allFileList = append(allFileList, directFileList...)
	pattern = fmt.Sprintf("%s/**/*%s%s", dir, namespaceSuffix, TPlSuffix)
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

// Glob adds double-star support to the core path/filepath Glob function.
// It's useful when your globs might have double-stars, but you're not sure.
func Glob(fsys fs.FS, pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		// passthru to core package if no double-star
		return fs.Glob(fsys, pattern)
	}
	var matches []string
	regStr := strings.ReplaceAll(pattern, "**", ".*")
	reg := regexp.MustCompile(regStr)
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if reg.MatchString(path) {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches, nil
}

//Interface2DataVolume convert interface to DataVolumeInterface
func Interface2DataVolume(input interface{}) (out DataVolumeInterface, ok bool) {
	for {
		if inputI, ok := input.(*interface{}); ok {
			input = *inputI
		} else {
			break
		}
	}
	out, ok = input.(DataVolumeInterface)
	return
}
