package models

import (
	"fmt"
	"log"
	"logistics-tracker/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&User{}, &Waybill{}, &Tracking{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedDefaultUser()
	log.Println("Database initialized successfully")
}

func seedDefaultUser() {
	var count int64
	DB.Model(&User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		user := &User{Username: "admin"}
		err := user.SetPassword("admin123")
		if err != nil {
			log.Printf("Failed to set default user password: %v", err)
			return
		}
		DB.Create(user)
		log.Println("Default user created: admin / admin123")
	}
}
