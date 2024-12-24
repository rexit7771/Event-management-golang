package structs

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate *validator.Validate

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password string `json:"password" gorm:"not null" validate:"required,min=5,max=10"`
	Role     string `json:"role" gorm:"default:member"`
}

func init() {
	validate = validator.New()
}

func (u *User) Validate() error {
	return validate.Struct(u)
}
