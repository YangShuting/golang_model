package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/Yangshuting/golang_model/lib"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const CNAME_KuaiMaoUser = "users"
const (
	EMAILSTATUS_NOTSET    = "NOTSET"
	EMAILSTATUS_NOTVERIFY = "NOTVERIFY"
	EMAILSTATUS_GOOD      = "good"
)

type Email struct {
	Addr   string `bson:"addr" json:"addr"`
	Status string `bson:"status" json:"status"`
}
type FillCode struct {
	Wechat bool `bson:"wechat" json:"wechat"`
	QQ     bool `bson:"qq" json:"qq"`
	Kabao  bool `bson:"kabao" json:"kakao"`
}
type KuaiMaoUser struct {
	ID             bson.ObjectId `bson:"_id" json:"_id"`
	UserName       int64         `bson:"username" json:"username"`
	NiCheng        string        `bson:"nicheng" json:"nicheng"`
	Index          int64         `bson:"index" json:"index"` // 第几个人物
	Telephone      string        `bson:"telephone" json:"-"`
	Icode          string        `bson:"icode" json:"-"`         // myself icode
	FIcode         string        `bson:"ficode" json:"whoyaowo"` // who invites me
	MyAvailableSub int64         `bson:"mysub" json:"-"`         // 我的一级下线人数还有几个空位
	CreateTime     int64         `bson:"createtime" json:"createtime"`
	Suanli         int64         `bson:"suanli" json:"suanli"`
	B              int64         `bson:"b" json:"b"`
	LastLoginAt    int64         `bson:"lla" json:"lla"` // 上次上线时间
	Email          *Email        `bson:"email" json:"-"`
	FillCode       *FillCode     `bson:"fillcode" json:"fillcode"` // 微信和qq验证码是否已经填过了
	WechatID       string        `bson:"wechatid" json:"-"`
	Country        string        `bson:"country" json:"-"`
	Device         string        `bson:"device" json:"-"`
	ChaosticID     string        `bson:"chaosticid" json:"chaosticid"`
	Channel        string        `bson:"channel" json:"-"`
	updater        bson.M
}

func NewKuaiMaoUser() *KuaiMaoUser {
	var kuaimaouser KuaiMaoUser
	kuaimaouser.Email = new(Email)
	kuaimaouser.Email.Status = EMAILSTATUS_NOTSET
	kuaimaouser.FillCode = new(FillCode)
	return &kuaimaouser
}

func (user *KuaiMaoUser) Insert(cc *lib.Cusctx) error {
	user.ID = bson.NewObjectId()
	user.CreateTime = time.Now().Unix()
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Insert(user)
	return err
}
func (user *KuaiMaoUser) SetChaosticID(cc *lib.Cusctx, chaosticID string) {
	user.ChaosticID = chaosticID
	change := bson.M{
		"$set": bson.M{
			"chaosticid": user.ChaosticID,
		},
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).UpdateId(user.ID, change)
}

func FindByTelephone(cc *lib.Cusctx, telephone string) *KuaiMaoUser {
	cont := bson.M{
		"telephone": telephone,
	}
	var user KuaiMaoUser
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).One(&user)
	if user.CreateTime == 0 {
		cc.Errf("stuser.FindByTelephone(),telephone=%v,err=%v", telephone, err)
		return nil
	}
	return &user
}

func FindByUserName(cc *lib.Cusctx, username int64) *KuaiMaoUser {
	cont := bson.M{
		"username": username,
	}
	var user KuaiMaoUser
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).One(&user)
	if user.CreateTime == 0 {
		cc.Errf("stuser.FindByUserName,username=%v,err=%v", username, err)
		return nil
	}
	return &user
}
func UIDhexByWechatID(cc *lib.Cusctx, WechatID string) string {
	cont := bson.M{
		"wechatid": WechatID,
	}
	selector := bson.M{
		"_id": 1,
	}
	var _user KuaiMaoUser
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).Select(selector).One(&_user)
	if err != nil {
		return ""
	}
	return _user.ID.Hex()
}
func WechatidByWechatID(cc *lib.Cusctx, WechatID string) string {
	cont := bson.M{
		"wechatid": WechatID,
	}
	selector := bson.M{
		"wechatid": 1,
	}
	var _user KuaiMaoUser
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).Select(selector).One(&_user)
	if err != nil {
		return ""
	}
	return WechatID
}
func FindByWechatID(cc *lib.Cusctx, WechatID string) *KuaiMaoUser {
	cont := bson.M{
		"wechatid": WechatID,
	}
	var user KuaiMaoUser
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).One(&user)
	if err != nil {
		return nil
	}
	if user.WechatID != WechatID {
		return nil
	}
	return &user
}

