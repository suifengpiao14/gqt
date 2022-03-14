package gqttpl

import (
	"fmt"
	"testing"
)

func TestGetTplFilesByDir(t *testing.T) {
	dir := "."
	suffix := ".sql.tpl"
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
