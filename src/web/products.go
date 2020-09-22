package web

import (
	"projecto/app"

	"github.com/labstack/echo"
)

type productsAPI struct {
	a *app.App
}

func (api *productsAPI) Init(a *app.App, g *echo.Group) (err error) {
	api.a = a

	g.Group("products").GET("/", api.products)
	return
}

func (api *productsAPI) products(ctx echo.Context) error {
	return nil
}
