package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Comments []Comment `json:"comments,omitempty"` // 关联查询用
}