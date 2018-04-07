package google

import (
	"github.com/labstack/echo"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

func GetContext(c echo.Context) context.Context {
	return appengine.NewContext(c.Request())
}
