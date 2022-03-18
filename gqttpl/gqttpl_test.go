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
		ID:            "data32",
		DataVolumeMap: &DataVolumeMap{},
	}
	data33 := &DataVolumeMapStruct{
		ID: "data33",
	}
	data34 := &DataVolumeMapStructRef{
		ID:            "data34",
		DataVolumeMap: &DataVolumeMap{},
	}
	data := []interface{}{data11, data12, data21, data22, data31, data32, data33, data34}
	for _, d := range data {
		dataVolume, ok := Interface2DataVolume(d)
		if ok {
			fmt.Printf("ok-%#v---", d)
			dataVolume.SetValue("hello", "world")
			dv, _ := d.(DataVolumeInterface)
			if dv != nil {
				hello, _ := dv.GetValue("hello")
				fmt.Println(hello)
			} else {
				fmt.Printf("is nil--")
			}

		} else {
			fmt.Printf("no-%#v-", d)
		}
		fmt.Printf("\n")
	}

}

func TestInterface2DataVolumeSpec(t *testing.T) {
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

func TestInterface2DataVolumeSpec32(t *testing.T) {
	d := &DataVolumeMapStructRef{
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
