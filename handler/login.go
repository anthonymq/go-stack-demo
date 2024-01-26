package handler

import (
	"github.com/anthonymq/go-stack-demo/view/login"
	"github.com/labstack/echo/v4"
)

type LoginHandler struct{}

func (h LoginHandler) HandleLoginShow(c echo.Context) error {
	return render(c, login.Show())
}
