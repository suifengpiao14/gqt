package gqt

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/suifengpiao14/gqt/v2/pkg"

	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

// Repository stores SQL templates.
type Repository struct {
	templates map[string]*template.Template // namespace: template
}

type DataVolume struct {
	Data  interface{}
	Extra *map[string]interface{}
}

// NewRepository creates a new Repository.
func NewRepository() *Repository {
	return &Repository{
		templates: make(map[string]*template.Template),
	}
}

var SQLSuffix = ".sql.tpl"
var DDLSuffix = ".ddl.tpl"

var MetaTplFlag = "metaTpl"
var LeftDelim = "{{"
var RightDelim = "}}"

// ddl namespace sufix . define name prefix
var DDLNamespaceSuffix = "ddl"

func (r *Repository) AddByDir(root string, funcMap template.FuncMap) (err error) {
	r.templates, err = pkg.AddTemplateByDir(root, SQLSuffix, funcMap, LeftDelim, RightDelim)
	if err != nil {
		return
	}
	ddlTemplates, err := pkg.AddTemplateByDir(root, DDLSuffix, funcMap, LeftDelim, RightDelim)
	if err != nil {
		return
	}
	for fullname, tpl := range ddlTemplates {
		r.templates[fullname] = tpl
	}
	return
}

func (r *Repository) AddByNamespace(namespace string, content string, funcMap template.FuncMap) (err error) {
	t, err := pkg.AddTemplateByStr(namespace, content, funcMap, LeftDelim, RightDelim)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	r.templates[namespace] = t
	return
}

type SQLRow struct {
	Name      string
	Namespace string
	SQL       string
	Statment  string
	Arguments []interface{}
	Result    interface{}
}

func (r *Repository) DefineResult2SQLRow(defineResult pkg.DefineResult) (sqlRow *SQLRow, err error) {
	sqlRow = &SQLRow{
		Name:      defineResult.Name,
		Namespace: defineResult.Namespace,
	}

	sqlNamed := pkg.StandardizeSpaces(defineResult.Output)
	if sqlNamed == "" {
		return
	}
	sqlRow.Statment, sqlRow.Arguments, err = sqlx.Named(sqlNamed, defineResult.Input)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	sqlRow.SQL = pkg.Statement2SQL(sqlRow.Statment, sqlRow.Arguments)

	return
}

// GetByNamespace get all template under namespace
func (r *Repository) GetByNamespace(namespace string, data interface{}) (sqlRowList []*SQLRow, err error) {
	data, err = interface2map(data)
	if err != nil {
		return nil, err
	}
	defineResultList, err := pkg.ExecuteNamespaceTemplate(r.templates, namespace, data)
	if err != nil {
		return nil, err
	}
	sqlRowList = make([]*SQLRow, 0)
	for _, defineResult := range defineResultList {
		sqlRow, err := r.DefineResult2SQLRow(*defineResult)
		if err != nil {
			return nil, err
		}
		sqlRowList = append(sqlRowList, sqlRow)
	}
	return
}

func (r *Repository) GetDDLNamespace() (ddlNamespace string, err error) {
	for namespace := range r.templates {
		if strings.HasSuffix(namespace, DDLNamespaceSuffix) {
			ddlNamespace = namespace
			break
		}
	}
	if ddlNamespace == "" {
		err = errors.Errorf("not find ddl namespace with sufix %s,you can set gqt.DDLNamespaceSuffix to change sufix", DDLNamespaceSuffix)
		return
	}
	return
}

func (r *Repository) GetDDLSQL() (ddlSQLRowList []*SQLRow, err error) {
	ddlSQLRowList = make([]*SQLRow, 0)
	ddlNamespace, err := r.GetDDLNamespace()
	if err != nil {
		return
	}
	sqlRowList, err := r.GetByNamespace(ddlNamespace, nil)
	if err != nil {
		return
	}
	for _, sqlRow := range sqlRowList {
		if len(sqlRow.SQL) < 6 {
			continue
		}
		createStr := sqlRow.SQL[:6]
		if strings.ToLower(createStr) == "create" {
			ddlSQLRowList = append(ddlSQLRowList, sqlRow)
		}
	}
	return
}

