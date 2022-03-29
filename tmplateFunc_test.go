package gqt

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/gqt/v2/gqttpl"
)

type Entity struct {
	Hello string
	gqttpl.TplEmptyEntity
}

func TestConvertMap(t *testing.T) {
	dataMap := map[string]interface{}{
		"a": 1,
	}
	interfac := interface{}(dataMap)
	a := gqttpl.TplEmptyEntity(dataMap)
	fmt.Printf("%#v----%#v", a, interfac)

}
