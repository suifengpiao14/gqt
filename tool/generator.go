package tool

import (
	"fmt"

	"github.com/suifengpiao14/gqt/v2"
)

func GenerateCrud() {

}
func GenerateModel(rep *gqt.Repository) (err error) {
	ddlList, err := getDDLFromRepository(rep)
	if err != nil {
		return
	}
	tableList, err := GenerateTable(ddlList)
	if err != nil {
		return
	}
	tableStructList, err := GenerateTableStruct(tableList)
	fmt.Println(tableStructList)
	return
}

func getDDLFromRepository(rep *gqt.Repository) (ddlList []string, err error) {
	ddlList = make([]string, 0)
	ddlMap, err := rep.GetByNamespace("ddl", nil)
	if err != nil {
		return
	}
	for _, ddl := range ddlMap {
		ddlList = append(ddlList, ddl)
	}
	return
}
