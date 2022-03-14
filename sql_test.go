// Copyright 2016 Davide Muzzarelli. All right reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gqt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

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
	testDir = filepath.Join(testDir, "../test")
	repo = NewRepository()
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
	sql, err := repo.GetSQL("parameter.getAllByAPIID", data)
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
	rpo := NewRepository()

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

func TestSQLIn(t *testing.T) {
	tpl := `
	{{define "testIn"}}
	 select * from aa where id {{in . .IDS}};
	 {{end}}
	`
	repo = NewRepository()
	err := repo.AddByNamespace("test", tpl, TemplatefuncMap)
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{
		"IDS": []int{1, 3, 4},
	}
	sqlrow, err := repo.GetSQL("test.testIn", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(sqlrow)
}
