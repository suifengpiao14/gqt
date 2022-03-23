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
	"github.com/suifengpiao14/gqt/v2/gqttpl"
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
	gqttpl.DataVolumeMap
}

func (t *ListEntity) TplName() string {
	return "sql.list"
}
func (t *ListEntity) TplOutput() (string, error) {
	return "sql.list", nil
}

func TestSubDefineWhere(t *testing.T) {
	err := repo.AddByDir(testDir, TemplatefuncMap)
	if err != nil {
		panic(err)
	}

	entity := &ListEntity{
		IDS:           []int{1, 2, 3},
		DataVolumeMap: gqttpl.DataVolumeMap{},
	}
	sql, err := repo.GetSQLByTplEntity(entity)
	if err != nil {
		panic(err)
	}
	fmt.Println(sql)
}

func TestStruct(t *testing.T) {

	for _, dir := range []string{"pkg1", "pkg2"} {
		err := repo.AddByDir(filepath.Join(testDir, dir), TemplatefuncMap)
		if err != nil {
			panic(err)
		}
	}
	type structData struct {
		APIID int
		Ids   string
		*gqttpl.DataVolumeMap
	}

	data := &structData{
		APIID: 1,
		Ids:   "1,2,3,4,5,6",
	}
	sql, err := repo.GetSQL("parameter.getAllByAPIID", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(sql)

}

func TestMap(t *testing.T) {
	for _, dir := range []string{"pkg1", "pkg2"} {
		err := repo.AddByDir(filepath.Join(testDir, dir), TemplatefuncMap)
		if err != nil {
			panic(err)
		}
	}

	data := make(map[string]interface{})
	data["APIID"] = 1
	data["Ids"] = "1,2,4"
	dataVolume := gqttpl.DataVolumeMap(data)
	sql, err := repo.GetSQL("parameter.getAllByAPIID", &dataVolume)
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

type DataVolumeInt struct {
	dataVolume []int
}

// 测试指针类型转换
func TestPtrConvert(t *testing.T) {
	mapPtr := &map[string]interface{}{
		"IDS": []int{1, 3, 4},
	}
	volumePtr := &DataVolumeInt{
		dataVolume: []int{5, 6, 7},
	}
	//interfac := interface{}(mapPtr)

	a := (*DataVolumeInt)(unsafe.Pointer(mapPtr))
	*a = *volumePtr
	fmt.Println(a)
	fmt.Println(volumePtr)
	//fmt.Println(mapPtr)

	fmt.Printf("a:%d--v:%d----m:%d", unsafe.Pointer(a), unsafe.Pointer(volumePtr), unsafe.Pointer(mapPtr))
	fmt.Print("\n")
	fmt.Printf("a:%#v--v:%#v----m:", *a, *volumePtr)
}

func TestSQLInMap(t *testing.T) {
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
	data := map[string]interface{}{
		"IDS": []int{1, 3, 4},
	}
	dataVolume := gqttpl.DataVolumeMap(data)
	sqlrow, err := repo.GetSQL("test.testIn", &dataVolume)
	if err != nil {
		panic(err)
	}
	fmt.Println(sqlrow)
}

type DataVolumeTest struct {
	Hello string
	IDS   []int
	gqttpl.DataVolumeMap
}

func TestSQLInDataVolume(t *testing.T) {
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
	data := &DataVolumeTest{
		Hello: "hell",
		IDS:   []int{1, 3, 4},
		//DataVolumeMap: make(DataVolumeMap),
	}
	sqlrow, err := repo.GetSQL("test.testIn", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(sqlrow)
}

//go:embed  example
var RepositoryFS embed.FS

func TestAddByFS(t *testing.T) {
	rpo := NewRepositorySQL()
	err := rpo.AddByFS(RepositoryFS, ".", TemplatefuncMap)
	if err != nil {
		panic(err)
	}
}
