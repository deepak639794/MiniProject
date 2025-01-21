package main

import (
	"Mymodule/middleware"
	"Mymodule/models"
	"Mymodule/utils"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func main() {
	dbConStr := "host=postgres user=DeepakAgrawal password=03July2003@@ dbname=CSVDB port=5432 sslmode=disable TimeZone=UTC"
	var err error
	db, err = gorm.Open(postgres.Open(dbConStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}

	if err := db.AutoMigrate(&models.UserDetails{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	r := gin.Default()
	r.MaxMultipartMemory = 30 << 30 // 30 GB

	r.Use(middleware.LogRequestResponse())

	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "hello sir"})
	})
	r.GET("/logs", func(c *gin.Context) {
		// Read and return logs from the file
		logs, err := utils.ReadLogsFromFile()
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
	r.Run(":8081")
}
