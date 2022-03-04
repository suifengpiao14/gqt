package gqt

import (
	"fmt"
	"testing"
)

func TestGetTplFilesByDir(t *testing.T) {
	dir := "."
	fileList, err := GetTplFilesByDir(dir)
	if err != nil {
		panic(err)
	}
	fmt.Println(fileList)
}
