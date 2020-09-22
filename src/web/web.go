package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"projecto/app"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const ComponentName = "web"

func New() *API {
	return new(API)
}

type componentAPI interface {
	Setup(a *app.App, group *echo.Group) (err error)
}

type webConfig interface {
	GetAddr() string
}

type logger interface {
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type API struct {
	addr   string
	echo   *echo.Echo
	logger logger
}

func (api *API) Setup(a *app.App) (err error) {
	api.logger = a.Logger
	api.echo = echo.New()
	api.addr = a.MustComponent("config").(webConfig).GetAddr()
	return api.registerComponents(a)
}

func (api *API) registerComponents(a *app.App) (err error) {
	var components []interface{}

	components = append(
		components,
		new(infoAPI),
	)

	group := api.makeMainGroup()
	for _, comp := range components {
		switch c := comp.(type) {
		case componentAPI:
			if err = c.Setup(a, group); err != nil {
				return
			}
		default:
			return fmt.Errorf("unexpected api component: %T", comp)
		}
	}
	return
}

func (api *API) makeMainGroup() (group *echo.Group) {
	corsConfig := middleware.DefaultCORSConfig
	corsConfig.AllowHeaders = []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Authorization"}
	api.echo.Use(middleware.CORSWithConfig(corsConfig))
	return api.echo.Group(
		"v1/",
		middleware.Recover(),
	)
}

func (*API) Name() (name string) {
	return ComponentName
}

func (api *API) Run() (err error) {
	errCh := make(chan error)
	go func() {
		if err = api.echo.Start(api.addr); err != nil {
			select {
			case errCh <- err:
				return
			default:
				if err != http.ErrServerClosed {
					log.Fatalf("can't start server: %v", err)
				}
			}
		}
	}()
	select {
	case err = <-errCh:
		return
	case <-time.After(time.Millisecond * 100):
		return
	}
}

func (api *API) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return api.echo.Shutdown(ctx)
}
