package web

import (
	"net/http"

	"projecto/app"

	"github.com/labstack/echo"
)

type infoAPI struct {
	a *app.App
}

func (api *infoAPI) Setup(a *app.App, g *echo.Group) (err error) {
	api.a = a

	g.Group("info").GET("", api.info)
	return
}

func (api *infoAPI) info(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{"name": api.a.Name(), "version": api.a.Version()})
}
