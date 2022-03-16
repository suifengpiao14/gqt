package gqt

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/gqt/v2/gqttpl"
)

var TemplatefuncMap = template.FuncMap{
	"zeroTime":      ZeroTime,
	"currentTime":   CurrentTime,
	"permanentTime": PermanentTime,
	"contains":      strings.Contains,
	"newPreComma":   NewPreComma,
	"in":            In,
	"toCamel":       gqttpl.ToCamel,
	"toLowerCamel":  gqttpl.ToLowerCamel,
	"snakeCase":     gqttpl.SnakeCase,
}

// DataVolumeInterface 如果模板中需要新增、修改数据，作为 sql prepared statement 值，则传递给模板的数据变量必须实现DataVolumeInterface接口
type DataVolumeInterface interface {
	SetValue(key string, value interface{})
	GetValue(key string) (value interface{}, ok bool)
	GetDynamicValus() (values map[string]interface{})
}

type DataVolumeMap map[string]interface{}

func (v *DataVolumeMap) init() {
	if *v == nil {
		*v = make(map[string]interface{})
	}
}

func (v *DataVolumeMap) SetValue(key string, value interface{}) {
	v.init()
	(*v)[key] = value // todo 并发lock
}

func (v *DataVolumeMap) GetValue(key string) (value interface{}, ok bool) {
	v.init()
	value, ok = (*v)[key]
	return
}

func (v *DataVolumeMap) GetDynamicValus() (values map[string]interface{}) {
	v.init()
	return *v
}

func Convert2DataVolume(data interface{}) (dataVolume DataVolumeInterface, err error) {
	for {
		dataI, ok := data.(*interface{})
		if ok {
			data = *dataI

		} else {
			break
		}
	}
	if dataMap, ok := data.(map[string]interface{}); ok {
		datavolumeMap := DataVolumeMap(dataMap)
		dataVolume = &datavolumeMap
		return
	}
	if dataMap, ok := data.(*map[string]interface{}); ok {
		a := DataVolumeMap(*dataMap)
		dataVolume = &a
		p := (*DataVolumeMap)(unsafe.Pointer(&dataMap))
		data = p
		return
	}

	dataVolume, ok := data.(DataVolumeInterface)
	if !ok {
		err = errors.Errorf("expected implement interface gqt.DataVolume ; got %#v ", data)
		return nil, err
	}

	return
}

func ZeroTime(data interface{}) (string, error) {
	dataVolume, err := Convert2DataVolume(data)
	if err != nil {
		return "", err
	}
	named := "ZeroTime"
	placeholder := ":" + named
	value := "0000-00-00 00:00:00"
	dataVolume.SetValue(named, value)
	return placeholder, nil
}

func CurrentTime(data interface{}) (string, error) {
	dataVolume, err := Convert2DataVolume(data)
	if err != nil {
		return "", err
	}
	named := "CurrentTime"
	placeholder := ":" + named
	value := time.Now().Format("2006-01-02 15:04:05")
	dataVolume.SetValue(named, value)
	return placeholder, nil
}

func PermanentTime(data interface{}) (string, error) {
	dataVolume, err := Convert2DataVolume(data)
	if err != nil {
		return "", err
	}
	named := "PermanentTime"
	placeholder := ":" + named
	value := "3000-12-31 23:59:59"
	dataVolume.SetValue(named, value)
	return placeholder, nil
}

type preComma struct {
	comma string
}

func NewPreComma() *preComma {
	return &preComma{}
}

func (c *preComma) PreComma() string {
	out := c.comma
	c.comma = ","
	return out
}

func In(obj interface{}, data interface{}) (str string, err error) {
	dataVolume, err := Convert2DataVolume(obj)
	if err != nil {
		return "", err
	}
	placeholders := make([]string, 0)
	inIndexKey := "InIndex_"
	inIndex := 0
	inIndexInterface, _ := dataVolume.GetValue(inIndexKey)
	if inIndexInterface != nil {
		inIndexInt, ok := inIndexInterface.(int)
		if ok {
			inIndex = inIndexInt
		}
	}

	v := reflect.Indirect(reflect.ValueOf(data))

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		num := v.Len()
		for i := 0; i < num; i++ {
			inIndex++
			named := fmt.Sprintf("in_%d", inIndex)
			placeholder := ":" + named
			placeholders = append(placeholders, placeholder)
			dataVolume.SetValue(named, v.Index(i).Interface())
		}

	case reflect.String:
		arr := strings.Split(v.String(), ",")
		num := len(arr)
		for i := 0; i < num; i++ {
			inIndex++
			named := fmt.Sprintf("in_%d", inIndex)
			placeholder := ":" + named
			placeholders = append(placeholders, placeholder)
			dataVolume.SetValue(named, arr[i])
		}
	default:
		err = fmt.Errorf("want slice/array/string ,have %s", v.Kind().String())
		if err != nil {
			return "", err
		}
	}
	dataVolume.SetValue(inIndexKey, inIndex) // 更新InIndex_
	str = strings.Join(placeholders, ",")
	return str, nil

}
