package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Item struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (i Item) Validate() error {
	return validation.ValidateStruct(&i,
		// Name cannot be empty, and the length must between 5 and 50
		validation.Field(&i.Name, validation.Required, validation.Length(3, 50)),
		// Price cannot be empty, and must be between 1 and 500
		validation.Field(&i.Price, validation.Required, validation.Min(1), validation.Max(500)),
	)
}
