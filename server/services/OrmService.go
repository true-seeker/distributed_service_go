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
	Answer           uint32
	TaskParts        []TaskPart `gorm:"foreignKey:TaskId"`
	BackpackCapacity uint32
}

type TaskPart struct {
	gorm.Model
	Answer uint32
	Items  []BackpackTaskItem `gorm:"foreignKey:TaskPartId"`
	TaskId uint
}

type BackpackTaskItem struct {
	gorm.Model
	Weight     uint32
	Price      uint32
	IsFixed    bool
	TaskPartId uint
}

type TaskUserSolution struct {
	gorm.Model
	TaskPartId uint
	TaskPart   TaskPart `gorm:"foreignKey:TaskPartId;references:ID"`
	UserId     uint
	User       User `gorm:"foreignKey:UserId;references:ID"`
	Answer     uint32
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

func GetUserByUsername(user User) User {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")
	existingUser := new(User)

	db.Where("username = ?", user.Username).First(existingUser)

	return *existingUser
}

func CheckIfUserAlreadyDidTheTask(user User, part TaskPart) bool {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")
	user = GetUserByUsername(user)

	taskUserSolution := new(TaskUserSolution)

	err = db.Where("user_id = ? AND task_part_id = ?", user.ID, part.ID).First(taskUserSolution).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func GetUser(user User) User {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")
	existingUser := new(User)

	db.Where("username = ? AND password = ?", user.Username, user.Password).First(existingUser)

	return *existingUser
}

func GetTaskById(taskId uint) Task {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{})
	FailOnError(err, "Failed to connect to DB")
	task := new(Task)

	db.Find(&task, taskId)
	return *task
}
