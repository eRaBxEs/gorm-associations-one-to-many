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
	ID       uint64 `gorm:"primaryKey"`
	Username string `gorm:"size:64"`
	Password string `gorm:"size:255"`
}

type Note struct {
	gorm.Model
	ID      uint64 `gorm:"primaryKey"`
	Name    string `gorm:"size:255"`
	Content string `gorm:"type:text"`
	UserID  uint64 `gorm:"index"`
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint64
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
	database, err := gorm.Open(mysql.Open("root:dbadev_sys@tcp(127.0.0.1:3306)/gorm_belongs_to?charset=utf8&parseTime=true"), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic("Failed to connect to databse!")
	}

	DB = database
}

func dbMigrate() {
	DB.AutoMigrate(&Note{}, &User{}, &CreditCard{})

	var note Note
	DB.First(&note)
	var user User
	// to get user to which the note belongs
	DB.Where("id = ?", note.UserID).First(&user)
	fmt.Printf("User from a note: %s\n", user.Username)

	fmt.Println("\n----------------")

	var notes []Note
	// to get all notes for a given user
	DB.Where("user_id = ?", user.ID).Find(&notes)

	fmt.Println("Notes from a user:")
	for _, element := range notes {
		fmt.Printf("%s - %s\n", element.Name, element.Content)
	}
	fmt.Println("\n----------------")

	var cc CreditCard
	// to get credit card details for given user
	DB.Where("user_id = ?", user.ID).First(&cc)
	fmt.Printf("Credit Card from a user: %s\n", cc.Number)
}

func main() {
	connectDatabase()
	dbMigrate()

}
