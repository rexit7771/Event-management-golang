package helpers

import "gorm.io/gorm"

func PreloadUser(db *gorm.DB) *gorm.DB {
	return db.Select("id", "name", "email")
}
