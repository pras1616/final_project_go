package database

import (
	"final_project/models"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Pras1616"
	dbname   = "postgres"
)

var (
	dbComment     *gorm.DB
	dbUser        *gorm.DB
	dbPhoto       *gorm.DB
	dbSocialMedia *gorm.DB
	err           error
)

func ConnectDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	dbComment, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbUser, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbPhoto, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbSocialMedia, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := dbUser.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetConnMaxIdleTime(time.Second * 2)
	sqlDB.SetConnMaxLifetime(time.Second * 2)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	fmt.Println("Success connect to DB using GORM")

	dbComment.AutoMigrate(&models.Comment{})
	dbUser.AutoMigrate(&models.User{})
	dbPhoto.AutoMigrate(&models.Photo{})
	dbSocialMedia.AutoMigrate(&models.SocialMedia{})
}

func GetDB_Comment() *gorm.DB {
	return dbComment
}

func GetDB_User() *gorm.DB {
	return dbUser
}

func GetDB_Photo() *gorm.DB {
	return dbPhoto
}

func GetDB_SocialMedia() *gorm.DB {
	return dbSocialMedia
}
