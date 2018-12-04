package handlers

import (
	"fmt"
	"net/http"

	"gitee.com/firewing_group/blue_kxq2/lib"
	"gitee.com/firewing_group/blue_kxq2/model"
	"github.com/labstack/echo"
)

func Login(c echo.Context) error {
	cc := c.Get("cc").(*lib.Cusctx)
	loginParam := new(model.LoginParams)
	c.Request().ParseForm()
	if err := c.Bind(loginParam); err != nil {
		cc.Errf("body_parse_err_%v", err.Error())
		return c.JSON(http.StatusBadRequest, lib.WXError("body_parse_err", lib.STATUS_BAD_REQUEST))
	}
	var user *model.KuaiMaoUser
	user = model.FindByTelephone(cc, loginParam.Telephone)
	if user == nil {
		// new user
		newuser := model.NewKuaiMaoUser()
		newuser.UserName = 0
		newuser.NiCheng = "user_" + loginParam.Telephone
		newuser.Telephone = loginParam.Telephone
		newuser.Index = 0
		newuser.Icode = loginParam.Telephone
		newuser.FIcode = loginParam.Telephone
		newuser.Suanli = 0
		newuser.ChaosticID = lib.Sha1(lib.Sha1(fmt.Sprintf("%v", loginParam.Telephone)))
		if channel, ok := cc.C.Get("channel").(string); ok {
			newuser.Channel = channel
		}
		err := newuser.Insert(cc)
		if err != nil {
			return c.JSON(http.StatusBadRequest, lib.WXError("new user error", lib.STATUS_BAD_REQUEST))
		}
		user = newuser
	}
	session, err := lib.SetSession(cc, user.ID.Hex())
	if err != nil {
		return c.JSON(http.StatusBadRequest, lib.WXError("set session error", lib.STATUS_BAD_REQUEST))
	}
	var loginReturn model.LoginReturn
	loginReturn.User = user
	loginReturn.Session = session
	return c.JSON(http.StatusOK, loginReturn)
}
