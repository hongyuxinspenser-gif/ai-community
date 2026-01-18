package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID   uint   `json:"post_id"`
	ParentID *uint  `json:"parent_id"` // 指针类型，允许为 null
	Content  string `json:"content"`
	
	// 虚拟字段，用于 JSON 输出
	Children []*Comment `json:"children,omitempty" gorm:"-"`
}