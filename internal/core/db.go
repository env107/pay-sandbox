package core

import (
	"log"
	"wepay-sandbox/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) {
	var err error
	if dsn == "" {
		dsn = "sandbox.db"
	}
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移
	err = DB.AutoMigrate(
		&model.Merchant{},
		&model.Transaction{},
		&model.CallbackLog{},
		&model.Refund{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
