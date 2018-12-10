package mid

import (
	"net/http"

	"github.com/Yangshuting/golang_model/lib"
	"github.com/Yangshuting/golang_model/model"
	"github.com/Yangshuting/golang_model/storage"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
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
