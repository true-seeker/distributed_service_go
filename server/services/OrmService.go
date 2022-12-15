package services

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgresConnectionString = fmt.Sprintf("host=%s "+
	"user=%s "+
	"password=%s "+
	"dbname=%s "+
	"port=%s "+
	"sslmode=disable TimeZone=Asia/Yekaterinburg",
	GetProperty("DataBase", "address"),
	GetProperty("DataBase", "user"),
	GetProperty("DataBase", "password"),
	GetProperty("DataBase", "dbname"),
	GetProperty("DataBase", "port"))

type Task struct {
	gorm.Model
	Answer           uint32
	BackpackCapacity uint32
	Items            []BackpackTaskItem `gorm:"foreignKey:TaskPartId"`
}

type BackpackTaskItem struct {
	ID         uint `gorm:"primarykey"`
	Weight     uint32
	Price      uint32
	TaskPartId uint
}

type TaskUserSolution struct {
	gorm.Model
	TaskId uint
	Task   Task `gorm:"foreignKey:TaskId;references:ID"`
	UserId uint
	User   User `gorm:"foreignKey:UserId;references:ID"`
	Answer uint32
}

type User struct {
	gorm.Model
	Username string
	Password string
}

func Migrate() {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")

	err = db.AutoMigrate(&Task{}, &BackpackTaskItem{}, &TaskUserSolution{}, &User{})
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

func GetUserByUsername(user User) User {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")
	existingUser := new(User)

	db.Where("username = ?", user.Username).First(existingUser)

	return *existingUser
}

func CheckIfUserAlreadyDidTheTask(user User, task Task) bool {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")
	user = GetUserByUsername(user)

	taskUserSolution := new(TaskUserSolution)

	err = db.Where("user_id = ? AND task_id = ?", user.ID, task.ID).First(taskUserSolution).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func GetUser(user User) User {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")
	existingUser := new(User)

	db.Where("username = ? AND password = ?", user.Username, user.Password).First(existingUser)

	return *existingUser
}

func CreateNewTaskUserSolution(solution TaskUserSolution) TaskUserSolution {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")

	db.Create(&solution)

	return solution
}
