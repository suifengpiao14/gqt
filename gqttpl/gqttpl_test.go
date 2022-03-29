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

type tplEntityMapStruct struct {
	ID string
	TplEmptyEntity
}

type tplEntityMapStructRef struct {
	ID string
	*TplEmptyEntity
}

func TestInterface2tplEntity(t *testing.T) {
	data11 := &TplEmptyEntity{
		"ID": "data11",
	}
	data12 := TplEmptyEntity{
		"ID": "data12",
	}
	data21 := &map[string]interface{}{
		"ID": "data21",
	}
	data22 := map[string]interface{}{
		"ID": "data22",
	}

	data31 := tplEntityMapStruct{
		ID: "data31",
	}
	data32 := tplEntityMapStructRef{
		ID: "data32",
	}
	data33 := &tplEntityMapStruct{
		ID: "data33",
	}
	data34 := &tplEntityMapStructRef{
		ID: "data34",
	}
	data := []interface{}{data11, data12, data21, data22, data31, data32, data33, data34}
	for _, d := range data {
		di := interface{}(d)
		tplEntity, ok := Interface2tplEntity(&di)
		if ok {
			tplEntity.SetValue("hello", "world")
			dv, _ := d.(TplEntityInterface)
			if dv != nil {
				hello, _ := dv.GetValue("hello")
				fmt.Print(hello)
				fmt.Printf("--dv--%#v", d)
			} else {
				fmt.Printf("%#v", d)
			}

		} else {
			fmt.Printf("no-%#v-", d)
		}
		fmt.Printf("\n")
	}

}

//TestInterface2tplEntitytplEntityMap 测试 传入 VolumetplEntityMap 类型
func TestInterface2tplEntitytplEntityMap(t *testing.T) {
	d := TplEmptyEntity{
		"ID": "data12",
	}
	tplEntity, ok := Interface2tplEntity(d)
	if !ok {
		fmt.Printf("no-%#v\n", d)
	} else {
		fmt.Printf("ok-%#v\n", d)
		tplEntity.SetValue("hello", "world")
		hello, _ := d.GetValue("hello")
		fmt.Println(hello)
	}
}

// TestInterface2tplEntityMap 测试传入 map[string]interface{}类型
func TestInterface2tplEntityMap(t *testing.T) {
	d := map[string]interface{}{
		"ID": "data12",
	}

	tplEntity, ok := Interface2tplEntity(d)
	if !ok {
		fmt.Printf("no-%#v\n", d)
	} else {
		tplEntity.SetValue("hello", "world")
		fmt.Printf("ok-%#v\n", d)
	}
}

//TestInterface2tplEntityStruct 测试传入 struct 类型
func TestInterface2tplEntityStruct(t *testing.T) {
	d := &tplEntityMapStruct{
		ID: "data31",
	}
	tplEntity, ok := Interface2tplEntity(d)
	if !ok {
		fmt.Printf("no-%#v\n", d)
	} else {

		tplEntity.SetValue("hello", "world")
		fmt.Printf("ok-%#v\n", d)
		hello, _ := d.GetValue("hello")
		fmt.Println(hello)
	}
	d.SetValue("a", "b")
	fmt.Printf("ok-%#v\n", d)
}

func TestInterface2tplEntityStructRef(t *testing.T) {
	d := tplEntityMapStructRef{
		ID: "data32",
	}
	tplEntity, ok := Interface2tplEntity(d)
	if !ok {
		fmt.Printf("no-%#v\n", d)
	} else {
		fmt.Printf("ok-%#v\n", d)
		tplEntity.SetValue("hello", "world")
		hello, _ := d.GetValue("hello")
		fmt.Println(hello)
	}
}

/**
&gqttpl.tplEntityMapStructRef{ID:"data32", tplEntityMap:(*gqttpl.tplEntityMap)(nil)}ok-&gqttpl.tplEntityMapStructRef{ID:"data32", tplEntityMap:(*gqttpl.tplEntityMap)(0xc00020e1f0)}
**/
