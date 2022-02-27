package tool

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/gqt/v2"
)

func TestGenerateModel(t *testing.T) {
	repo := gqt.NewRepository()
	err := repo.AddByDir("../example", gqt.TemplatefuncMap)
	if err != nil {
		panic(err)
	}
	tableModelList, err := GenerateTableModel(repo)
	if err != nil {
		panic(err)
	}
	fmt.Println(tableModelList)
}
