package utils

import (
	"fmt"
	"os"
	"reflect"
	"sync"
)

func RegisterModel(models ...interface{}) {
	for _, model := range models {
		registModel(model)
	}
}
func registModel(model interface{}) {

	val := reflect.ValueOf(model)
	ind := reflect.Indirect(val)
	typ := ind.Type()

	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("cannot use non-ptr model struct `%s`", getFullName(typ)))
	}

	table := getTableName(val)

	if _, ok := modelCache.get(table); ok {
		fmt.Printf("table name `%s` repeat register, must be unique\n", table)
		os.Exit(2)
	}

	info := newModelInfo(val)

	info.table = table
	info.pkg = typ.PkgPath()
	info.model = model
	modelCache.set(table, info)
}

/*
const (
	od_CASCADE            = "cascade"
	od_SET_NULL           = "set_null"
	od_SET_DEFAULT        = "set_default"
	od_DO_NOTHING         = "do_nothing"
	defaultStructTagName  = "orm"
	defaultStructTagDelim = ";"
)
*/
var (
	modelCache = &_modelCache{
		cache: make(map[string]*modelInfo),
	}
)

// model info collection
type _modelCache struct {
	sync.RWMutex
	cache map[string]*modelInfo
	done  bool
}

// get model info by table name
func (mc *_modelCache) get(table string) (mi *modelInfo, ok bool) {
	mi, ok = mc.cache[table]
	return
}

// set model info to collection
func (mc *_modelCache) set(table string, mi *modelInfo) {
	mc.cache[table] = mi
}

// clean all model info.
func (mc *_modelCache) clean() {
	mc.cache = make(map[string]*modelInfo)
	mc.done = false
}
