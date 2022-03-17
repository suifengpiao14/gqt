package gqt

import (
	"io/fs"
	"reflect"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/suifengpiao14/gqt/v2/gqttpl"

	"github.com/pkg/errors"
)

// RepositorySQL stores SQL templates.
type RepositorySQL struct {
	templates map[string]*template.Template // namespace: template
}

// NewRepositorySQL create a new Repository.
func NewRepositorySQL() *RepositorySQL {
	return &RepositorySQL{
		templates: make(map[string]*template.Template),
	}
}

var LeftDelim = "{{"
var RightDelim = "}}"

type SQLRow struct {
	Name      string
	Namespace string
	SQL       string
	Statment  string
	Arguments []interface{}
	Result    interface{}
}

type TplEntity interface {
	TplName() string
}

func (r *RepositorySQL) AddByDir(root string, funcMap template.FuncMap) (err error) {
	r.templates, err = gqttpl.AddTemplateByDir(root, gqttpl.SQLNamespaceSuffix, funcMap, LeftDelim, RightDelim)
	if err != nil {
		return
	}
	ddlTemplates, err := gqttpl.AddTemplateByDir(root, gqttpl.DDLNamespaceSuffix, funcMap, LeftDelim, RightDelim)
	if err != nil {
		return
	}
	for fullname, tpl := range ddlTemplates {
		r.templates[fullname] = tpl
	}
	return
}

func (r *RepositorySQL) AddByFS(fsys fs.FS, root string, funcMap template.FuncMap) (err error) {
	r.templates, err = gqttpl.AddTemplateByFS(fsys, root, gqttpl.SQLNamespaceSuffix, funcMap, LeftDelim, RightDelim)
	if err != nil {
		return
	}
	ddlTemplates, err := gqttpl.AddTemplateByFS(fsys, root, gqttpl.DDLNamespaceSuffix, funcMap, LeftDelim, RightDelim)
	if err != nil {
		return
	}
	for fullname, tpl := range ddlTemplates {
		r.templates[fullname] = tpl
	}
	return
}

func (r *RepositorySQL) AddByNamespace(namespace string, content string, funcMap template.FuncMap) (err error) {
	t, err := gqttpl.AddTemplateByStr(namespace, content, funcMap, LeftDelim, RightDelim)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	r.templates[namespace] = t
	return
}

func (r *RepositorySQL) DefineResult2SQLRow(defineResult gqttpl.TPLDefine) (sqlRow *SQLRow, err error) {
	sqlRow = &SQLRow{
		Name:      defineResult.Name,
		Namespace: defineResult.Namespace,
	}

	sqlNamed := gqttpl.StandardizeSpaces(defineResult.Output)
	if sqlNamed == "" {
		return
	}
	data, err := getNamedData(defineResult.Input)
	if err != nil {
		return
	}
	sqlRow.Statment, sqlRow.Arguments, err = sqlx.Named(sqlNamed, data)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	sqlRow.SQL = gqttpl.Statement2SQL(sqlRow.Statment, sqlRow.Arguments)

	return
}

// GetByNamespace get all template under namespace
func (r *RepositorySQL) GetByNamespace(namespace string, data interface{}) (sqlRowList []*SQLRow, err error) {
	defineResultList, err := gqttpl.ExecuteNamespaceTemplate(r.templates, namespace, data)
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

func (r *RepositorySQL) GetDDLNamespace() (ddlNamespace string, err error) {
	for namespace := range r.templates {
		if strings.HasSuffix(namespace, gqttpl.DDLNamespaceSuffix) {
			ddlNamespace = namespace
			break
		}
	}
	if ddlNamespace == "" {
		err = errors.Errorf("not find ddl namespace with sufix %s,you can set gqt.DDLNamespaceSuffix to change sufix", gqttpl.DDLNamespaceSuffix)
		return
	}
	return
}

func (r *RepositorySQL) GetDDLSQL() (ddlSQLRowList []*SQLRow, err error) {
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
		sqlRow.SQL = gqttpl.StandardizeSpaces(sqlRow.SQL)
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

// 将模板名称，模板中的变量，封装到结构体中，使用结构体访问，避免拼写错误以及分散的硬编码，可以配合 gqttool 自动生成响应的结构体
func (r *RepositorySQL) GetSQLByTplEntity(t TplEntity) (sqlRow *SQLRow, err error) {
	return r.GetSQL(t.TplName(), t)
}

// GetSQLByTplEntityRef 支持只返回error 函数签名
func (r *RepositorySQL) GetSQLRawByTplEntityRef(t TplEntity, sqlStr *string) (err error) {
	sqlRow, err := r.GetSQL(t.TplName(), t)
	if err != nil {
		return err
	}
	*sqlStr = sqlRow.SQL
	return
}

//无sql注入的安全方式
func (r *RepositorySQL) GetSQL(fullname string, data interface{}) (sqlRow *SQLRow, err error) {
	defineResult, err := gqttpl.ExecuteTemplate(r.templates, fullname, data)
	if err != nil {
		return nil, err
	}
	sqlRow, err = r.DefineResult2SQLRow(*defineResult)
	return
}

func (r *RepositorySQL) NewSQLChain() *SQLChain {
	return &SQLChain{
		sqlRows:       make([]*SQLRow, 0),
		sqlRepository: func() *RepositorySQL { return r },
	}
}

func getNamedData(data interface{}) (out map[string]interface{}, err error) {
	out = make(map[string]interface{})
	if data == nil {
		return
	}
	for {
		dataI, ok := data.(*interface{})
		if ok {
			data = *dataI
		} else {
			break
		}
	}
	mapOut, ok := data.(map[string]interface{})
	if ok {
		out = mapOut
		return
	}
	if mapOut, ok := data.(DataVolumeMap); ok {
		out = mapOut
		return
	}

	// dataVolume, err := Convert2DataVolume(data)
	// if err != nil {
	// 	return
	// }

	v := reflect.Indirect(reflect.ValueOf(data))

	if v.Kind() != reflect.Struct {
		return
	}
	vt := v.Type()
	// 提取结构体field字段
	fieldNum := v.NumField()
	for i := 0; i < fieldNum; i++ {
		fv := v.Field(i)
		ft := fv.Type()
		fname := vt.Field(i).Name
		if fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
			ft = fv.Type()
		}
		ftk := ft.Kind()
		switch ftk {
		case reflect.Int:
			out[fname] = fv.Int()
		case reflect.Int64:
			out[fname] = int64(fv.Int())
		case reflect.Float64:
			out[fname] = fv.Float()
		case reflect.String:
			out[fname] = fv.String()
		case reflect.Struct, reflect.Map:
			subOut, err := getNamedData(fv.Interface())
			if err != nil {
				return out, err
			}
			for k, v := range subOut {
				if _, ok := out[k]; !ok {
					out[k] = v
				}
			}

		default:
			out[fname] = fv.Interface()
		}
	}
	// // 模板中动态数据增加/覆盖
	// dynamicData := dataVolume.GetDynamicValus()
	// for key, val := range dynamicData {
	// 	out[key] = val
	// }
	return
}
