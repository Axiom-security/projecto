package db

import (
	"projecto/app"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const ComponentName = "db"

type config interface {
	GetDB() string
}

type Database struct {
	Database *gorm.DB
}

func New() *Database {
	return &Database{}
}

func (db *Database) Setup(a *app.App) (err error) {
	conn := a.MustComponent("config").(config).GetDB()
	if db.Database, err = gorm.Open(postgres.Open(conn), &gorm.Config{}); err != nil {
		return err
	}
	return
}

func (db Database) GetDatabase() *gorm.DB {
	return db.Database
}

func (db Database) Name() string {
	return ComponentName
}

func (db *Database) Close() error {
	return nil
}
