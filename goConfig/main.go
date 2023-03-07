package goConfig

import (
	"os"
	"reflect"
	"sync"

	"github.com/hnchenkai/goConfig/goConfig/store"
	"github.com/hnchenkai/goConfig/goConfig/types"
)

type ConfigBase struct {
	store.ConfigBase
}

// 返回一个自定义文件地址，nacos分组等信息的对象
func NewInstance(structT reflect.Type, file string, dataId string, groupId string) *ConfigUnit {
	ut := &ConfigUnit{
		list:           make([]store.IFConfigBase, 0),
		lock:           &sync.RWMutex{},
		fastUnit:       reflect.Value{},
		configBaseType: structT,
	}
	types.NewLocalStore(file, structT, ut.Append)
	types.NewNacosWatch(dataId, groupId, structT, ut.Append)
	return ut
}

// 返回一个默认设置的对象
func Default(structT reflect.Type) *ConfigUnit {
	dataId := os.Getenv("nacos.file")
	groupId := "DEFAULT_GROUP"
	return NewInstance(structT, "./config.json", dataId, groupId)
}
