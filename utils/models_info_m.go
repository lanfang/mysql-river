// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

// single model info
type modelInfo struct {
	pkg       string
	name      string
	fullName  string
	table     string
	model     interface{}
	fields    *fields
	addrField reflect.Value
}

// new model info
func newModelInfo(val reflect.Value) (info *modelInfo) {
	var (
		err error
		fi  *fieldInfo
		sf  reflect.StructField
	)

	info = &modelInfo{}
	info.fields = newFields()

	ind := reflect.Indirect(val)
	typ := ind.Type()

	info.addrField = val

	info.name = typ.Name()
	info.fullName = getFullName(typ)

	for i := 0; i < ind.NumField(); i++ {
		field := ind.Field(i)
		sf = ind.Type().Field(i)
		if sf.PkgPath != "" {
			continue
		}
		fi, err = newFieldInfo(info, field, sf)

		if err != nil {
			if err == errSkipField {
				err = nil
				continue
			}
			break
		}

		added := info.fields.Add(fi)
		if added == false {
			err = errors.New(fmt.Sprintf("duplicate column name: %s", fi.column))
			break
		}

		if fi.pk {
			if info.fields.pk != nil {
				err = errors.New(fmt.Sprintf("one model must have one pk field only"))
				break
			} else {
				info.fields.pk = fi
			}
		}

		fi.fieldIndex = i
		fi.mi = info
		fi.inModel = true
	}

	if err != nil {
		fmt.Println(fmt.Errorf("field: %s.%s, %s", ind.Type(), sf.Name, err))
		os.Exit(2)
	}

	return
}
