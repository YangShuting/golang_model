package mid

import (
	"time"

	"github.com/Yangshuting/golang_model/lib"
	"github.com/Yangshuting/golang_model/storage"
	"github.com/labstack/echo"
)

func SuperCtx(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s := storage.GetSession()
		defer s.Close()
		cc := lib.NewContext(c, nil, nil)
		reqLimiter := lib.NewReqLimiterService(10*time.Second, 5)
		cc.M = s
		c.Set("cc", cc)
		c.Set("reqLim", reqLimiter)
		return next(c)
	}
}
