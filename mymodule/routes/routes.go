package routes

import (
	// "Mymodule/mymodule/middleware"

	"Mymodule/mymodule/middleware"
	"Mymodule/mymodule/utils"
	"fmt"

	// "Mymodule/mymodule/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes initializes all the routes for the application
func SetupRoutes(r *gin.Engine, logReader utils.LogReader, db *gorm.DB) {
	// Apply middleware
	r.Use(middleware.LogRequestResponse())

	// Define routes
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "hello sir"})
	})

	r.GET("/logs", func(c *gin.Context) {
		fmt.Println("#############inside the  logs function ")
		logs, err := logReader.ReadLogsFromFile()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch logs"})
			return
		}
		c.JSON(200, gin.H{"logs": logs})
	})

	r.POST("/uploadFile", utils.UploadFileDB)

	r.GET("/users", func(ctx *gin.Context) {
		utils.FetchUsers(db, ctx)
	})

}
