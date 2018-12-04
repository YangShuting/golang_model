package storage

import (
	"fmt"

	"gitee.com/firewing_group/blue_kxq2/config"
	mgo "gopkg.in/mgo.v2"
)

var mongoHost string
var selfDBName string
var animalDBName string
var sessionOrigin *mgo.Session

func InitMongo() {
	mongoHost = config.RawEnv("MONGO_HOST")
	selfDBName = config.RawEnv("SELF_DB_NAME")
	animalDBName = config.RawEnv("KETH_DB")
	if session, err := mgo.Dial(mongoHost); err != nil {
		panic("mongo init" + err.Error())
	} else {
		sessionOrigin = session
	}
	fmt.Printf("<<<初始化mongo完成>>>")
}

func GetSession() *mgo.Session {
	return sessionOrigin.Copy()
}
