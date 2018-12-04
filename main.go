package main

import (
	"os"

	"gitee.com/firewing_group/blue_kxq2/config"
	"gitee.com/firewing_group/blue_kxq2/mid"
	"gitee.com/firewing_group/blue_kxq2/model"
	"gitee.com/firewing_group/blue_kxq2/router"
	"gitee.com/firewing_group/blue_kxq2/storage"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/unrolled/secure"
)

//echo log middleware
func main() {

	{
		config.Init()
		storage.InitMongo()
		// storage.InitRedisConfig()
		model.InitModel()
	}

	e := echo.New()
	e.Use(mid.SuperCtx)
	//secure： https://github.com/unrolled/secure
	secureMiddleware := secure.New()
	e.Use(echo.WrapMiddleware(secureMiddleware.Handler))
	e.Use(middleware.CORS())
	// e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{}))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static",
		Browse: true,
	}))
	// logger
	e.Use(middleware.Logger())
	router.Start(e)
	e.Start(os.Getenv("bind"))
}
