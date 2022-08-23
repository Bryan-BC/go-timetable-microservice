package db

import (
	"log"

	"github.com/Bryan-BC/go-timetable-microservice/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	DataBase *gorm.DB
}

func Init(url string) DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Panicf("Error opening database: %s \n", err)
		return DB{}
	}

	db.AutoMigrate(&models.Timetable{})
	return DB{DataBase: db}
}
