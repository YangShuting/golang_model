package model

import "github.com/Yangshuting/golang_model/config"

var DBNAME_KuaiMaoUser = ""

func InitModel() {
	DBNAME_KuaiMaoUser = config.RawEnv("KETH_DB")
}
