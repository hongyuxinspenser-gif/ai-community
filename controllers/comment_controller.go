package controllers

import (
	"ai-community/config"
	"ai-community/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&comment)
	c.JSON(http.StatusOK, comment)
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Comment{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}

// GetCommentsByPost 获取树状评论
func GetCommentsByPost(c *gin.Context) {
	postID := c.Param("post_id")
	
	var allComments []models.Comment
	// 查出该帖子下所有评论
	if err := config.DB.Where("post_id = ?", postID).Find(&allComments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}

	// 1. 构建树
	tree := buildTree(allComments)

	// 2. 分页逻辑 (针对顶级评论)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	total := len(tree)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > total { start = total }
	if end > total { end = total }

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"page":  page,
		"list":  tree[start:end],
	})
}

// buildTree 辅助函数：将扁平切片转为树状切片
// 注意：这个函数只在包内部使用，所以小写开头
func buildTree(comments []models.Comment) []*models.Comment {
	commentMap := make(map[uint]*models.Comment)
	for i := range comments {
		commentMap[comments[i].ID] = &comments[i]
	}

	var roots []*models.Comment
	for i := range comments {
		curr := &comments[i]
		if curr.ParentID != nil {
			if parent, ok := commentMap[*curr.ParentID]; ok {
				parent.Children = append(parent.Children, curr)
			}
		} else {
			roots = append(roots, curr)
		}
	}
	return roots
}