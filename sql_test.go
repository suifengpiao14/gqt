// Copyright 2016 Davide Muzzarelli. All right reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gqt

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/jmoiron/sqlx"
)

var testDir string
var repo *RepositorySQL

func init() {
	var err error
	testDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	testDir = filepath.Join(testDir, "example")
	repo = NewRepositorySQL()
}

type ListEntity struct {
	IDS []int
	TplEmptyEntity
}

func (t *ListEntity) TplName() string {
	return "sql.list"
}
func (t *ListEntity) TplOutput(tplEntity TplEntityInterface) (string, error) {
	return "sql.list", nil
}

func TestSubDefineWhere(t *testing.T) {
	err := repo.AddByDir(testDir, TemplatefuncMap)
	if err != nil {
		panic(err)
	}

	entity := &ListEntity{
		IDS:            []int{1, 2, 3},
		TplEmptyEntity: TplEmptyEntity{},
	}
	sql, err := repo.GetSQL(entity)
	if err != nil {
		panic(err)
	}
	fmt.Println(sql)
}

func TestSQLNamed(t *testing.T) {
	namedSql := "select * from `test_table` where `id`=:id and name=:name"
	data := map[string]interface{}{
		"id":   1,
		"name": "hahha",
		"more": "more",
	}
	sqlStatement, vars, err := sqlx.Named(namedSql, data)
	if err != nil {
		panic(err)
	}
	fmt.Println(vars)
	fmt.Println(sqlStatement)
}

func TestGetDDLSQL(t *testing.T) {
	rpo := NewRepositorySQL()

	err := rpo.AddByDir("example", TemplatefuncMap)
	if err != nil {
		panic(err)
	}
	ddlMap, err := rpo.GetDDLSQL()
	if err != nil {
		panic(err)
	}
	fmt.Println(ddlMap)
}

type ModelStruct struct {
	TableName  string
	PrimaryKey string
}

func (s *ModelStruct) PrimaryKeyCamel() string {
	return "ID"
}

type tplEntityInt struct {
	tplEntity []int
}

// 测试指针类型转换
func TestPtrConvert(t *testing.T) {
	mapPtr := &map[string]interface{}{
		"IDS": []int{1, 3, 4},
	}
	volumePtr := &tplEntityInt{
		tplEntity: []int{5, 6, 7},
	}
	//interfac := interface{}(mapPtr)

	a := (*tplEntityInt)(unsafe.Pointer(mapPtr))
	*a = *volumePtr
	fmt.Println(a)
	fmt.Println(volumePtr)
	//fmt.Println(mapPtr)

	fmt.Printf("a:%d--v:%d----m:%d", unsafe.Pointer(a), unsafe.Pointer(volumePtr), unsafe.Pointer(mapPtr))
	fmt.Print("\n")
	fmt.Printf("a:%#v--v:%#v----m:", *a, *volumePtr)
}

type tplEntityTest struct {
	Hello string
	IDS   []int
	TplEmptyEntity
}

func TestSQLIntplEntity(t *testing.T) {
	tpl := `
	{{define "testIn"}}
	 select * from aa where id {{in . .IDS}};
	 {{end}}
	`
	repo = NewRepositorySQL()
	err := repo.AddByNamespace("test", tpl, TemplatefuncMap)
	if err != nil {
		panic(err)
	}
	data := &tplEntityTest{
		Hello: "hell",
		IDS:   []int{1, 3, 4},
		//tplEntityMap: make(tplEntityMap),
	}
	sqlrow, err := repo.GetSQL(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(sqlrow)
}

//go:embed  example
var repositoryFS embed.FS

func TestAddByFS(t *testing.T) {
	rpo := NewRepositorySQL()
	err := rpo.AddByFS(repositoryFS, ".", TemplatefuncMap)
	if err != nil {
		panic(err)
	}
}
