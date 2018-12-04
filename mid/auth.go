package mid

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"gitee.com/firewing_group/blue_kxq2/lib"
	"gitee.com/firewing_group/blue_kxq2/model"
	"gitee.com/firewing_group/blue_kxq2/storage"
	"github.com/labstack/echo"
)

func AuthMid(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session := c.QueryParam("session")
		user, err := auth(c, session)
		if err != nil {
			return c.JSON(http.StatusBadRequest, lib.WXError(err.Error(), lib.STATUS_BAD_REQUEST))
		}
		lib.SetUser(c, user)
		return next(c)
	}
}

func auth(c echo.Context, session string) (*model.KuaiMaoUser, error) {
	cc := c.Get("cc").(*lib.Cusctx)
	uid, err := storage.GetRedis(session)
	if err != nil {
		return nil, err
	}
	user := model.FindByID(cc, bson.ObjectIdHex(uid))
	return user, nil
}
