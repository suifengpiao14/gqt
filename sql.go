package gqt

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm/logger"

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

var Suffix = ".sql.tpl"

func (r *Repository) AddByDir(root string, funcMap template.FuncMap) (err error) {
	// List the directories
	allFileList, err := GetTplFilesByDir(root)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	for _, filename := range allFileList {

		namespace := FileName2Namespace(filename, root, Suffix)
		t, err := template.New(namespace).Funcs(funcMap).ParseFiles(filename)
		if err != nil {
			return errors.WithStack(err)
		}
		r.templates[namespace] = t
	}
	return
}

func (r *Repository) AddByNamespace(namespace string, content string, funcMap template.FuncMap) (err error) {
	t, err := template.New(namespace).Funcs(funcMap).Parse(content)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	r.templates[namespace] = t
	return
}

// GetByNamespace get all template under namespace
func (r *Repository) GetByNamespace(namespace string, data interface{}) (sqlMap map[string]string, err error) {
	t, ok := r.templates[namespace]
	if !ok {
		err = errors.Errorf("not found namespace:%s", namespace)
		return
	}
	mapData, err := interface2map(data)
	if err != nil {
		return nil, err
	}
	sqlMap = make(map[string]string, 0)
	templates := t.Templates()
	for _, tpl := range templates {
		name := tpl.Name()
		var b bytes.Buffer
		err = tpl.Execute(&b, &mapData)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		fullName := fmt.Sprintf("%s.%s", namespace, name)
		content := strings.Trim(b.String(), "\r\n")
		if len(content) == 0 {
			continue
		}
		sqlNamed := b.String()
		sqlStatement, vars, err := sqlx.Named(sqlNamed, mapData)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		sqlStr := r.Statement2SQL(sqlStatement, vars)
		sqlMap[fullName] = sqlStr
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
func (r *Repository) GetSQLByTplEntity(t TplEntity) (sqlStr string, err error) {
	return r.GetSQL(t.TplName(), t)
}

// GetSQLByTplEntityRef 支持只返回error 函数签名
func (r *Repository) GetSQLByTplEntityRef(t TplEntity, sqlStr *string) (err error) {
	(*sqlStr), err = r.GetSQL(t.TplName(), t)
	return
}

//无sql注入的安全方式
func (r *Repository) GetSQL(name string, data interface{}) (sqlStr string, err error) {
	sqlStatement, vars, err := r.GetStatement(name, data)
	if err != nil {
		return
	}
	sqlStr = r.Statement2SQL(sqlStatement, vars)
	return
}

func (r *Repository) Statement2SQL(sqlStatement string, vars []interface{}) (sqlStr string) {
	sqlStr = logger.ExplainSQL(sqlStatement, nil, `'`, vars...)
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

type SQLRow struct {
	Tag       string
	SQL       string
	Statment  string
	Arguments []interface{}
	Result    interface{}
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
	sql, err := s.sqlRepository().GetSQL(tplName, args)
	if err != nil {
		s.err = err
		return s
	}
	sqlRow := &SQLRow{
		Tag:    tplName,
		SQL:    sql,
		Result: result,
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
func (s *SQLChain) AddSQL(tag string, sql string, result interface{}) {
	sqlRow := &SQLRow{
		Tag:    tag,
		SQL:    sql,
		Result: result,
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

var defaultRepository = NewRepository()

// AddByDir method for the default repository.
func AddByDir(dir string, funcMap template.FuncMap) error {
	return defaultRepository.AddByDir(dir, funcMap)
}

// AddByNamespace method for the default repository.
func AddByNamespace(filename string, content string, funcMap template.FuncMap) error {
	return defaultRepository.AddByNamespace(filename, content, funcMap)
}

func GetByNamespace(namespace string, data interface{}) (sqlMap map[string]string, err error) {
	return defaultRepository.GetByNamespace(namespace, data)
}

// Get method for the default repository.
func GetStatement(name string, data interface{}) (sql string, vars interface{}, err error) {
	return defaultRepository.GetStatement(name, data)
}

// Exec method for the default repository.
func GetSQL(name string, data interface{}) (sql string, e error) {
	return defaultRepository.GetSQL(name, data)
}

// Parse method for the default repository.
func Parse(name string, data interface{}) (string, error) {
	return defaultRepository.Parse(name, data)
}
