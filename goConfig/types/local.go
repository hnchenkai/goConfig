package types

import (
	"reflect"
	"time"

	"github.com/hnchenkai/goConfig/goConfig/store"

	"github.com/jinzhu/configor"
)

type localStore struct {
	dataT   reflect.Type
	appendF store.AppendCallback
}

func NewLocalStore(file string, tp reflect.Type, append store.AppendCallback) *localStore {
	uu := localStore{
		dataT:   tp,
		appendF: append,
	}
	uu.initLocal(file)
	return &uu
}

func (w *localStore) initLocal(file string) {
	localConfig := reflect.New(w.dataT).Interface()
	// 配置文件没10s自动更新
	wt := configor.New(&configor.Config{AutoReload: true,
		AutoReloadInterval: 10 * time.Second,
		Silent:             true,
		AutoReloadCallback: func(config interface{}) {
			w.appendF(store.Local, config.(store.IFConfigBase), []string{})
		}})
	wt.Load(localConfig, file)
	w.appendF(store.Local, localConfig.(store.IFConfigBase), []string{})
}
