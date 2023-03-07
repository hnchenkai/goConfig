package goConfig

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"sync"

	"github.com/hnchenkai/goConfig/goConfig/store"
)

// 配置列表类
type ConfigUnit struct {
	list []store.IFConfigBase
	lock *sync.RWMutex
	// 记录备份的配置信息
	fastUnit       reflect.Value
	configBaseType reflect.Type
}

// 初始化结构
func (p *ConfigUnit) Init() {
	p.list = make([]store.IFConfigBase, 0)
	p.lock = &sync.RWMutex{}
}

// 删除一个元素
func (p *ConfigUnit) del(Type int) {
	var pIdx int = -1
	for idx, p1 := range p.list {
		if p1.GetType() == Type {
			pIdx = idx
			break
		}
	}
	if pIdx >= 0 {
		p.lock.Lock()
		if pIdx == len(p.list)-1 {
			p.list = p.list[:pIdx]
		} else {
			p.list = append(p.list[:pIdx], p.list[pIdx+1])
		}
		p.lock.Unlock()
	}
}

// 添加一个插入操作
func (p *ConfigUnit) Append(Type int, cu store.IFConfigBase, keys []string) {
	defer p.freshValidClone()

	cu.SetKeys(keys)
	cu.SetType(Type)
	cu.Init()
	p.del(cu.GetType())
	p.lock.Lock()
	defer p.lock.Unlock()
	p.list = append(p.list, cu)
	sort.Slice(p.list, func(i, j int) bool {
		return p.list[i].GetType() > p.list[j].GetType()
	})
	// 这里表示配置出现了变动，可以考虑刷新一个配置的快速信息表
}

func (p *ConfigUnit) clone() []store.IFConfigBase {
	p.lock.RLock()
	defer p.lock.RUnlock()
	listCopy := make([]store.IFConfigBase, len(p.list))
	copy(listCopy, p.list)

	return listCopy
}

func isBoolTrue(isOk []reflect.Value) bool {
	if len(isOk) == 0 || isOk[0].Kind() != reflect.Bool || !isOk[0].Bool() {
		return false
	} else {
		return true
	}
}

// 找到IsValid的方法
func findMethod(va reflect.Value) reflect.Value {
	if !va.IsValid() {
		return reflect.Value{}
	}
	method := va.MethodByName("IsValid")
	if method.IsValid() {
		return method
	}

	if va.Kind() != reflect.Ptr {
		va = va.Addr()
	}

	method = va.MethodByName("IsValid")
	if method.IsValid() {
		return method
	}

	return reflect.Value{}
}

func (p *ConfigUnit) GetInfo(key string, args ...interface{}) (result reflect.Value, ok bool) {
	ok = false
	numIn := len(args)
	argsIn := make([]reflect.Value, 0)
	for _, v := range args {
		argsIn = append(argsIn, reflect.ValueOf(v))
	}
	list := p.clone()
	for _, pice := range list {
		valueT := reflect.ValueOf(pice)
		if !pice.HasKey(key) {
			continue
		}
		// fmt.Println(valueT.Elem().Type(), valueT.NumMethod())
		va := valueT.Elem().FieldByName(key)
		method := findMethod(va)
		if method.IsValid() && method.Type().NumIn() == numIn && !isBoolTrue(method.Call(argsIn)) {
			continue
		}
		result = va
		ok = true
		break
	}

	return
}

// 以下是获得数据
func GetInfo[T any](app *ConfigUnit, key string, args ...interface{}) (result T, ok bool) {
	value, ok := app.GetInfo(key, args...)
	if !ok {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			ok = false
		}
	}()
	return value.Interface().(T), true
}

// 导出一份正确的数据
func (p *ConfigUnit) freshValidClone() {
	// 构建一个新的对象
	newUt := reflect.New(p.configBaseType)
	for i := 0; i < p.configBaseType.NumField(); i++ {
		fd := p.configBaseType.Field(i)
		// 私有数据我们就pass了
		if !fd.IsExported() {
			continue
		}

		if uu, ok := p.GetInfo(fd.Name); ok {
			newUt.Elem().FieldByName(fd.Name).Set(uu)
		}
	}
	// 赋值上去
	p.fastUnit = newUt

	log.Println("freshValidClone")
}

func (p *ConfigUnit) SetConfigStruct(t reflect.Type) {
	p.configBaseType = t
}
func (p *ConfigUnit) NewConfigBase() store.IFConfigBase {
	return reflect.New(p.configBaseType).Interface().(store.IFConfigBase)
}

// 导出一份正确的数据
func (p *ConfigUnit) Static() (result interface{}) {
	if !p.fastUnit.IsValid() {
		p.freshValidClone()
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			result = nil
		}
	}()

	result = p.fastUnit.Interface()
	return
}

func (p *ConfigUnit) Get(key string) (result interface{}) {
	ut := p.fastUnit.Elem().FieldByName(key)
	if !ut.IsValid() {
		return nil
	}
	return ut.Interface()
}

func (p *ConfigUnit) GetString(key string) (result string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			result = ""
		}
	}()
	ut := p.fastUnit.Elem().FieldByName(key)
	if !ut.IsValid() {
		return ""
	}
	if ut.Type().String() == "string" {
		fmt.Println(ut.Type(), ut.Kind())
		return ut.Interface().(string)
	}
	uu := ut.Interface().(StringValue)
	return uu.Value()
}

func (p *ConfigUnit) GetOrgin(key string, args ...interface{}) interface{} {
	result, ok := GetInfo[interface{}](p, key, args...)
	if ok {
		return result
	} else {
		return nil
	}
}
