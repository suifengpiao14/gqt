package tool

import (
	"testing"

	"github.com/suifengpiao14/gqt/v2"
)

func TestGenerateModel(t *testing.T) {
	repo := gqt.NewRepository()
	err := repo.AddByDir("../example", gqt.TemplatefuncMap)
	if err != nil {
		panic(err)
	}
	err = GenerateModel(repo)
	if err != nil {
		panic(err)
	}
}
