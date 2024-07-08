package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	ConnectDB()
}

func ConnectDB() {
	var err error
	dsn := "host=localhost user=postgres password=password dbname=recipe port=5432 sslmode=disable"
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	if err := Db.AutoMigrate(); err != nil {
		panic("failed to migrate database")
	}

	fmt.Println("Database connected successfully!")

}
