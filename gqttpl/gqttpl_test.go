package gqttpl

import (
	"embed"
	"fmt"
	"reflect"
	"testing"
)

func TestGetTplFilesByDir(t *testing.T) {
	dir := "."
	suffix := "sql"
	fileList, err := GetTplFilesByDir(dir, suffix)
	if err != nil {
		panic(err)
	}
	fmt.Println(fileList)
}
func TestStandardizeSpaces(t *testing.T) {
	s := `
	a     b 
	c	d
	`
	ns := StandardizeSpaces(s)
	fmt.Println(ns)
}

//go:embed  test
var RepositoryFS embed.FS

func TestAddByFS(t *testing.T) {
	pattern := "test/data/**/*.tpl"
	files, err := Glob(RepositoryFS, pattern)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)
}

func TestGlobDirectory(t *testing.T) {
	dir := "../"
	pattern := "**.go"
	files, err := GlobDirectory(dir, pattern)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)
}

type GenExampleSQLInsertEntity struct {
	APIID            string
	Auth             string
	Body             string
	ContentType      string
	ExampleID        string
	Headers          string
	Method           string
	Parameters       string
	PreRequestScript string
	Response         string
	ServiceID        string
	Summary          string
	Tag              string
	TestScript       string
	Title            string
	URL              string
	TplEmptyEntity
}

func (t *GenExampleSQLInsertEntity) TplName() string {
	return "gen.example.sql.Insert"
}
func (t *GenExampleSQLInsertEntity) TplType() string {
	return "sql_insert"
}

type BatchArgs struct {
	InsertEntity *GenExampleSQLInsertEntity
}

func TestNilTplEntity(t *testing.T) {
	args := &BatchArgs{}
	SetTplEntity(args.InsertEntity)

}

func SetTplEntity(t TplEntityInterface) {
	if t == nil {
		fmt.Println(111)
		return
	}
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Ptr {
		fmt.Println("ptr")
	}
	fmt.Println(nil)
	fmt.Println(t)
	// rv := reflect.ValueOf(t)
	// println(rv.IsNil())
	// err := fmt.Sprintf("%#v must be not nil", t)
	// fmt.Println(err)
	//t.SetValue("a", "ok")
}
