package tool

//RepositoryEntity 根据数据表ddl和sql tpl 生成 sql tpl 调用的输入、输出实体
func RepositoryEntity(table *Table, sqlTpl string) (entities map[string]string, err error) {

	return
}

func ParsSqlTplVariable(sqlTpl string) (variableList []string) {
	variableMap := make(map[string]string)
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
		if variable != "" {
			variableMap[variable] = variable
		}

	}

	// parse sql variable
	sqlVariableDelim := byte(':')

	for {
		variable, pos := parsePrefixVariable(byteArr, sqlVariableDelim)
		if variable == "" {
			break
		}
		variableMap[variable] = variable
		pos += len(variable)
		byteArr = byteArr[pos:]

	}
	for variable := range variableMap {
		variableList = append(variableList, variable)
	}

	return
}

// 找到第一个变量
func parsePrefixVariable(item []byte, variableStart byte) (variable string, pos int) {
	variableBegin := false
	pos = 0
	variableNameByte := make([]byte, 0)
	for j := 0; j < len(item); j++ {
		c := item[j]
		if c == variableStart {
			if j == 0 {
				variableBegin = true
				pos = j
				continue
			}
			if !IsNameChar(item[j-1]) {
				variableBegin = true
				pos = j
				continue
			}
		}
		if variableBegin {
			if IsNameChar(c) {
				variableNameByte = append(variableNameByte, c)
			} else if len(variableNameByte) > 0 {
				break
			} else {
				variableBegin = false
			}
		}
	}
	variable = string(variableNameByte)
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
