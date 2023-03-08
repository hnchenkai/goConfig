package goConfig

type RedisOption struct {
	Network string
	// host:port address.
	Addr string

	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string

	// Database to be selected after connecting to the server.
	DB int
}
type IPasswd struct {
	// 客户id
	Appid string
	// 账号名字
	Account string
	// 登录用的密码
	Passwd string
	// 手机号，主要给发送验证码使用
	Phone string
}

type DebugLevel int32

const (
	Unknow DebugLevel = 0
	Error  DebugLevel = 1
	Warn   DebugLevel = 2
	Info   DebugLevel = 3 // 开启接口请求信息打印
)

func (c *DebugLevel) ShowApiLog() bool {
	return *c > Warn
}

// 判断是否有效
func (c *DebugLevel) IsValid() bool {
	return *c > Unknow
}

type IPasswdValue []IPasswd

func (c *IPasswdValue) IsValid() bool {
	return len(*c) > 0
}

// 定义一下有效字符串的
func isValidString(v string) bool {
	if v == "" {
		return false
	}

	if v == "null" {
		return false
	}

	if v == "nil" {
		return false
	}

	return true
}

// 字符串的累
type StringValue string

// 判断是否有效
func (c *StringValue) IsValid() bool {
	return isValidString(c.Value())
}

func (c *StringValue) Value() string {
	return string(*c)
}

type HttpValue struct {
	Port int
}

func (c *HttpValue) IsValid() bool {
	return c.Port > 0
}

func (c *RedisOption) IsValid() bool {
	return isValidString(c.Addr)
}

func (c *RedisOption) Value() *RedisOption {
	return c
}

type IntValue int

func (c *IntValue) IsValid() bool {
	return *c != 0
}

// 灰度规则
type BetaValue map[string]string

func (c *BetaValue) IsValid(activeId string) bool {
	if c == nil {
		return false
	}
	_, ok := (*c)[activeId]
	return ok
}

func (c *BetaValue) Value(activeId string) string {
	ver, ok := (*c)[activeId]
	if !ok {
		return ""
	}
	return ver
}

type LogkitValue struct {
	Platform  string `json:"platform"`
	Url       string `json:"url"`
	ProjectId string `json:"projectid"`
	Path      string `json:"path"`
}

func (c *LogkitValue) IsValid() bool {
	return len(c.Url) > 0
}

func (c *LogkitValue) Value() *LogkitValue {
	return c
}

// 使用字符串 "true" 才表示true
type StringBoolValue string

func (c *StringBoolValue) IsValid() bool {
	return len(*c) > 0
}

func (c *StringBoolValue) Value() bool {
	return *c == "true"
}
