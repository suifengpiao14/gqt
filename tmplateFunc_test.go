package gqt

import (
	"fmt"
	"testing"
)

type Entity struct {
	Hello string
	TplEmptyEntity
}

func TestConvertMap(t *testing.T) {
	dataMap := map[string]interface{}{
		"a": 1,
	}
	interfac := interface{}(dataMap)
	a := TplEmptyEntity(dataMap)
	fmt.Printf("%#v----%#v", a, interfac)

}
