package services

import (
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
	//Task   Task `gorm:"foreignKey:TaskId;references:ID"`
}

type BackpackTaskItem struct {
	gorm.Model
	Weight     float64
	Price      float64
	IsFixed    bool
	TaskPartId uint
}

type TaskClientSolution struct {
	gorm.Model
	TaskPartId uint
	TaskPart   TaskPart `gorm:"foreignKey:TaskPartId;references:ID"`
	ClientId   uint
	Client     Client `gorm:"foreignKey:ClientId;references:ID"`
}

type Client struct {
	gorm.Model
	username string
}

func Migrate() {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")

	err = db.AutoMigrate(&Task{}, &TaskPart{}, &BackpackTaskItem{}, &TaskClientSolution{}, &Client{})
	FailOnError(err, "Failed to migrate")
}

func SaveNewTaskParts(task Task) Task {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")

	db.Create(&task)

	return task
}