func FindSpecByID(cc *lib.Cusctx, id bson.ObjectId, fields ...string) *KuaiMaoUser {
	if len(fields) == 0 {
		return nil
	}
	selector := bson.M{}
	for idx := range fields {
		selector[fields[idx]] = 1
	}
	var res KuaiMaoUser
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).FindId(id).Select(selector).One(&res)
	if err != nil {
		return nil
	}
	return &res
}

func FindByIndex(cc *lib.Cusctx, index int64) *KuaiMaoUser {
	cont := bson.M{
		"index": index,
	}
	var user KuaiMaoUser
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).One(&user)
	if err != nil {
		return nil
	}
	return &user
}
func FindByNiCheng(cc *lib.Cusctx, nicheng string) *KuaiMaoUser {
	cont := bson.M{
		"nicheng": nicheng,
	}
	var user KuaiMaoUser
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).One(&user)
	if user.CreateTime == 0 {
		cc.Errf("stuser.FindByNiCheng,nicheng=%v,err=%v", nicheng, err)
		return nil
	}
	return &user
}

func FindByID(cc *lib.Cusctx, id bson.ObjectId) *KuaiMaoUser {
	var user KuaiMaoUser
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).FindId(id).One(&user)
	if user.CreateTime == 0 {
		cc.Errf("stuser.FindByID,err=%v", err)
		return nil
	}
	return &user
}

func FindByIDV2(cc *lib.Cusctx, id bson.ObjectId) (*KuaiMaoUser, error) {
	var user *KuaiMaoUser
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).FindId(id).One(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FindByIcode(cc *lib.Cusctx, Icode string) *KuaiMaoUser {
	cont := bson.M{
		"icode": Icode,
	}
	var user KuaiMaoUser
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).One(&user)
	if err != nil {
		return nil
	}
	if user.CreateTime == 0 {
		cc.Errf("stuser.FindByIcode,Icode=%v,err=%v", Icode, err)
		return nil
	}
	return &user
}

func HasThisIcode(cc *lib.Cusctx, icode string) (*KuaiMaoUser, bool) {
	n, _ := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(nil).Count()
	if n == 0 {
		return nil, true
	}
	cont := bson.M{
		"icode": icode,
	}
	selector := bson.M{
		"mysub":      1,
		"ficode":     1,
		"createtime": 1,
	}
	var user KuaiMaoUser
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).Select(selector).One(&user)
	if user.CreateTime == 0 {
		return nil, false
	}
	if user.MyAvailableSub <= 0 {
		return nil, false
	}
	return &user, true
}

func GetCount(cc *lib.Cusctx) int64 {
	n, err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(nil).Count()
	if err != nil {
		cc.Errf("stuser.GetCount(),err=%v", err)
		return 0
	}
	return int64(n)
}

func (user *KuaiMaoUser) GetCountry(cc *lib.Cusctx) string {
	return user.Country
}
func (user *KuaiMaoUser) SetCountry(cc *lib.Cusctx) {
	change := bson.M{
		"$set": bson.M{
			"country": user.Country,
		},
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).UpdateId(user.ID, change)
}
func (user *KuaiMaoUser) GetDevice(cc *lib.Cusctx) string {
	return user.Device
}
func (user *KuaiMaoUser) SetDevice(cc *lib.Cusctx) {
	change := bson.M{
		"$set": bson.M{
			"device": user.Device,
		},
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).UpdateId(user.ID, change)
}
func (user *KuaiMaoUser) AddSuanLi(cc *lib.Cusctx, AddValue int64) {
	cc.Logf("userindex=%v,cur_has_suanli=%v,addvalue=%v,final=%v", user.Index, user.Suanli, AddValue, user.Suanli+AddValue)

	cont := bson.M{
		"_id": user.ID,
	}
	Change := mgo.Change{
		Update: bson.M{
			"$inc": bson.M{
				"suanli": AddValue,
			},
		},
		Upsert:    true,
		ReturnNew: true,
	}
	var newuser KuaiMaoUser
	_, err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).Apply(
		Change,
		&newuser,
	)
	if err == nil {
		user.Suanli = newuser.Suanli
		cc.Logf("AddSuanLi_success:%v", user.Suanli)
	}
}

