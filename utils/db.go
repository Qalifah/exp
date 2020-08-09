package utils

import (
	"exp/models"
	"fmt"
	"log"
	"os"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres" //gorm's postgre dialect interface
	"github.com/joho/godotenv"
)

func ConnectDB() *gorm.DB {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("databaseUser")
	password := os.Getenv("databasePassword")
	databaseName := os.Getenv("databaseName")


	dbURI := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", username, databaseName, password)
	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Println("error", err)
		panic(err)
	}
	// defer db.Close()
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	fmt.Println("Connected Successfully!", db)
	return db
}