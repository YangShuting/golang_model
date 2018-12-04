package handlers

import (
	"net/http"

	"gitee.com/firewing_group/blue_kxq2/lib"
	"gitee.com/firewing_group/blue_kxq2/model"
	"github.com/labstack/echo"
)

func FindUserNiCheng(c echo.Context) error {
	cc := c.Get("cc").(*lib.Cusctx)
	user := lib.GetUser(c).(*model.KuaiMaoUser)
	// fmt.Printf("user_%+v", user)
	findU := model.FindByTelephone(cc, user.Telephone)
	return c.JSON(http.StatusOK, findU.NiCheng)
}
