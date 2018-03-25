// +build !appengine

package app

import (
	"github.com/labstack/echo"
)

func createMux() *echo.Echo {
	e := echo.New()
	return e
}
