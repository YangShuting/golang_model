package model

import (
	"gitee.com/firewing_group/blue_kxq2/config"
)

var DBNAME_KuaiMaoUser = ""

func InitModel() {
	DBNAME_KuaiMaoUser = config.RawEnv("KETH_DB")
}
