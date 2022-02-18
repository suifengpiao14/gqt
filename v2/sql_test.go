// Copyright 2016 Davide Muzzarelli. All right reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"path/filepath"
	"testing"
)

var testDir string

func init() {
	var err error
	testDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	testDir = filepath.Join(testDir, "../test")
}

func TestStruct(t *testing.T) {
	for _, dir := range []string{"pkg1", "pkg2"} {
		err := AddByDir(filepath.Join(testDir, dir), TemplatefuncMap)
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
	sql, err := GetSQL("parameter.getAllByAPIID", data)
	if err != nil {
		panic(err)
	}
	fmt.Printf(sql)

}

func TestMap(t *testing.T) {
	for _, dir := range []string{"pkg1", "pkg2"} {
		err := AddByDir(filepath.Join(testDir, dir), TemplatefuncMap)
		if err != nil {
			panic(err)
		}
	}

	data := make(map[string]interface{})
	data["APIID"] = 1
	data["Ids"] = "1,2,4"
	sql, err := GetSQL("parameter.getAllByAPIID", data)
	if err != nil {
		panic(err)
	}
	fmt.Printf(sql)

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
