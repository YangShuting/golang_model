package config

import (
	"fmt"
	"os"
	"sync"

	ini "gopkg.in/ini.v1"
)

type KV struct {
	V  string
	ts int64
}

var configLock sync.Mutex

type ConfigMap map[string]*KV

var configMap = make(ConfigMap)

func LoadEnv() string {
	getEnv := os.Getenv("mode")
	if getEnv == "" {
		return "debug"
	}
	return getEnv
}

func Init() {
	env := LoadEnv()
	fmt.Printf("<<<正在加载环境配置 %v\n", env)
	dataByte, err := ini.Load("./config/" + env + ".ini")
	if err != nil {
		panic(err)
	}
	for _, key := range dataByte.Section("").Keys() {
		newkv := new(KV)
		newkv.V = key.String()
		configMap[key.Name()] = newkv
	}
	for k, kv := range configMap {
		fmt.Printf("<<<configenv: %v: %v>>> \n", k, kv.V)
	}
}

func RawEnv(k string) string {
	configLock.Lock()
	defer configLock.Unlock()
	var found string
	if v, found := configMap[found]; found {
		return v.V
	} else {
		if vv, vFound := configMap[k]; vFound {
			return vv.V
		}
		return ""
	}
}
