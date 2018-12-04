package mid

import (
	"gitee.com/firewing_group/blue_kxq2/lib"
	"gitee.com/firewing_group/blue_kxq2/storage"
	"github.com/labstack/echo"
)

func SuperCtx(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s := storage.GetSession()
		defer s.Close()
		cc := lib.NewContext(c, nil, nil)
		cc.M = s
		c.Set("cc", cc)
		return next(c)
	}
}
