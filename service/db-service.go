package service

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/lowspoons-server/model"

	//postgres implicits
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PGService struct {
	DB *gorm.DB
}

func (i *PGService) Connect() (*gorm.DB, error) {
	var err error
	i.DB, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=lowspoons sslmode=disable")

	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
		return nil, err
	}

	i.DB.LogMode(true)

	i.DB.AutoMigrate(&model.User{}, &model.Todo{}, &model.Group{})

	return i.DB, nil
}
