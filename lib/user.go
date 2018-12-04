package lib

import (
	"github.com/labstack/echo"
)

const USER = "USER"

func GetUser(ctx echo.Context) interface{} {
	if user := ctx.Get(USER); user != nil {
		return user
	}
	return nil
}

func SetUser(ctx echo.Context, user interface{}) {
	ctx.Set(USER, user)
}
