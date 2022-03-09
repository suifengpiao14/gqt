package gqt

import (
	"fmt"
	"testing"
)

func TestGetMapData(t *testing.T) {
	dataVolume := &map[string]interface{}{
		"InIndex_1": 10,
	}
	inIndexKey := "InIndex_"
	v := GetMapData(dataVolume, inIndexKey)
	fmt.Println(v)
}
