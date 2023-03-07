package store

import "strings"

const (
	Local = 1
	Nacos = 2
)

type IFConfig interface {
	IsValid() bool
	Value() interface{}
}

type IFConfigBase interface {
	Init()
	SetType(Type int)
	GetType() int
	SetKeys(keys []string)
	HasKey(key string) bool
}

// 配置单元，新增配置需要
type ConfigBase struct {
	__type__ int // 配置类型，这个是设置的，不是来自配置文件
	__keys__ []string
}

func (c *ConfigBase) SetType(Type int) {
	c.__type__ = Type
}

func (c *ConfigBase) SetKeys(keys []string) {
	c.__keys__ = keys
}

func (c *ConfigBase) HasKey(key string) bool {
	if c.__type__ == Local {
		return true
	}
	key = strings.ToLower(key)
	for _, v := range c.__keys__ {
		if v == key {
			return true
		}
	}
	return false
}

func (c *ConfigBase) GetType() int {
	return c.__type__
}

func (c *ConfigBase) Init() {

}
