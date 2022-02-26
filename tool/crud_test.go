package tool

import (
	"fmt"
	"testing"
)

func TestCrud(t *testing.T) {
	ddlList := GetDDL()
	tables, err := GenerateTable(ddlList)
	if err != nil {
		panic(err)
	}
	for _, table := range tables {
		tplMap, err := Crud(table)
		if err != nil {
			panic(err)
		}
		fmt.Println(tplMap)
	}

}
