package utils

import (
	"os"
	"strconv"
	"strings"
)

var envCache map[string]string = map[string]string{}
var envKeys []string

// 初始化环境变量信息
func init() {
	for _, v := range os.Environ() {
		keys := strings.Split(v, "=")
		key := strings.ToLower(strings.ReplaceAll(keys[0], ".", "_"))
		value := strings.Join(keys[1:], "=")
		envCache[key] = value
	}
	envKeys = make([]string, 0, len(envCache))
	for k := range envCache {
		envKeys = append(envKeys, k)
	}
}

// 返回所有环境变量的key
func GetEnvKeys() []string {
	return envKeys
}

// 从环境变量中获取一下对应的信息
// 环境变量名字支持 xx.xxx  xx_xxx 两者等价
func GetEnv(key string) string {
	key = strings.ToLower(strings.ReplaceAll(key, ".", "_"))
	v, ok := envCache[key]
	if !ok {
		return ""
	}
	return v
}

// 转换成数字返回
func GetEnvInt64(key string) int64 {
	v, err := strconv.ParseInt(GetEnv(key), 10, 64)
	if err != nil {
		return 0
	}

	return v
}

// 转换成数字返回
func GetEnvInt(key string) int {
	v, err := strconv.Atoi(GetEnv(key))
	if err != nil {
		return 0
	}

	return v
}

// 转换成浮点数返回
func GetEnvFloat(key string) float64 {
	v, err := strconv.ParseFloat(GetEnv(key), 64)
	if err != nil {
		return 0
	}

	return v
}

// 返回值 true :内容是true或者1的时候
func GetEnvBool(key string) bool {
	v := GetEnv(key)
	if v == "true" || v == "1" {
		return true
	} else {
		return false
	}
}
