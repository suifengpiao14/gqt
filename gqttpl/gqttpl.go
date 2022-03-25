package gqttpl

import (
	"bytes"
	"fmt"
	"io/fs"
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

const TEMPLATE_MAP_KEY = "_templateMap"

// TplEntityInterface 模板参数对象，由于sql、curl经常需要在模板中增加数据，所以直接在模板输入实体接口融合TplEntityInterface 接口功能，实体包含隐藏字段类型tplEntityMap，即可实现TplEntityInterface 功能
type TplEntityInterface interface {
	TplName() string
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

type TPLDefineList []*TPLDefine

const (
	TPL_DEFINE_TYPE_CURL_REQUEST  = "curl_request"
	TPL_DEFINE_TYPE_CURL_RESPONSE = "curl_response"
	TPL_DEFINE_TYPE_SQL           = "sql"
	TPL_DEFINE_TYPE_TEXT          = "text"
)

type TPLDefine struct {
	Name      string
	Namespace string
	Output    string
	Type      string
	Input     TplEntityInterface
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

// 判断一个(变量)名词是否和define 名称相同
func (dl TPLDefineList) IsDefineNameCamel(variableName string) bool {
	for _, TPLDefine := range dl {
		if ToCamel(TPLDefine.Name) == variableName {
			return true
		}
	}
	return false
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

//Interface2tplEntity convert interface to TplEntityInterface 核心思路：使得 input 和 out 指向同一个内存地址
func Interface2tplEntity(input interface{}) (out TplEntityInterface, ok bool) {
	if inputI, ok := input.(*interface{}); ok {
		input = *inputI
	}
	if tplEntity, ok := input.(TplEmptyEntity); ok {
		out = &tplEntity
		return out, ok
	}
	if tplEntityMapRef, ok := input.(*TplEmptyEntity); ok {
		out = tplEntityMapRef
		return out, ok
	}

	if inputMap, ok := input.(map[string]interface{}); ok {
		inputvolumeMap := TplEmptyEntity(inputMap)
		out = &inputvolumeMap
		return out, ok
	}
	if inputMap, ok := input.(*map[string]interface{}); ok { // 同时更新input 内的对象，使得input、out指向同一个地址 data21
		tmp := TplEmptyEntity(*inputMap)
		out = &tmp
		return out, ok
	}

	if out, ok := input.(TplEntityInterface); ok {
		v := reflect.Indirect(reflect.ValueOf(out))
		t := v.Type()
		if t.Kind() == reflect.Struct {
			defaulttplEntityMap := &TplEmptyEntity{}
			targetType := reflect.TypeOf(defaulttplEntityMap)
			n := t.NumField()
			for i := 0; i < n; i++ {
				fv := v.Field(i)
				ft := fv.Type()
				if ft == targetType && fv.IsValid() && fv.IsNil() {
					if fv.CanSet() {
						fv.Set(reflect.ValueOf((defaulttplEntityMap))) //解决结构体无名称方式引用 *tplEntityMap ,实例化时，并没有实力该字段，导致地址为空,解决测试用例data34 填充值问题
					} else {
						// todo resolve data32 panic
					}
				}
			}
		}
		return out, ok
	}
	return
}
