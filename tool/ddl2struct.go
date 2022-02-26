package tool

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	executor "github.com/bytewatch/ddl-executor"
	"github.com/iancoleman/strcase"
)

//map for converting mysql type to golang types
var typeForMysqlToGo = map[string]string{
	"int":                "int64",
	"integer":            "int64",
	"tinyint":            "int64",
	"smallint":           "int64",
	"mediumint":          "int64",
	"bigint":             "int64",
	"int unsigned":       "int64",
	"integer unsigned":   "int64",
	"tinyint unsigned":   "int64",
	"smallint unsigned":  "int64",
	"mediumint unsigned": "int64",
	"bigint unsigned":    "int64",
	"bit":                "int64",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "time.Time", // time.Time or string
	"datetime":           "time.Time", // time.Time or string
	"timestamp":          "time.Time", // time.Time or string
	"time":               "time.Time", // time.Time or string
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

func isTimeMysqlType(mysqlType string) bool {
	timeMap := map[string]string{
		"date":      "date",
		"datetime":  "datetime",
		"timestamp": "timestamp",
		"time":      "time",
	}
	_, ok := timeMap[mysqlType]
	return ok
}
func mysql2GoType(mysqlType string, time2str bool) (goType string, err error) {
	if time2str {
		typeForMysqlToGo["date"] = "string"
		typeForMysqlToGo["datetime"] = "string"
		typeForMysqlToGo["timestamp"] = "string"
		typeForMysqlToGo["time"] = "string"
	}
	index := strings.Index(mysqlType, "(")
	if index > -1 {
		mysqlType = mysqlType[:index]
	}
	goType, ok := typeForMysqlToGo[mysqlType]
	if !ok {
		err = errors.Errorf("mysql2GoType: not found mysql type %s to go type", mysqlType)
	}
	return

}

type Column struct {
	CamelName string
	Name      string
	Type      string
	Comment   string
	Tag       string
	Nullable  bool
}

type Table struct {
	TableName       string
	TableNameCamel  string
	PrimaryKey      string
	PrimaryKeyCamel string
	DeleteAtColumn  string
	DeleteAtCamel   string
	CreatedAtColumn string
	CreatedAtCamel  string
	UpdateAtColumn  string
	UpdateAtCamel   string
	Columns         []*Column
}

func GenerateTableStruct(tableList []*Table) (tableStructList []string, err error) {
	tableStructList = make([]string, 0)
	tableTpl := structTpl()
	tl, err := template.New("").Parse(tableTpl)
	if err != nil {
		return
	}

	for i := 0; i < len(tableList); i++ {
		buf := new(bytes.Buffer)
		table := tableList[i]
		err = tl.Execute(buf, table)
		if err != nil {
			return
		}
		tableStructList = append(tableStructList, buf.String())
	}
	return
}

func GenerateTable(ddlList []string) (tables []*Table, err error) {
	tables = make([]*Table, 0)
	conf := executor.NewDefaultConfig()
	inst := executor.NewExecutor(conf)
	databaseName := "test"
	ddlDB := fmt.Sprintf("create database `%s`;use `%s`;", databaseName, databaseName)
	ddlAll := make([]string, 0)
	ddlAll = append(ddlAll, ddlDB)
	ddlAll = append(ddlAll, ddlList...)

	ddls := strings.Join(ddlAll, ";")

	err = inst.Exec(ddls)
	if err != nil {
		return nil, err
	}

	tableList, err := inst.GetTables(databaseName)
	if err != nil {
		return
	}

	for i := 0; i < len(tableList); i++ {
		tableName := tableList[i]
		tableDef, err := inst.GetTableDef(databaseName, tableName)
		if err != nil {
			return nil, err
		}

		table := &Table{
			TableName:      tableName,
			TableNameCamel: strcase.ToCamel(tableName),
			Columns:        make([]*Column, 0),
		}
		for _, indice := range tableDef.Indices {
			if indice.Name == "PRIMARY" {
				table.PrimaryKey = indice.Columns[0] // 暂时取第一个为主键，不支持多字段主键
				table.PrimaryKeyCamel = strcase.ToCamel(table.PrimaryKey)
			}
		}
		for _, columnDef := range tableDef.Columns {

			goType, err := mysql2GoType(columnDef.Type, true)
			if err != nil {
				return nil, err
			}
			if columnDef.OnUpdate {
				table.UpdateAtColumn = columnDef.Name
				table.UpdateAtCamel = strcase.ToCamel(columnDef.Name)
			}
			if strings.Contains(columnDef.DefaultValue, "current_timestamp") && !columnDef.OnUpdate {
				table.CreatedAtColumn = columnDef.Name
				table.CreatedAtCamel = strcase.ToCamel(columnDef.Name)
			}
			if isTimeMysqlType(columnDef.Type) && strings.Contains(columnDef.Name, "deleted_at") {
				table.DeleteAtColumn = columnDef.Name
				table.DeleteAtCamel = strcase.ToCamel(columnDef.Name)
			}
			columnPt := &Column{
				CamelName: strcase.ToCamel(columnDef.Name),
				Name:      columnDef.Name,
				Type:      goType,
				Comment:   columnDef.Comment,
				Nullable:  columnDef.Nullable,
				Tag:       fmt.Sprintf("`json:\"%s\"`", strcase.ToLowerCamel(columnDef.Name)),
			}
			table.Columns = append(table.Columns, columnPt)
		}
		tables = append(tables, table)
	}
	return
}

func structTpl() string {
	return `
		type {{.TableName}} struct{
			{{range .Columns }} 
			// {{.Comment}}
			{{.CamelName}} {{.Type}} {{if .Tag}} {{.Tag}} {{end}}
			{{end}}
		}
	`
}
