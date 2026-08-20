package database

import (
	"Server/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: "root:app_root_pwd@tcp(127.0.0.1:3306)/app_db?charset=utf8mb4&parseTime=True&loc=Local",
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	DB = db

	// AutoMigrate 创建/更新表结构（开发阶段用；生产建议用正式迁移工具）
	if err := db.AutoMigrate(&models.UserModel{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
