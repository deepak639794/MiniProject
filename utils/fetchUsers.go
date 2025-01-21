package utils

import (
	"Mymodule/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FetchUsers(db *gorm.DB, c *gin.Context) {
	fmt.Println("inside the fetch user")
	pageSizeStr := c.DefaultQuery("page_size", "100")
	pageStr := c.DefaultQuery("page_no", "")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 100
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // Default to 1 if invalid
	}

	offset := (page - 1) * pageSize
	fmt.Println("########")
	var userEntry []models.UserDetails
	if err := db.Order("id").Limit(pageSize).Offset(offset).Find(&userEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching users from the database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page_n0":   page,
		"page_size": pageSize,
		"users":     userEntry,
	})

}
