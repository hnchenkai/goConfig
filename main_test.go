package mainT_test

import (
	"fmt"
	"goConfig/goConfig"
	"reflect"
	"testing"
	"time"
)

type httpUnit struct {
	Port int
}

type UserUnit struct {
	Port int
}

func (p *UserUnit) IsValid(aaa interface{}) bool {
	return true
}

type ConfigUnits struct {
	goConfig.ConfigBase
	Http    httpUnit // 启动的端口
	Metaid  string   // 消息队列的地址
	Checkid string
	User    UserUnit
	// Log      goConfig.DebugLevel  `json:"log"` // "info" 模式下会显示接口的信息
	// LogKit   goConfig.LogkitValue `json:"logMgr"`
	// IsvCheck goConfig.StringValue `json:"isvcheck"` // 是否关闭isv检查
	// // 服务配置信息
	// Redis          goConfig.RedisOption  `json:"redis"` // redis的地址
	// Amqp           goConfig.StringValue  // 密钥信息是本地缓存的，监听mq然后刷新缓存
	// Passwds        goConfig.IPasswdValue // 定制项目固定的账号密码 这个是临时操作，适合本地化的项目
	// JwtSecret      goConfig.StringValue  `json:"jwtsecret"` // jwt签名使用的密钥
	// Beta           goConfig.BetaValue    // 活动服务的灰度配置
	// BetaApp        goConfig.BetaValue    // 渠道服务的灰度配置
	// SystemUser     goConfig.StringValue  `json:"system"` // 系统登录的白名单
	// PrometheusGate goConfig.StringValue  `json:"prometheusgate"`
	// ShopLimit      goConfig.StringValue  // 开启商户的隔离检查
}

func TestStruct(t *testing.T) {
	inst := goConfig.Default(reflect.TypeOf(ConfigUnits{}))
	// 顺序，越前面优先级越高

	for {
		// cl := inst.Static().(*ConfigUnits)
		httpV := inst.Get("Http")
		if httpV != nil {
			fmt.Println("port", httpV.(httpUnit).Port)
		}
		user := inst.GetOrgin("User", "123s")
		if user != nil {
			fmt.Println("User", user)
		}
		fmt.Println("Metaid", inst.GetString("Metaid"))
		fmt.Println("Checkid", inst.GetString("Checkid"))
		time.Sleep(time.Second * 2)
	}

	// for i := 0; i < 10000; i++ {
	// 	time.Sleep(time.Second)
	// }
}
