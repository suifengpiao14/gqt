package gqt

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/gqt/v2/gqttpl"
)

type ExampleModel struct {
	ExampleID string
	ServiceID int
	Request   bool
	// 创建时间
	CreatedAt string `json:"createdAt"`

	// 修改时间
	UpdatedAt string `json:"updatedAt"`

	// 删除时间
	DeletedAt string `json:"deletedAt"`
}

type ExampleEntity struct {
	ExampleID string
	ServiceID int
	Request   bool
	gqttpl.DataVolumeMap
}

func (e *ExampleEntity) TplName() string {
	return "example"
}
func (e *ExampleEntity) TplOutput() (string, error) {
	return "tpl exec output", nil
}

func TestModel2Entity(t *testing.T) {
	from := &ExampleModel{
		ExampleID: "aa",
		ServiceID: 1,
		Request:   true,
		CreatedAt: "2022-03-06",
	}
	to := &ExampleEntity{}
	Model2TplEntity(from, to)
	fmt.Println(to)
}
