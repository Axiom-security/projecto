package itemrepo

//go:generate sh -c "go run $GOPATH/$GENERATOR -c $GOPATH/$TEMPLATE -n Item -p itemrepo -m '*model.Item' > $GOPATH/src/service/item/itemrepo/repo.gen.go"

import (
	"projecto/app"
	"projecto/db"
	"projecto/model"
)

const ComponentName = "repository/item"

type Repository interface {
	iItemRepo
	app.Component
}

func New() Repository {
	return new(impl)
}

type impl struct {
	*repoItem
	db *db.Database
}

func (r *impl) Setup(a *app.App) error {
	r.db = a.MustComponent(db.ComponentName).(*db.Database)
	r.repoItem = newRepoItem(r.db)
	if err := r.db.GetDatabase().AutoMigrate(&model.Item{}); err != nil {
		return err
	}
	return nil
}

func (r *impl) Name() string {
	return ComponentName
}
