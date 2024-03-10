package migrator

import (
	"go-image-api/internal/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.History{}); err != nil {
		return err
	}

	return nil
}

func Drop(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&entity.History{}); err != nil {
		return err
	}

	return nil
}
