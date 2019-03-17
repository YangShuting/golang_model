package handlers

import (
	"fmt"
	"net/http"

	"github.com/Yangshuting/golang_model/lib"
	"github.com/Yangshuting/golang_model/model"
	"github.com/labstack/echo"
)

func NewCommands(c echo.Context) error {
	cc := c.Get("cc").(*lib.Cusctx)
	user := lib.GetUser(c).(*model.KuaiMaoUser)
	command := c.QueryParam("command")

	newCom, err := model.NewCommand(cc, user.ID.Hex(), command)
	if err != nil {
		return c.JSON(http.StatusBadRequest, lib.WXError(err.Error(), lib.STATUS_BAD_REQUEST))
	}

	model.RedisListPush("command_"+user.ID.Hex(), command)
	return c.JSON(http.StatusOK, newCom)
}

func GetCommandLists(c echo.Context) error {
	cc := c.Get("cc").(*lib.Cusctx)
	user := lib.GetUser(c).(*model.KuaiMaoUser)
	redisCom, databaseCom := model.GetRedisList(cc, "command_"+user.ID.Hex(), 0, 10)
	if redisCom != nil {
		cc.Logf("redis_command_%+v", redisCom)
		return c.JSON(http.StatusOK, redisCom.Val())
	}
	return c.JSON(http.StatusOK, databaseCom)
}

func TestSpeedLimiter(c echo.Context) error {
	bool, err := model.SpeedLimiter("test_speed_limiter_hello_world", 60)
	fmt.Printf("ifSuccess_%+v", bool)
	fmt.Printf("err_%+v", err)
	return nil

}
