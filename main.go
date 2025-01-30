package main

import (
	"Mymodule/mymodule/middleware"
	"Mymodule/mymodule/models"
	"Mymodule/mymodule/routes"
	"Mymodule/mymodule/utils"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	logReader utils.LogReader
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

	r.Use(middleware.LogRequestResponse()) //Middleware in Gin is a function that runs before and after
	logReader = utils.NewLogReader()
	routes.SetupRoutes(r, logReader, db)

	r.Run(":8081")

}
