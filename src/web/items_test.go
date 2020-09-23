package web

import (
	"net/http"
	"net/http/httptest"
	"projecto/app"
	"projecto/model"

	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type itemRepoMock struct {
	mock.Mock
}

type itemMock struct {
	mock.Mock
}

func (m *itemMock) Setup(*app.App) error {
	return nil
}

func (m *itemMock) Create(*model.Item) error {
	return nil
}

func (m *itemMock) List() ([]*model.Item, error) {
	items := []*model.Item{{ID: 1, Name: "Foo", Price: 123}}
	return items, nil
}

func (m *itemMock) Name() string {
	return "service/item"
}

func (m *itemRepoMock) Setup(*app.App) error {
	return nil
}

func (m *itemRepoMock) Name() string {
	return "itemrepo"
}

func TestItems(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)

	a := &app.App{}
	require.Nil(t, a.Start())

	repoMock := &itemRepoMock{}
	item := &itemMock{}
	a.Register(repoMock)
	a.Register(item)

	api := itemsAPI{
		item: item,
	}
	h := &api

	assert.NoError(t, h.items(c))
	require.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(
		t,
		`[{"id":1,"name":"Foo","price":123}]`,
		rec.Body.String(),
	)
}
