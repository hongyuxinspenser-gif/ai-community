package config

import (
	"log"
	"ai-community/models" // 引用 models 包，用于自动建表

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	// 连接 SQLite 数据库
	DB, err = gorm.Open(sqlite.Open("community.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移模式：根据 Model 结构体创建/更新表结构
	// 注意：这里需要传入具体的 Model 实例
	err = DB.AutoMigrate(&models.Post{}, &models.Comment{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully.")
}