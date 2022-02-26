package tool

import (
	"fmt"
	"strings"
)

func Backquote(s string) (out string) {
	out = fmt.Sprintf("`%s`", s)
	return
}
func Crud(table *Table) (tplMap map[string]string, err error) {
	tplMap = make(map[string]string)
	tplMap["TplGetByPrimaryKey"] = TplGetByPrimaryKey(table)
	tplMap["TplTotal"] = TplTotal(table)
	tplMap["TplPaginate"] = TplPaginate(table)
	tplMap["TplInsert"] = TplInsert(table)
	tplMap["TplUpdate"] = TplUpdate(table)
	tplMap["TplDel"] = TplDel(table)
	return
}

func isIgnoreColumn(columnName string, table *Table) (yes bool) {
	ignoreColumnMap := make(map[string]string)
	ignoreColumnMap[table.CreatedAtColumn] = table.CreatedAtColumn
	ignoreColumnMap[table.UpdateAtColumn] = table.UpdateAtColumn
	ignoreColumnMap[table.DeleteAtColumn] = table.DeleteAtColumn
	_, yes = ignoreColumnMap[columnName]
	return
}

func TplGetByPrimaryKey(table *Table) (tpl string) {
	tpl = fmt.Sprintf("select * from `%s`  where `%s`=:%s", table.TableName, table.PrimaryKey, table.PrimaryKeyCamel)
	if table.DeleteAtColumn != "" {
		tpl = fmt.Sprintf("%s  and `%s` is null", tpl, table.DeleteAtColumn)
	}
	tpl = tpl + ";\n"
	return
}

func TplTotal(table *Table) (tpl string) {
	tpl = fmt.Sprintf("select count(*) as `count` from `%s`  where 1=1 ", table.TableName)
	if table.DeleteAtColumn != "" {
		tpl = fmt.Sprintf("%s  and `%s` is null", tpl, table.DeleteAtColumn)
	}
	tpl = tpl + ";\n"
	return
}

func TplPaginate(table *Table) (tpl string) {
	tpl = fmt.Sprintf("select * from `%s`  where 1=1 ", table.TableName)
	if table.DeleteAtColumn != "" {
		tpl = fmt.Sprintf("%s  and `%s` is null ", tpl, table.DeleteAtColumn)
	}
	if table.UpdateAtColumn != "" {
		tpl = fmt.Sprintf(" %s order by `%s` desc  limit :Offset,:Limit ", tpl, table.UpdateAtColumn)
	}
	tpl = tpl + ";\n"
	return
}

func TplInsert(table *Table) (tpl string) {
	columns := make([]string, 0)
	values := make([]string, 0)
	for _, column := range table.Columns {
		if isIgnoreColumn(column.Name, table) {
			continue
		}
		columns = append(columns, Backquote(column.Name))
		values = append(values, ":"+column.CamelName)
	}

	columnStr := strings.Join(columns, ",")
	valueStr := strings.Join(values, ",")
	tpl = fmt.Sprintf("insert into `%s` (%s)values\n (%s);\n", table.TableName, columnStr, valueStr)
	return
}

func TplUpdate(table *Table) (tpl string) {
	updataList := make([]string, 0)
	for _, column := range table.Columns {
		if isIgnoreColumn(column.Name, table) {
			continue
		}
		str := fmt.Sprintf("{{if .%s}} {{$preComma.PreComma}} `%s`=:%s {{end}} ", column.CamelName, column.Name, column.CamelName)
		updataList = append(updataList, str)
	}
	updateTpl := strings.Join(updataList, "\n")
	tpl = fmt.Sprintf("{{$preComma:=newPreComma}}\n update `%s` set %s where `%s`=:%s;\n", table.TableName, updateTpl, table.PrimaryKey, table.PrimaryKeyCamel)
	return
}

func TplDel(table *Table) (tpl string) {
	tpl = fmt.Sprintf("update `%s` set `%s`={{currentTime}} where `%s`=:%s;\n", table.TableName, table.DeleteAtColumn, table.PrimaryKey, table.PrimaryKeyCamel)
	return
}
