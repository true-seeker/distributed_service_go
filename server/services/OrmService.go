package services

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgresConnectionString = fmt.Sprintf("host=localhost "+
	"user=%s "+
	"password=%s "+
	"dbname=%s "+
	"port=%s "+
	"sslmode=disable TimeZone=Asia/Yekaterinburg",
	GetProperty("DataBase", "user"),
	GetProperty("DataBase", "password"),
	GetProperty("DataBase", "dbname"),
	GetProperty("DataBase", "port"))

type Task struct {
	gorm.Model
	Answer    float64
	TaskParts []TaskPart `gorm:"foreignKey:TaskId"`
}

type TaskPart struct {
	gorm.Model
	Answer float64
	Items  []BackpackTaskItem `gorm:"foreignKey:TaskPartId"`
	TaskId uint
}

type BackpackTaskItem struct {
	gorm.Model
	Weight     float64
	Price      float64
	IsFixed    bool
	TaskPartId uint
}

type TaskUserSolution struct {
	gorm.Model
	TaskPartId uint
	TaskPart   TaskPart `gorm:"foreignKey:TaskPartId;references:ID"`
	UserId     uint
	User       User `gorm:"foreignKey:UserId;references:ID"`
}

type User struct {
	gorm.Model
	Username string
	Password string
}

func Migrate() {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")

	err = db.AutoMigrate(&Task{}, &TaskPart{}, &BackpackTaskItem{}, &TaskUserSolution{}, &User{})
	FailOnError(err, "Failed to migrate")
}

func SaveNewTaskParts(task Task) Task {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")

	db.Create(&task)

	return task
}

func RegisterNewUser(user User) error {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")

	existingUser := new(User)

	db.Where("username = ?", user.Username).First(existingUser)
	if existingUser.Username == "" {
		db.Create(&user)
	} else {
		return errors.New("user already exists")
	}

	return nil
}
