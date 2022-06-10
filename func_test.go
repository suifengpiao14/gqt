package gqt

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFlight(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			nubmer1 := 0
			callFlight(&nubmer1)
			fmt.Println(nubmer1)
		}()
	}

}

func callFlight(number interface{}) {
	Flight("aa", number, func() (interface{}, error) {
		rv := reflect.Indirect(reflect.ValueOf(number))
		rv.SetInt(10)
		fmt.Println("hello world")
		return number, nil
	})
}