func (user *KuaiMaoUser) AddB(cc *lib.Cusctx, AddValue int64) {
	cc.Logf("userindex=%v,cur_has_b=%v,addvalue=%v,final=%v", user.Index, user.B, AddValue, user.B+AddValue)

	cont := bson.M{
		"_id": user.ID,
		"b":   user.B,
	}
	Change := mgo.Change{
		Update: bson.M{
			"$inc": bson.M{
				"b": AddValue,
			},
		},
		Upsert:    true,
		ReturnNew: true,
	}
	var newuser KuaiMaoUser
	_, err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).Apply(
		Change,
		&newuser,
	)
	if err == nil {
		user.B = newuser.B
		cc.Logf("AddB_success:%v", user.B)
	}
}
func (user *KuaiMaoUser) IncB(cc *lib.Cusctx, AddValue int64) {
	Change := bson.M{
		"$inc": bson.M{
			"b": AddValue,
		},
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).UpdateId(user.ID, Change)
}
func (user *KuaiMaoUser) IncSuanli(cc *lib.Cusctx, AddValue int64) {
	Change := bson.M{
		"$inc": bson.M{
			"suanli": AddValue,
		},
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).UpdateId(user.ID, Change)
}
func (user *KuaiMaoUser) SetEmail(cc *lib.Cusctx, email *Email) {
	user.Email.Addr = email.Addr
	user.Email.Status = email.Status
	updater := bson.M{
		"$set": bson.M{
			"email": user.Email,
		},
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).UpdateId(user.ID, updater)
}

func (user *KuaiMaoUser) SetFillCode(cc *lib.Cusctx, fillcode *FillCode) {
	user.FillCode.QQ = fillcode.QQ
	user.FillCode.Wechat = fillcode.Wechat
	user.FillCode.Kabao = fillcode.Kabao
	updater := bson.M{
		"$set": bson.M{
			"fillcode": user.FillCode,
		},
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).UpdateId(user.ID, updater)
}

func (user *KuaiMaoUser) SetFIcode(cc *lib.Cusctx, FIcode string) {
	updater := bson.M{
		"$set": bson.M{
			"ficode": FIcode,
		},
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).UpdateId(user.ID, updater)
}

func GetSuanLiLeaderBoard(cc *lib.Cusctx, count int) *mgo.Iter {
	return cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(nil).Sort("-suanli", "username").Limit(count).Iter()
}

func (user *KuaiMaoUser) GetSelfRanking(cc *lib.Cusctx) int {
	cond := bson.M{
		"suanli": bson.M{
			"$gte": user.Suanli,
		},
	}
	rank, _ := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cond).Count()
	return rank
}

