package item

import (
	"projecto/app"
	"projecto/model"
	"projecto/service/item/itemrepo"
)

const ComponentName = "service/item"

type Service struct {
	db itemrepo.Repository
}

type Item interface {
	List() ([]*model.Item, error)
	Create(*model.Item) error
}

func New() *Service {
	return &Service{}
}

func (s *Service) Setup(a *app.App) (err error) {
	s.db = a.MustComponent("repository/item").(itemrepo.Repository)
	return
}

func (s *Service) List() (items []*model.Item, err error) {
	return s.db.Find()
}

func (s *Service) Create(item *model.Item) (err error) {
	return s.db.Create(item)
}

func (s Service) Name() string {
	return ComponentName
}

func (*Service) Close() error {
	return nil
}
