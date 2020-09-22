package app

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	version string
	name    string
)

type Component interface {
	Setup(a *App) (err error)
	Name() (name string)
}

type logger interface {
	Warnf(format string, args ...[]interface{})
	Infof(format string, args ...[]interface{})
	Debugf(format string, args ...[]interface{})
}

type ComponentRunnable interface {
	Component
	Run() error
	Close() error
}

type App struct {
	components []Component
	mu         sync.RWMutex
	logger     logger
}

func (app *App) Name() string {
	return name
}

func (app *App) Version() string {
	return version
}

func (app *App) Register(s Component) *App {
	app.mu.Lock()
	defer app.mu.Unlock()
	for _, es := range app.components {
		if s.Name() == es.Name() {
			panic(fmt.Errorf("component '%s' already registered", s.Name()))
		}
	}
	app.components = append(app.components, s)
	return app
}

func (app *App) Component(name string) Component {
	app.mu.RLock()
	defer app.mu.RUnlock()
	for _, s := range app.components {
		if s.Name() == name {
			return s
		}
	}
	return nil
}

func (app *App) MustComponent(name string) Component {
	s := app.Component(name)
	if s == nil {
		panic(fmt.Errorf("component '%s' not registered", name))
	}
	return s
}

func (app *App) ComponentNames() (names []string) {
	app.mu.RLock()
	defer app.mu.RUnlock()
	names = make([]string, len(app.components))
	for i, c := range app.components {
		names[i] = c.Name()
	}
	return
}

func (app *App) Start() (err error) {
	app.mu.RLock()
	defer app.mu.RUnlock()

	closeServices := func(idx int) {
		for i := idx; i >= 0; i-- {
			if serviceClose, ok := app.components[i].(ComponentRunnable); ok {
				if e := serviceClose.Close(); e != nil {
					app.log.Warnf("Component '%s' close error: %v", serviceClose.Name(), e)
				}
			}
		}
	}

	for i, s := range app.components {
		if err = s.Init(app); err != nil {
			closeServices(i)
			return fmt.Errorf("can't init service '%s': %v", s.Name(), err)
		}
	}

	for i, s := range app.components {
		if serviceRun, ok := s.(ComponentRunnable); ok {
			if err = serviceRun.Run(); err != nil {
				closeServices(i)
				return fmt.Errorf("can't run service '%s': %v", serviceRun.Name(), err)
			}
		}
	}

	return
}

func (app *App) Close() error {
	app.log.Infof("Close components...")
	app.mu.RLock()
	defer app.mu.RUnlock()
	done := make(chan struct{})
	go func() {
		select {
		case <-done:
			return
		case <-time.After(time.Minute):
			panic("app.Close timeout")
		}
	}()

	var errs []string
	for i := len(app.components) - 1; i >= 0; i-- {
		if serviceClose, ok := app.components[i].(ComponentRunnable); ok {
			app.log.Debugf("Close '%s'", serviceClose.Name())
			if e := serviceClose.Close(); e != nil {
				errs = append(errs, fmt.Sprintf("Component '%s' close error: %v", serviceClose.Name(), e))
			}
		}
	}
	close(done)
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
