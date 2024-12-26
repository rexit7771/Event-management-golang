package structs

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	Date        string `json:"date" gorm:"type:date" validate:"required,date_after_today"`
	Location    string `json:"location" gorm:"not null"`
	Image_url   string `json:"image_url"`
	Approved    bool   `json:"approved" gorm:"default:0"`
	Created_by  uint   `json:"created_by" gorm:"not null"`
	User        User   `gorm:"foreignKey:Created_by;references:ID"`
}

func init() {
	validate := validator.New()
	validate.RegisterValidation("date_after_today", validateDateAfterToday)
}

func validateDateAfterToday(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	return dateStr == tomorrow
}
