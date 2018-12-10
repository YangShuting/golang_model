package main

import (
	"os"

	"github.com/Yangshuting/golang_model/config"
	"github.com/Yangshuting/golang_model/mid"
	"github.com/Yangshuting/golang_model/model"
	"github.com/Yangshuting/golang_model/router"
	"github.com/Yangshuting/golang_model/storage"
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