func GetActivedUsers(cc *lib.Cusctx, limit int64) (int64, int64) {
	cont := []bson.M{
		bson.M{
			"$match": bson.M{
				"lla": bson.M{
					"$gt": limit,
				},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "redis_black_tel",
				"localField":   "telephone",
				"foreignField": "_id",
				"as":           "matched_docs",
			},
		},
		bson.M{
			"$match": bson.M{
				"matched_docs": bson.M{
					"$eq": []bson.M{},
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "",
				"amount": bson.M{
					"$sum": "$suanli",
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}

	type result struct {
		Amount int64 `bson:"amount"`
		Count  int64 `bson:"count"`
	}
	resp := []result{}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Pipe(cont).All(&resp)

	fmt.Printf("result=%v\n", resp)
	if len(resp) == 0 {
		return 0, 0
	}
	return resp[0].Amount, resp[0].Count
}

func (user *KuaiMaoUser) SetLastLoginData(cc *lib.Cusctx, lla int64) {
	updater := bson.M{
		"$set": bson.M{
			"lla": lla,
		},
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).UpdateId(user.ID, updater)
}

func ConsumeOneSub(cc *lib.Cusctx, icode string) {
	cont := bson.M{
		"icode": icode,
	}
	updater := bson.M{
		"$inc": bson.M{
			"mysub": -1,
		},
	}
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Update(cont, updater)
}
func ConcurConsumeOneSub(cc *lib.Cusctx, icode string) int {
	cont := bson.M{
		"icode": icode,
		"mysub": bson.M{
			"$gt": 0,
		},
	}
	Change := mgo.Change{
		Update: bson.M{
			"$inc": bson.M{
				"mysub": -1,
			},
		},
		Upsert:    false,
		ReturnNew: true,
	}
	var newuser KuaiMaoUser
	changeInfo, err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).Apply(
		Change,
		&newuser,
	)
	if err != nil {
		return 0
	}
	return changeInfo.Updated
}
func (user *KuaiMaoUser) ModifyNiCheng(cc *lib.Cusctx, nicheng string) {
	updater := bson.M{
		"$set": bson.M{
			"nicheng": nicheng,
		},
	}
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Update(bson.M{"_id": user.ID}, updater)
	if err != nil {
		fmt.Printf("modify nicheng err=%v\n", err)
	}
}

func (user *KuaiMaoUser) GetCountryNum() string {
	tele := strings.Split(user.Telephone, "_")
	if len(tele) < 2 {
		return "86"
	}
	return tele[0]
}

func (user *KuaiMaoUser) HasEnoughB(B int64) bool {
	return user.B >= B
}

func (user *KuaiMaoUser) CousumeB(cc *lib.Cusctx, B int64) {
	cc.Logf("cur_has_b=%v,CousumeB=-%v,final=%v", user.B, B, user.B-B)
	if B <= 0 {
		return
	}
	cont := bson.M{
		"_id": user.ID,
		"b":   user.B,
	}
	Change := mgo.Change{
		Update: bson.M{
			"$inc": bson.M{
				"b": -B,
			},
		},
		Upsert:    true,
		ReturnNew: true,
	}
	var newuser KuaiMaoUser
	_, err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).Apply(
		Change,
		&newuser,
	)
	if err == nil {
		user.B = newuser.B
		cc.Logf("CousumeB_success:%v", user.B)
	}
}

func GetAllUser(cc *lib.Cusctx) []*KuaiMaoUser {

	cont := bson.M{}
	var user []*KuaiMaoUser
	cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).All(&user)

	return user
}

func ChangeUserSub(cc *lib.Cusctx) {
	cont := bson.M{
		"index": bson.M{
			"$lte": 9999,
		},
	}
	Change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"mysub": 10,
			},
		},
		Upsert:    true,
		ReturnNew: true,
	}
	var newuser KuaiMaoUser
	_, err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).Apply(
		Change,
		&newuser,
	)
	cc.Logf("username=%v,subfrom%vto10", newuser.UserName, newuser.MyAvailableSub)
	if err == nil {
		cc.Logf("_changesubsuccess")
	} else {
		cc.Logf("_changesubfailed")
	}
}

func ExistWechatID(cc *lib.Cusctx, wechatid string) bool {
	cont := bson.M{
		"wechatid": wechatid,
	}
	n, err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).Count()
	if err != nil {
		cc.Errf("ExistWechatID_err=%v\n", err)
	}
	return n > 0
}

func (user *KuaiMaoUser) SetWechatID(cc *lib.Cusctx, wechatid string) bool {
	change := bson.M{
		"$set": bson.M{
			"wechatid": wechatid,
		},
	}
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).UpdateId(user.ID, change)
	if err != nil {
		return false
	}
	return true
}

func (user *KuaiMaoUser) UnsetWechatID(cc *lib.Cusctx) bool {
	change := bson.M{
		"$unset": bson.M{
			"wechatid": "",
		},
	}
	err := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).UpdateId(user.ID, change)
	if err != nil {
		return false
	}
	return true
}
func (user *KuaiMaoUser) CountMyInvitedNum(cc *lib.Cusctx) int {
	cont := bson.M{
		"ficode": bson.RegEx{
			Pattern: user.Icode,
			Options: "",
		},
	}
	n, _ := cc.M.DB(DBNAME_KuaiMaoUser).C(CNAME_KuaiMaoUser).Find(cont).Count()
	return n
}
