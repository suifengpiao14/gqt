package gqttpl

import (
	"embed"
	"fmt"
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
