package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	ID         uint64 `gorm:"primaryKey"`
	Username   string `gorm:"size:64"`
	Password   string `gorm:"size:255"`
	Notes      []Note
	CreditCard *CreditCard
}

type Note struct {
	gorm.Model
	ID      uint64 `gorm:"primaryKey"`
	Name    string `gorm:"size:255"`
	Content string `gorm:"type:text"`
	UserID  uint64 `gorm:"index"`
	User    User
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint64
	User   User // we can also get the user if we have the credit card by adding user to the credit card struct
}

var DB *gorm.DB

func connectDatabase() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	dsn := "henry_dev:devdba_user@tcp(127.0.0.1:3307)/gorm_testdb?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic("Failed to connect to databse!")
	}

	DB = database
}

func dbMigrate() {
	DB.AutoMigrate(&Note{}, &User{}, &CreditCard{})

	var note Note
	// Using single Preload function from gorm for the user
	DB.Preload("User").First(&note)
	// to get user to which the note belongs
	fmt.Printf("User from a note: %s\n", note.User.Username)

	fmt.Println("\n----------------")

	// Using chained double Preload function from gorm for the notes and credit card details on the user
	var user User
	DB.Preload("Notes").Preload("CreditCard").Where("username = ?", "erabxes").Find(&user)

	fmt.Println("Notes from a user:")
	for _, element := range user.Notes {
		fmt.Printf("%s - %s\n", element.Name, element.Content)
	}
	fmt.Println("\n----------------")

	// to get credit card details for given user
	fmt.Printf("Credit Card from a user: %s\n", user.CreditCard.Number)
}

func main() {
	connectDatabase()
	dbMigrate()

}
