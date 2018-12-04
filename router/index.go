package router

import (
	"gitee.com/firewing_group/blue_kxq2/handlers"
	"gitee.com/firewing_group/blue_kxq2/mid"
	"github.com/labstack/echo"
)

type ModelSetting struct {
	NeedAuth bool
}

func Init() {}

func Start(e *echo.Echo) {
	e.POST("/blue_kxq2/hello", handlers.Hello, mid.AuthMid)
	e.POST("/blue_kxq2/getUserNiCheng", handlers.FindUserNiCheng, mid.AuthMid)
	e.POST("/blue_kxq2/login", handlers.Login)
	e.POST("blue_kxq2/newCommand", handlers.NewCommands, mid.AuthMid)
	e.POST("blue_kxq2/getCommandList", handlers.GetCommandLists, mid.AuthMid)
	e.POST("/blue_kxq/testSpeedLimiter", handlers.TestSpeedLimiter)
}
