package namings

import (
	"goConfig/goConfig/utils"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosConfig struct {
	Namespace string
	Addr      string
	Port      uint64
	Passwd    string
	Username  string
	LogDir    string
	CacheDir  string
	LogLevel  string
}

func getEvnNacosPort() uint64 {
	port := utils.GetEnv("nacos.serverPort")
	if len(port) == 0 {
		return 8848
	}
	uPort, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		return 8848
	}
	return uPort
}

func getNacosConfig() *NacosConfig {
	nacosConfig := NacosConfig{
		Namespace: utils.GetEnv("nacos.namespace"),
		Addr:      utils.GetEnv("nacos.serverAddr"),
		Port:      getEvnNacosPort(),
		Username:  utils.GetEnv("nacos.user"),
		Passwd:    utils.GetEnv("nacos.passwd"),
	}
	LogDir := utils.GetEnv("nacos.logdir")
	if len(LogDir) == 0 {
		LogDir = "tmp/nacos/log"
	}
	nacosConfig.LogDir = LogDir

	CacheDir := utils.GetEnv("nacos.cachedir")
	if len(CacheDir) == 0 {
		CacheDir = "tmp/nacos/cache"
	}
	nacosConfig.CacheDir = CacheDir

	LogLevel := utils.GetEnv("nacos.loglevel")
	if len(CacheDir) == 0 {
		LogLevel = "info"
	}
	nacosConfig.LogLevel = LogLevel

	return &nacosConfig
}

var configClient config_client.IConfigClient

func GetNacosClient() config_client.IConfigClient {
	return configClient
}

func init() {
	nacosConfig := getNacosConfig()
	if len(nacosConfig.Addr) == 0 {
		return
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacosConfig.Addr,
			Port:   nacosConfig.Port,
		},
	}

	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         nacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              nacosConfig.LogDir,
		CacheDir:            nacosConfig.CacheDir,
		Username:            nacosConfig.Username,
		Password:            nacosConfig.Passwd,
		LogLevel:            nacosConfig.LogLevel,
	}
	// 创建动态配置客户端的另一种方式 (推荐)
	c, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	configClient = c
}
