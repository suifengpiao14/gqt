package gqttpl

import (
	"embed"
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

type DataVolumeMapStruct struct {
	ID string
	DataVolumeMap
}

type DataVolumeMapStructRef struct {
	ID string
	*DataVolumeMap
}

func TestInterface2DataVolume(t *testing.T) {
	data11 := &DataVolumeMap{
		"ID": "data11",
	}
	data12 := DataVolumeMap{
		"ID": "data12",
	}
	data21 := &map[string]interface{}{
		"ID": "data21",
	}
	data22 := map[string]interface{}{
		"ID": "data22",
	}

	data31 := DataVolumeMapStruct{
		ID: "data31",
	}
	data32 := DataVolumeMapStructRef{
		ID: "data32",
	}
	data33 := &DataVolumeMapStruct{
		ID: "data33",
	}
	data34 := &DataVolumeMapStructRef{
		ID: "data34",
	}
	data := []interface{}{data11, data12, data21, data22, data31, data32, data33, data34}
	for _, d := range data {
		di := interface{}(d)
		dataVolume, ok := Interface2DataVolume(&di)
		if ok {
			dataVolume.SetValue("hello", "world")
			dv, _ := d.(DataVolumeInterface)
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

//TestInterface2DataVolumeDataVolumeMap 测试 传入 VolumeDataVolumeMap 类型
func TestInterface2DataVolumeDataVolumeMap(t *testing.T) {
	d := DataVolumeMap{
		"ID": "data12",
	}
	dataVolume, ok := Interface2DataVolume(d)
	if !ok {
		fmt.Printf("no-%#v\n", d)
	} else {
		fmt.Printf("ok-%#v\n", d)
		dataVolume.SetValue("hello", "world")
		hello, _ := d.GetValue("hello")
		fmt.Println(hello)
	}
}

// TestInterface2DataVolumeMap 测试传入 map[string]interface{}类型
func TestInterface2DataVolumeMap(t *testing.T) {
	d := map[string]interface{}{
		"ID": "data12",
	}

	dataVolume, ok := Interface2DataVolume(d)
	if !ok {
		fmt.Printf("no-%#v\n", d)
	} else {
		dataVolume.SetValue("hello", "world")
		fmt.Printf("ok-%#v\n", d)
	}
}

//TestInterface2DataVolumeStruct 测试传入 struct 类型
func TestInterface2DataVolumeStruct(t *testing.T) {
	d := &DataVolumeMapStruct{
		ID: "data31",
	}
	dataVolume, ok := Interface2DataVolume(d)
	if !ok {
		fmt.Printf("no-%#v\n", d)
	} else {

		dataVolume.SetValue("hello", "world")
		fmt.Printf("ok-%#v\n", d)
		hello, _ := d.GetValue("hello")
		fmt.Println(hello)
	}
	d.SetValue("a", "b")
	fmt.Printf("ok-%#v\n", d)
}

func TestInterface2DataVolumeStructRef(t *testing.T) {
	d := DataVolumeMapStructRef{
		ID: "data32",
	}
	dataVolume, ok := Interface2DataVolume(d)
	if !ok {
		fmt.Printf("no-%#v\n", d)
	} else {
		fmt.Printf("ok-%#v\n", d)
		dataVolume.SetValue("hello", "world")
		hello, _ := d.GetValue("hello")
		fmt.Println(hello)
	}
}

/**
&gqttpl.DataVolumeMapStructRef{ID:"data32", DataVolumeMap:(*gqttpl.DataVolumeMap)(nil)}ok-&gqttpl.DataVolumeMapStructRef{ID:"data32", DataVolumeMap:(*gqttpl.DataVolumeMap)(0xc00020e1f0)}
**/
