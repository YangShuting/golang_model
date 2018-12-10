package handlers

import (
	"net/http"

	"github.com/Yangshuting/golang_model/lib"
	"github.com/Yangshuting/golang_model/model"
	"github.com/labstack/echo"
)

func FindUserNiCheng(c echo.Context) error {
	cc := c.Get("cc").(*lib.Cusctx)
	user := lib.GetUser(c).(*model.KuaiMaoUser)
	// fmt.Printf("user_%+v", user)
	findU := model.FindByTelephone(cc, user.Telephone)
	return c.JSON(http.StatusOK, findU.NiCheng)
}
