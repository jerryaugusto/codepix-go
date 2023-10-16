package model

import (
	"time"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Base struct {
	ID        string    `json:"id,omitempty" valid:"uuid"`
	CreatedAt time.Time `json:"created_at,omitempty" valid:"-"`
	UpdatedAt time.Time `json:"updated_at,omitempty" valid:"-"`
}
