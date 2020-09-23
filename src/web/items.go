package web

import (
	"encoding/json"
	"net/http"
	"projecto/app"
	"projecto/model"
	"projecto/service/item"

	"github.com/labstack/echo"
)

type itemsAPI struct {
	a    *app.App
	item item.Item
}

func (api *itemsAPI) Setup(a *app.App, g *echo.Group) (err error) {
	api.a = a
	api.item = a.MustComponent("service/item").(item.Item)

	items := g.Group("items")
	items.GET("", api.items)
	items.POST("", api.create)
	return
}

func (api *itemsAPI) items(ctx echo.Context) error {
	items, err := api.item.List()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err})
	}
	return ctx.JSON(http.StatusOK, items)
}

func (api *itemsAPI) create(ctx echo.Context) error {
	var item model.Item
	if err := json.NewDecoder(ctx.Request().Body).Decode(&item); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err})
	}
	if err := item.Validate(); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err})
	}
	if err := api.item.Create(&item); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err})
	}
	return ctx.JSON(http.StatusOK, item)
}
