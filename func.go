package gqt

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/suifengpiao14/gqt/v2/gqttpl"
	"golang.org/x/sync/singleflight"
)

func GetMD5LOWER(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Model2Entity copy model to entity ,some times input used to insert and update ,in this case input mybe model, copy model value to insertEntity and updateEntity
func Model2TplEntity(from interface{}, to gqttpl.TplEntityInterface) {
	err := copier.Copy(to, from)
	if err != nil {
		panic(err)
	}
}

var g = singleflight.Group{}

func Flight(sqlStr string, output interface{}, fn func() (interface{}, error)) (err error) {
	if sqlStr == "" {
		err = errors.New("sql must not be empty")
		return
	}
	output, err, _ = g.Do(GetMD5LOWER(sqlStr), fn)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

//ConvertStruct 转换结构体
func ConvertStruct(from interface{}, to interface{}) {
	err := copier.Copy(to, from)
	if err != nil {
		panic(err)
	}
}
