# goConfig

提供 go 语言环境下动态获取配置信息的能力，支持本地配置文件和 nacos 远程加载

## 配置说明

支持自定义配置结构提，提供安全快速的配置加载方案

## 如何定义结构

配置文件支持一层的合法值判定机制  
使用的时候需要通过符合 IFConfig 结构的类型作为元素

```
type IFConfig interface {
	IsValid() bool
	Value() interface{}
}
```

没有定义 IsValid 方法，或者传入参数，返回结果类型错误的，都会导致合法判定失败。  
默认优先级 nacos > file (配置中不存在也会判定存在)

配置文件的 struct 需要继承 ConfigBase

## 获取环境变量

内置的 utils 中支持获取环境变量中的配置项  
环境变量 key 中 . 等价于 \_ 即 http.port 等价于 http_port

## 启用 nacos 说明

支持 nacos 连接，需要注入环境变量

1. nacos_serveraddr 服务器地址 不配置地址的话无法启用 nacos
1. nacos_serverport 服务器端口 默认端口 8848
1. nacos_namespace 命名空间
1. nacos_user 用户名
1. nacos_passwd 密码
1. nacos_logdir 日志导出目录默认 tmp/nacos/log
1. nacos_cachedir 配置缓存目录默认 tmp/nacos/cache
1. nacos_loglevel 日志等级默认 info

## 使用案例

```
type ConfigUnits struct {
	goConfig.ConfigBase
	Http     goConfig.HttpValue   // 启动的端口
	Log      goConfig.DebugLevel  `json:"log"` // "info" 模式下会显示接口的信息
	PrometheusGate goConfig.StringValue  `json:"prometheusgate"`
	ShopLimit      goConfig.StringValue  // 开启商户的隔离检查
}

func TestStruct(t *testing.T) {
	inst := goConfig.Default(reflect.TypeOf(ConfigUnits{}))
	// 顺序，越前面优先级越高

	for {
		cl := inst.Static().(*ConfigUnits)
		fmt.Println("port", inst.Get("Http").(goConfig.HttpValue).Port, cl.Http.Port)
		fmt.Println("PrometheusGate", inst.GetString("JwtSecret"))
		time.Sleep(time.Second * 2)
	}
}
```
