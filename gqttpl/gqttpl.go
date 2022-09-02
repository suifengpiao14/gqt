package gqttpl

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"gorm.io/gorm/logger"
)

var TPlSuffix = ".tpl"
var SQLNamespaceSuffix = "sql"
var DDLNamespaceSuffix = "ddl"
var ConfigNamespaceSuffix = "config"
var MetaNamespaceSuffix = "meta"
var CURLNamespaceSuffix = "curl"
var LeftDelim = "{{"
var RightDelim = "}}"

const TEMPLATE_MAP_KEY = "_templateMap"
const URI_KEY = "__URI__" // 记录资源地址key(curl 请求地址、db 连接地址等)

const (
	EOF                  = "\n"
	WINDOW_EOF           = "\r\n"
	HTTP_HEAD_BODY_DELIM = EOF + EOF
)

type TPLDefine struct {
	Name      string
	Namespace string
	Output    string
	Input     interface{}
}

// TplEntityInterface 模板参数对象，由于sql、curl经常需要在模板中增加数据，所以直接在模板输入实体接口融合TplEntityInterface 接口功能，实体包含隐藏字段类型tplEntityMap，即可实现TplEntityInterface 功能
type TplEntityInterface interface {
	TplName() string
	TplType() string // 返回 TPL_DEFINE_TYPE 类型，方便后续根据类型获取资源(db、curl) 自动获取数据
	SetValue(key string, value interface{})
	GetValue(key string) (value interface{}, ok bool)
	GetDynamicValus() (values map[string]interface{})
}

type TplEmptyEntity map[string]interface{}

func (v *TplEmptyEntity) init() {
	if v == nil {
		err := errors.Errorf("*TplEmptyEntity must init")
		panic(err)
	}
	if *v == nil {
		*v = TplEmptyEntity{} // 解决 data33 情况
	}
}

func (v *TplEmptyEntity) SetValue(key string, value interface{}) {
	v.init()
	(*v)[key] = value // todo 并发lock
}

func (v *TplEmptyEntity) GetValue(key string) (value interface{}, ok bool) {
	v.init()
	value, ok = (*v)[key]
	return
}

func (v *TplEmptyEntity) GetDynamicValus() (values map[string]interface{}) {
	v.init()
	return *v
}

func (v *TplEmptyEntity) TplName() string {
	err := errors.Errorf("*TplEmptyEntity.TplName is empty")
	panic(err)
}
func (v *TplEmptyEntity) TplType() string {
	err := errors.Errorf("*TplEmptyEntity.TplType is empty")
	panic(err)
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
func ExecuteNamespaceTemplate(templateMap map[string]*template.Template, namespace string, tplEntity TplEntityInterface) (tplDefineList []*TPLDefine, err error) {
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
	if tplEntity == nil { // 确保数据容器对象不为空
		tplEntity = &TplEmptyEntity{}
	}
	tplEntity.SetValue(TEMPLATE_MAP_KEY, templateMap) // 将模板传入，方便在模板中执行模板
	for _, tpl := range templates {
		tplDefine, err := execTpl(tpl, namespace, tplEntity)
		if err != nil {
			return nil, err
		}

		if TrimSpaces(tplDefine.Output) == "" {
			continue // skip default empty template define
		}
		tplDefineList = append(tplDefineList, tplDefine)
	}
	return
}

func execTpl(tpl *template.Template, namespace string, tplEntity TplEntityInterface) (tplDefine *TPLDefine, err error) {
	var b bytes.Buffer
	tplEntity.SetValue("tpl", tpl)
	err = tpl.Execute(&b, &tplEntity) // 此处使用引用地址，方便在模板中增加数据，返回到data中
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	out := b.String()
	tplDefine = &TPLDefine{
		Name:      tpl.Name(),
		Namespace: namespace,
		Output:    out,
		Input:     tplEntity,
	}
	return
}

// ExecuteNamespaceTemplate execute all template under namespace
func ExecuteTemplate(templateMap map[string]*template.Template, fullname string, tplEntity TplEntityInterface) (tplDefine *TPLDefine, err error) {
	if tplEntity == nil {
		err = errors.Errorf("ExecuteTemplate tplEntity must not nil ")
		return nil, err
	}
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
	tplEntityR := reflect.ValueOf(tplEntity)
	if tplEntityR.IsNil() {
		err := errors.Errorf("%#v must not nil", tplEntity)
		return nil, err

	}
	tplEntity.SetValue(TEMPLATE_MAP_KEY, templateMap) // 将模板传入，方便在模板中执行模板
	tplDefine, err = execTpl(tpl, namespace, tplEntity)
	if err != nil {
		return nil, err
	}
	return
}

//ExecuteTemplateTry 找不到模板的时候，返回null，不报错(curl 模板需要先执行xxxBody_模板)
func ExecuteTemplateTry(templateMap map[string]*template.Template, fullname string, tplEntity TplEntityInterface) (tplDefine *TPLDefine, err error) {
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
	if tplEntity == nil { // 确保数据容器对象不为空
		tplEntity = &TplEmptyEntity{}
	}
	tplEntity.SetValue(TEMPLATE_MAP_KEY, templateMap) // 将模板传入，方便在模板中执行模板
	tplDefine, err = execTpl(tpl, namespace, tplEntity)
	if err != nil {
		return nil, err
	}
	return
}

func GetTplFilesByFS(fsys fs.FS, dir string, namespaceSuffix string) (allFileList []string, err error) {
	pattern := fmt.Sprintf("%s/**%s%s", dir, namespaceSuffix, TPlSuffix)
	return Glob(fsys, pattern)
}

// GetTplFilesByDir get current and reverse dir tpl file
func GetTplFilesByDir(dir string, namespaceSuffix string) (allFileList []string, err error) {
	pattern := fmt.Sprintf("%s/**%s%s", strings.TrimRight(dir, "/"), namespaceSuffix, TPlSuffix)
	return GlobDirectory(dir, pattern)
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

//TrimSpaces  去除开头结尾的非有效字符
func TrimSpaces(s string) string {
	return strings.Trim(s, "\r\n\t\v\f ")
}

func Statement2SQL(sqlStatement string, vars []interface{}) (sqlStr string) {
	sqlStr = logger.ExplainSQL(sqlStatement, nil, `'`, vars...)
	return
}

// Glob adds double-star support to the core path/filepath Glob function.
// It's useful when your globs might have double-stars, but you're not sure.
func Glob(fsys fs.FS, pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		// passthru to core package if no double-star
		return fs.Glob(fsys, pattern)
	}
	var matches []string
	regStr := strings.ReplaceAll(pattern, ".", "\\.")
	regStr = strings.ReplaceAll(regStr, "**", ".*")
	reg := regexp.MustCompile(regStr)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if reg.MatchString(path) {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches, err
}

func GlobDirectory(dir string, pattern string) ([]string, error) {
	dir = strings.TrimRight(dir, "/")
	if !strings.Contains(pattern, "**") {
		pattern = fmt.Sprintf("%s/*%s", dir, pattern)
		// passthru to core package if no double-star
		return filepath.Glob(pattern)
	}
	var matches []string
	regStr := strings.ReplaceAll(pattern, "\\", "/")
	regStr = strings.ReplaceAll(regStr, ".", "\\.")
	regStr = strings.ReplaceAll(regStr, "**", ".*")
	reg := regexp.MustCompile(regStr)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			err := errors.Errorf("dir:%s filepath.Walk info is nil", dir)
			return err
		}
		if !info.IsDir() {
			path = strings.ReplaceAll(path, "\\", "/")
			if reg.MatchString(path) {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches, err
}
