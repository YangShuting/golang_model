package test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"gitee.com/firewing_group/blue_kxq2/model"
	"gitee.com/firewing_group/blue_kxq2/storage"
	"github.com/stretchr/testify/assert"

	ini "gopkg.in/ini.v1"
	resty "gopkg.in/resty.v1"
)

type KV struct {
	V  string
	ts int64
}

var config_lock sync.Mutex
var configMap map[string]*KV

func TestConfigEnv(t *testing.T) {
	var configByte = []byte(
		`ConfigMongoDB=global_config
		ConfigMongoCollection=blue_kxq
		testVariable=keth
		`,
	)
	fmt.Printf("configenv=> %+v \n", configByte)
	iniLoadData, err := ini.Load(configByte)
	fmt.Printf("configByte %+v \n", iniLoadData)
	fmt.Printf("read err: %+v \n", err)
	for _, key := range iniLoadData.Section("").Keys() {
		_newkv := new(KV)
		_newkv.V = key.String()
		fmt.Printf("_newkv.V => %+v", _newkv.V)
		configMap[key.Name()] = _newkv
		// _newkv := new(KV)
		// _newkv.V = key.String()
		// configmap[key.Name()] = _newkv
		// fmt.Printf("_newkv=> %+v", _newkv)
		// fmt.Printf(" \n key => %+v \n", key)
		// fmt.Printf(" \n key.Name() => %+v \n", key.Name())
	}
	for k, kv := range configMap {
		fmt.Printf("cconfigenv, k=> %+v", k)
		fmt.Printf("cconfigenv, kv=> %+v", kv)
	}
}

func TestIniLoadData(t *testing.T) {
	// env := os.Getenv("mode")
	// dataByte, err := ini.Load(fmt.Sprintf("../config/", env))
	// fmt.Printf("err is => %+v", err)
	// fmt.Printf("datByte => %+v", dataByte)
}

func TestContextParam(t *testing.T) {
	// t := reflect.ValueOf(TestIniLoadData)
	// fmt.Printf("t_%+v", t)
}

//测试 lib 中的资源池
func TestPool(t *testing.T) {

}

func TestHello(t *testing.T) {
	resp, err := resty.R().
		Post("http://127.0.0.1:1344/blue_kxq2/hello")
	if err != nil {
		assert.Errorf(t, err, "post error")
	}
	fmt.Printf("resp_%+v", resp)
}

func TestGetUserSuanli(t *testing.T) {
	resp, err := resty.R().
		Post("http://127.0.0.1:1344/blue_kxq2/getUserSuanli")
	if err != nil {
		assert.Errorf(t, err, "get error")
	}
	fmt.Printf("resp_%+v", resp.Body())
}

func TestRedis(t *testing.T) {
	//set
	err := storage.SetRedis("testRedisYST", "YesItPass")
	fmt.Printf("set_redis_err_%+v", err)

	//get
	val, err := storage.GetRedis("testRedisYST")
	fmt.Printf("get_redis_%+v", val)
}

func TestLogin(t *testing.T) {
	resp, err := resty.R().
		Post("http://127.0.0.1:1344/blue_kxq2/login")
	if err != nil {
		assert.Errorf(t, err, "get error")
	}
	fmt.Printf("resp_%+v", resp.Body())
}

var AUTO_NUM int = 0

func TestCAS(t *testing.T) {
	flg := make(chan int)
	go autoIn(flg)
	go autoIn2(flg)
	<-flg
	fmt.Printf("auto_num_%+v", AUTO_NUM)
}
func autoIn(f chan int) {
	for i := 0; i < 100; i++ {
		AUTO_NUM = AUTO_NUM + 1
	}
	fmt.Println("autoIn")
}

func autoIn2(f chan int) {
	for i := 0; i < 100; i++ {
		AUTO_NUM = AUTO_NUM + 1
	}
	fmt.Println("autoIn2")

	f <- 1
}

func producer(c chan int) {
	defer close(c)
	for i := 0; i < 10; i++ {
		c <- i
	}
}

func consumer(c, f chan int) {
	for {
		if v, ok := <-c; ok {
			fmt.Print(v)
		} else {
			break
		}
	}
	f <- 1
}

func sayHello() {
	for i := 0; i < 10; i++ {
		fmt.Printf("hello \n")
		runtime.Gosched()
	}
}

func sayWorld() {
	for i := 0; i < 10; i++ {
		fmt.Printf("world \n")
		runtime.Gosched()
	}
}

func TestSpeedLimiter(t *testing.T) {
	bool, err := model.SpeedLimiter("test_speed_limiter_@", 30)
	fmt.Printf("ifSuccess_%+v", bool)
	fmt.Printf("err_%+v", err)
}
