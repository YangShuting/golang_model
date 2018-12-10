package model

import (
	"github.com/Yangshuting/golang_model/lib"
	"github.com/Yangshuting/golang_model/storage"
	"gopkg.in/mgo.v2/bson"
)

const CNAME_CLICK = "count_click"

type Click struct {
	ID  string `bson:"_id" json:"_id"`
	Num string `bson:"num" json:"num"`
}

func NewClickCount(cc *lib.Cusctx, id string) *Click {
	click := &Click{
		ID:  id,
		Num: "1",
	}

	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_CLICK).Insert(&click)
	return click
}

func FindAClick(cc *lib.Cusctx, id string) (*Click, error) {
	var click *Click
	findQ := bson.M{
		"_id": id,
	}
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_CLICK).Find(findQ).One(&click)
	if err != nil {
		return nil, err
	}
	return click, nil
}

func IncClick(cc *lib.Cusctx, id string) (int64, error) {
	click, findErr := FindAClick(cc, id)
	if findErr != nil {
		return 0, findErr
	}
	updateC := bson.M{
		"$set": bson.M{
			"num": storage.StringToInt(click.Num) + 1,
		},
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_CLICK).UpdateId(click.ID, updateC)
	return int64(storage.StringToInt(click.Num) + 1), nil
}
