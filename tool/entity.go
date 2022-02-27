package tool

import (
	"bytes"
	"sort"
	"strings"
	"text/template"
)

//RepositoryEntity 根据数据表ddl和sql tpl 生成 sql tpl 调用的输入、输出实体
func RepositoryEntity(table *Table, sqlTpl string) (entityStruct string, err error) {
	variableMap := ParsSqlTplVariable(sqlTpl)
	entityTplData, err := GetEntityData(table, variableMap)
	if err != nil {
		return
	}

	tp, err := template.New("").Parse(EntityTpl())
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	err = tp.Execute(buf, entityTplData)
	if err != nil {
		return
	}
	entityStruct = buf.String()
	return
}

type EntityTplData struct {
	StructName string
	Attributes Variables
}

func GetEntityData(table *Table, variableMap map[string]*Variable) (entityTplData *EntityTplData, err error) {
	entityTplData = &EntityTplData{
		StructName: table.TableNameCamel,
		Attributes: make(Variables, 0),
	}
	tableColumnMap := make(map[string]*Column)
	for _, column := range table.Columns {
		tableColumnMap[column.CamelName] = column
	}

	for name, variable := range variableMap { // 使用数据库字段定义校正变量类型
		tableColumn, ok := tableColumnMap[name]
		if ok {
			if err != nil {
				return nil, err
			}
			variable.Type = tableColumn.Type
			entityTplData.Attributes = append(entityTplData.Attributes, variable)
			continue
		}
		// 不属于数据表字段变量，直接添加
		entityTplData.Attributes = append(entityTplData.Attributes, variable)
	}
	sort.Sort(entityTplData.Attributes)
	return
}
func EntityTpl() (tpl string) {
	tpl = `
		type {{.StructName}} struct{
			{{range .Attributes }}
				{{.Name}} {{.Type}} 
			{{end}}
		}
	`
	return
}

type Variable struct {
	Name       string
	Type       string
	AllowEmpty bool
}

// 按照 Person.Age 从大到小排序
type Variables []*Variable

func (v Variables) Len() int { // 重写 Len() 方法
	return len(v)
}
func (v Variables) Swap(i, j int) { // 重写 Swap() 方法
	v[i], v[j] = v[j], v[i]
}
func (v Variables) Less(i, j int) bool { // 重写 Less() 方法， 从小到大排序
	return v[i].Name < v[j].Name
}

func ParsSqlTplVariable(sqlTpl string) (variableMap map[string]*Variable) {
	variableMap = make(map[string]*Variable)
	byteArr := []byte(sqlTpl)
	leftDelim := byte('{')
	rightDelim := byte('}')
	itemBegin := false
	itemArr := make([][]byte, 0)
	item := make([]byte, 0)
	byteLen := len(byteArr)
	for i := 0; i < byteLen; i++ {
		c := byteArr[i]
		if c == leftDelim && i+1 < byteLen && byteArr[i+1] == leftDelim && !itemBegin {
			itemBegin = true
			item = make([]byte, 0)
			i++
			continue
		}
		if c == rightDelim && i+1 < byteLen && byteArr[i+1] == rightDelim && itemBegin {
			itemBegin = false
			itemArr = append(itemArr, item)
			i++
			continue
		}
		if itemBegin {
			item = append(item, c)
		}
	}
	for _, item := range itemArr {
		variable, _ := parsePrefixVariable(item, byte('.'))
		if variable.Name != "" {
			variableMap[variable.Name] = &variable
		}

	}

	// parse sql variable
	sqlVariableDelim := byte(':')

	for {
		variable, pos := parsePrefixVariable(byteArr, sqlVariableDelim)
		if variable.Name == "" {
			break
		}
		variableMap[variable.Name] = &variable
		pos += len(variable.Name)
		byteArr = byteArr[pos:]
	}
	limitVariabeMap := GetLimitVariable(sqlTpl)
	for name, variable := range limitVariabeMap {
		variableMap[name] = variable
	}
	return
}

func GetLimitVariable(sqlTpl string) (variableMap map[string]*Variable) {
	variableMap = make(map[string]*Variable)
	index := strings.Index(strings.ToLower(sqlTpl), "limit")
	if index < 0 {
		return
	}
	byteArr := []byte(sqlTpl[index:])
	// parse sql variable
	sqlVariableDelim := byte(':')

	for {
		variable, pos := parsePrefixVariable(byteArr, sqlVariableDelim)
		if variable.Name == "" {
			break
		}
		variable.Type = "int"
		variable.AllowEmpty = false
		variableMap[variable.Name] = &variable
		pos += len(variable.Name)
		byteArr = byteArr[pos:]
	}

	return
}

// 找到第一个变量
func parsePrefixVariable(item []byte, variableStart byte) (variable Variable, pos int) {
	variableBegin := false
	pos = 0
	variableNameByte := make([]byte, 0)
	itemLen := len(item)
	for j := 0; j < itemLen; j++ {
		c := item[j]
		if c == variableStart && j+1 < itemLen && IsNameChar(item[j+1]) { // c=变量标示字符，并且后面跟字符，是变量名的必要条件
			if j == 0 || !IsNameChar(item[j-1]) { // 变量符号开始或者前面不为字符，为变量名的充要条件
				variableBegin = true
				pos = j
				continue
			}
		}
		if variableBegin {
			if IsNameChar(c) {
				variableNameByte = append(variableNameByte, c)
			} else {
				break
			}
		}
	}
	variable = Variable{
		Name:       string(variableNameByte),
		Type:       "interface{}",
		AllowEmpty: true,
	}
	return
}

// 判断是否可以作为名称的字符
func IsNameChar(c byte) (yes bool) {
	yes = false
	a := byte('a')
	z := byte('z')
	A := byte('A')
	Z := byte('Z')
	underline := byte('_')
	if (a <= c && c <= z) || (A <= c && c <= Z) || c == underline {
		yes = true
	}
	return
}
