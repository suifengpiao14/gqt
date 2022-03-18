package gqt

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/gqt/v2/gqttpl"
)

type Entity struct {
	Hello string
	gqttpl.DataVolumeMap
}

func TestConvert2DataVolume(t *testing.T) {
	entity := &Entity{}

	volume, err := Convert2DataVolume(entity)

	if err != nil {
		panic(err)
	}
	key := "key1"

	volume.SetValue(key, "value1")
	getValue, _ := volume.GetValue(key)
	fmt.Printf("%#v", volume)
	fmt.Printf("%#v", getValue)
}

func TestConvertMap(t *testing.T) {
	dataMap := map[string]interface{}{
		"a": 1,
	}
	interfac := interface{}(dataMap)
	a := gqttpl.DataVolumeMap(dataMap)
	fmt.Printf("%#v----%#v", a, interfac)

}
