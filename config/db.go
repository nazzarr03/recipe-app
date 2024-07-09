package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/nazzarr03/recipe-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Db  *gorm.DB
	Rdb *redis.Client
)

func init() {
	ConnectDB()
	ConnectRedis()
}

func ConnectDB() {
	var err error
	dsn := "host=localhost user=postgres password=password dbname=recipe port=5432 sslmode=disable"
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	if err := Db.AutoMigrate(&models.User{}, &models.Recipe{}); err != nil {
		panic("failed to migrate database")
	}

	fmt.Println("Database connected successfully!")

}

func ConnectRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		panic("failed to connect to Redis")
	}

	fmt.Println("Redis connected successfully:", pong)
}