// 支持返回Prepared Statement ,该模式优势1. 提升性能，避免重复解析 SQL 带来的开销，2. 避免 SQL 注入 缺点： 1. 存在两次与数据库的通信，在密集进行 SQL 查询的情况下，可能会出现 I/O 瓶颈
func (r *Repository) GetStatement(name string, data interface{}) (sqlStatement string, vars []interface{}, err error) {
	if name == "" {
		err = errors.New("name not be empty")
		return "", nil, err
	}
	mapData, err := interface2map(data)
	if err != nil {
		return "", nil, err
	}
	sqlNamed, err := r.Parse(name, &mapData) // 当data为map[string]interface{}时，模板内可以改变data值
	if err != nil {
		return "", nil, err
	}
	sqlStatement, vars, err = sqlx.Named(sqlNamed, mapData)
	sqlStatement = strings.ReplaceAll(sqlStatement, "\r", "")
	sqlStatement = strings.ReplaceAll(sqlStatement, "\n", "")
	sqlStatement = strings.ReplaceAll(sqlStatement, "  ", " ")
	sqlStatement = strings.Trim(sqlStatement, " ")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

type TplEntity interface {
	TplName() string
}

// 将模板名称，模板中的变量，封装到结构体中，使用结构体访问，避免拼写错误以及分散的硬编码，可以配合 gqttool 自动生成响应的结构体
func (r *Repository) GetSQLByTplEntity(t TplEntity) (sqlRow *SQLRow, err error) {
	return r.GetSQL(t.TplName(), t)
}

// GetSQLByTplEntityRef 支持只返回error 函数签名
func (r *Repository) GetSQLRawByTplEntityRef(t TplEntity, sqlStr *string) (err error) {
	sqlRow, err := r.GetSQL(t.TplName(), t)
	*sqlStr = sqlRow.SQL
	return
}

//无sql注入的安全方式
func (r *Repository) GetSQL(fullname string, data interface{}) (sqlRow *SQLRow, err error) {
	defineResult, err := pkg.ExecuteTemplate(r.templates, fullname, data)
	if err != nil {
		return nil, err
	}
	sqlRow, err = r.DefineResult2SQLRow(*defineResult)
	return
}

// Parse executes the template and returns the resulting SQL or an error.
func (r *Repository) Parse(name string, data interface{}) (string, error) {
	// Prepare namespace and block name
	if name == "" {
		return "", errors.Errorf("unnamed block")
	}
	path := strings.Split(name, ".")
	namespace := strings.Join(path[0:len(path)-1], ".")
	if namespace == "." {
		namespace = ""
	}
	block := path[len(path)-1]

	// Execute the template
	var b bytes.Buffer
	t, ok := r.templates[namespace]
	if !ok {
		return "", errors.Errorf("unknown namespace \"%s\"", namespace)
	}
	err := t.ExecuteTemplate(&b, block, data)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return b.String(), nil
}

func (r *Repository) NewSQLChain() *SQLChain {
	return &SQLChain{
		sqlRows:       make([]*SQLRow, 0),
		sqlRepository: func() *Repository { return r },
	}
}

type SQLChain struct {
	sqlRows       []*SQLRow
	sqlRepository func() *Repository
	err           error
}

func (s *SQLChain) ParseSQL(tplName string, args interface{}, result interface{}) *SQLChain {
	if s.sqlRepository == nil {
		s.err = errors.Errorf("want SQLChain.sqlRepository ,have %#v", s)
	}
	if s.err != nil {
		return s
	}
	sqlRow, err := s.sqlRepository().GetSQL(tplName, args)
	if err != nil {
		s.err = err
		return s
	}
	s.sqlRows = append(s.sqlRows, sqlRow)
	return s
}

func (s *SQLChain) ParseTpEntity(entity TplEntity, result interface{}) *SQLChain {
	if s.sqlRepository == nil {
		s.err = errors.Errorf("want SQLChain.sqlRepository ,have %#v", s)
	}
	if s.err != nil {
		return s
	}
	sqlRow, err := s.sqlRepository().GetSQLByTplEntity(entity)
	if err != nil {
		s.err = err
		return s
	}
	s.sqlRows = append(s.sqlRows, sqlRow)
	return s
}

//GetAllSQL get all sql from SQLChain
func (s *SQLChain) SQLRows() (sqlRowList []*SQLRow, err error) {
	return s.sqlRows, s.err
}

//Exec exec sql
func (s *SQLChain) Exec(fn func(sqlRowList []*SQLRow) (e error)) (err error) {
	if s.err != nil {
		return s.err
	}
	s.err = fn(s.sqlRows)
	return s.err
}

//Exec exec sql ,get data
func (s *SQLChain) Scan(fn func(sqlRowList []*SQLRow) (e error)) (err error) {
	if s.err != nil {
		return
	}
	s.err = fn(s.sqlRows)
	return s.err
}

//AddSQL add one sql to SQLChain
func (s *SQLChain) AddSQL(namespace string, name string, sql string, result interface{}) {
	sqlRow := &SQLRow{
		Name:      name,
		Namespace: name,
		SQL:       sql,
		Result:    result,
	}
	s.sqlRows = append(s.sqlRows, sqlRow)
}

func (s *SQLChain) SetError(err error) {
	if s.err != nil {
		return
	}
	if err != nil {
		err = errors.WithStack(err)
		s.err = err
	}
}

func (s *SQLChain) Error() (err error) {
	return s.err
}

// 批量获取sql记录
func NewSQLChain(sqlRepository func() *Repository) (s *SQLChain) {
	s = &SQLChain{
		sqlRows:       make([]*SQLRow, 0),
		sqlRepository: sqlRepository,
	}
	return
}

var g = singleflight.Group{}

func Flight(sqlStr string, fn func() (interface{}, error)) (err error) {
	if sqlStr == "" {
		err = errors.New("sql must not be empty")
		return
	}
	_, err, _ = g.Do(GetMD5LOWER(sqlStr), fn)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

func interface2map(data interface{}) (out map[string]interface{}, err error) {
	out = make(map[string]interface{}, 0)
	if data == nil {
		return
	}
	v := reflect.Indirect(reflect.ValueOf(data))
	switch v.Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		for _, key := range keys {
			v := v.MapIndex(key).Interface()
			out[key.String()] = v
		}
	case reflect.Struct:
		num := v.NumField()
		for i := 0; i < num; i++ {
			name := v.Type().Field(i).Name
			v := v.Field(i).Interface()
			out[name] = v
		}
	default:
		err = errors.Errorf("not support type %#v", data)
	}
	return
}
