package types

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"

	"github.com/hnchenkai/goConfig/goConfig/store"

	"github.com/hnchenkai/goConfig/goConfig/namings"

	"github.com/nacos-group/nacos-sdk-go/vo"
)

type nacosWatch struct {
	dataT   reflect.Type
	appendF store.AppendCallback
}

func NewNacosWatch(dataId string, groupId string, dataT reflect.Type, append store.AppendCallback) *nacosWatch {
	uu := nacosWatch{
		dataT:   dataT,
		appendF: append,
	}

	uu.InitNacos(dataId, groupId)

	return &uu
}

// 初始化nacos操作
func (w *nacosWatch) InitNacos(dataId string, groupId string) {
	configClient := namings.GetNacosClient()
	if configClient == nil {
		return
	}
	// dataId := os.Getenv("nacos.file")
	// groupId := "DEFAULT_GROUP"
	//监听配置
	if err := configClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  groupId,
		OnChange: func(namespace, group, dataId, data string) {
			// log.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
			w.parseConfig(data)
		},
	}); err != nil {
		log.Printf("err: %v\n", err)
	}

	if data, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  groupId,
	}); err == nil {
		w.parseConfig(data)
		log.Println("nacos 配置加载成功", len(data))
	} else {
		log.Println("nacos 配置加载失败", err)
	}

	// time.Sleep(time.Second * 1000)
}

func (w *nacosWatch) parseConfig(data string) {
	cu := reflect.New(w.dataT).Interface()
	if err := json.Unmarshal([]byte(data), cu); err != nil {
		log.Printf("err: %+v\n", err)
	} else {
		mp := map[string]interface{}{}
		if err := json.Unmarshal([]byte(data), &mp); err != nil {
			log.Printf("err: %+v\n", err)
		}
		keys := make([]string, 0)
		for k := range mp {
			keys = append(keys, strings.ToLower(k))
		}
		w.appendF(store.Nacos, cu.(store.IFConfigBase), keys)
	}
}
