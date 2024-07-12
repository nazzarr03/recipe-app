package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/nazzarr03/recipe-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found")
	}
	ConnectDB()
	ConnectRedis()
	ConnectRabbitMQ()
}

func ConnectDB() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	if err := Db.AutoMigrate(&models.User{}, &models.Recipe{}); err != nil {
		panic("failed to migrate database")
	}

	fmt.Println("Database connected successfully!")

}
