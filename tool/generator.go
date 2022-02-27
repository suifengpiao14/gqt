package tool

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/suifengpiao14/errorformatter"
	"github.com/suifengpiao14/gqt/v2"
)

func main() {
	var (
		ddlFile   string
		modelDir  string
		tplDir    string
		entiryDir string
		err       error
	)
	modelCmd := flag.NewFlagSet("model", flag.ExitOnError)
	crudCmd := flag.NewFlagSet("crud", flag.ExitOnError)
	entiryCmd := flag.NewFlagSet("entity", flag.ExitOnError)
	modelCmd.StringVar(&ddlFile, "ddl", "template/ddl.sql.tpl", "ddl template file name")
	modelCmd.StringVar(&modelDir, "dir", "model", "ddl template file name")

	crudCmd.StringVar(&ddlFile, "ddl", "template/ddl.sql.tpl", "ddl template file name")
	crudCmd.StringVar(&tplDir, "dir", "template", "sql template dir")

	entiryCmd.StringVar(&tplDir, "tplDir", "template", "sql template dir")
	entiryCmd.StringVar(&entiryDir, "entiryDir", "entity", "entity dir")
	testing.Init()

	if len(os.Args) < 3 {
		fmt.Println("expected 'model' or 'crud' or 'entiry' subcommands")
		os.Exit(1)
	}
	args := os.Args[2:]
	switch os.Args[1] {
	case "model":
		modelCmd.Parse(args)
		repo := gqt.NewRepository()
		errChain := errorformatter.NewErrorChain()
		var content string
		var modelMap map[string]string
		errChain.Run(func() (err error) {
			content, err = GetFileContent(ddlFile)
			return
		}).
			SetError(repo.AddByNamespace("ddl", content, gqt.TemplatefuncMap)).
			Run(func() (err error) {
				modelMap, err = GenerateTableModel(repo)
				return
			})
		err = errChain.Error()
		if err != nil {
			panic(err)
		}
		for name, model := range modelMap {
			snakeName := strcase.ToSnake(name)
			filename := fmt.Sprintf("%s/%s.model.go", modelDir, snakeName)
			err = saveFile(filename, model)
			if err != nil {
				panic(err)
			}
		}

	case "crud":
		crudCmd.Parse(args)
		repo := gqt.NewRepository()
		errChain := errorformatter.NewErrorChain()
		var content string
		var tplMap map[string]string
		errChain.Run(func() (err error) {
			content, err = GetFileContent(ddlFile)
			return
		}).
			SetError(repo.AddByNamespace("ddl", content, gqt.TemplatefuncMap)).
			Run(func() (err error) {
				tplMap, err = GenerateCrud(repo)
				return
			})

		err = errChain.Error()
		if err != nil {
			panic(err)
		}
		for name, tpl := range tplMap {
			snakeName := strcase.ToSnake(name)
			filename := fmt.Sprintf("%s/%s.sql.tpl", tplDir, snakeName)
			err = saveFile(filename, tpl)
			if err != nil {
				panic(err)
			}
		}
	case "entity":
		entiryCmd.Parse(args)
		repo := gqt.NewRepository()
		errChain := errorformatter.NewErrorChain()
		var entityMap map[string]string
		errChain.SetError(repo.AddByDir(entiryDir, gqt.TemplatefuncMap)).
			Run(func() (err error) {
				entityMap, err = GenerateEntity(repo)
				return
			})

		err = errChain.Error()
		if err != nil {
			panic(err)
		}
		for name, entity := range entityMap {
			snakeName := strcase.ToSnake(name)
			filename := fmt.Sprintf("%s/%s.entity.go", entiryDir, snakeName)
			err = saveFile(filename, entity)
			if err != nil {
				panic(err)
			}
		}
	default:
		fmt.Println("expected 'model' or 'crud' or 'entiry' subcommands")
		os.Exit(1)
	}
}

func GetFileContent(file string) (content string, err error) {
	f, err := os.OpenFile(file, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return
	}
	content = string(b)
	return

}

func saveFile(filename string, content string) (err error) {
	if IsExist(filename) {
		return
	}
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	_, err = f.WriteString(content)
	if err != nil {
		return
	}
	return

}

func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func GenerateEntity(rep *gqt.Repository) (entityMap map[string]string, err error) {
	ddlList, err := getDDLFromRepository(rep)
	if err != nil {
		return
	}
	tableList, err := GenerateTable(ddlList)
	if err != nil {
		return
	}
	for _, table := range tableList {
		//todo get sqlTpl
		sqlTpl := ""
		entityStruct, err := RepositoryEntity(table, sqlTpl)
		if err != nil {
			return nil, err
		}
		entityMap[table.TableNameCamel] = entityStruct
	}

	return
}

func GenerateCrud(rep *gqt.Repository) (tplMap map[string]string, err error) {
	ddlList, err := getDDLFromRepository(rep)
	if err != nil {
		return
	}
	tableList, err := GenerateTable(ddlList)
	if err != nil {
		return
	}
	for _, table := range tableList {
		tplMap, err = Crud(table)
		if err != nil {
			return
		}
	}

	return
}
func GenerateTableModel(rep *gqt.Repository) (modelMap map[string]string, err error) {
	ddlList, err := getDDLFromRepository(rep)
	if err != nil {
		return
	}
	tableList, err := GenerateTable(ddlList)
	if err != nil {
		return
	}
	modelMap, err = GenerateModel(tableList)
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
