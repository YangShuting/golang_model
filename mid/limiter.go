package mid

import (
	"fmt"

	"github.com/Yangshuting/golang_model/lib"
	"github.com/labstack/echo"
)

func ReqLimitMid(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		RequestLimit := c.Get("reqLim").(*lib.ReqLimiterService)
		key := RequestLimit.GetIPAndUri(c)
		if RequestLimit.IsAvaliable(key) {
			RequestLimit.Increase(key)
		} else {
			fmt.Printf("Reach request limiting......!")
		}
		return next(c)
	}
}
