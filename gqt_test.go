// Copyright 2016 Davide Muzzarelli. All right reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gqt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var testDir string

func init() {
	var err error
	testDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	testDir = filepath.Join(testDir, "test")
}

func Test(t *testing.T) {
	for _, dir := range []string{"pkg1", "pkg2"} {
		err := Add(filepath.Join(testDir, dir), TemplatefuncMap)
		if err != nil {
			t.Error(err)
		}
	}

	data := make(map[string]interface{})
	data["APIID"] = "1548452447"
	var sql string
	err := GetSafeSQL("parameter.getAllByAPIID", data, &sql)
	if err != nil {
		panic(err)
	}
	fmt.Printf(sql)

}
