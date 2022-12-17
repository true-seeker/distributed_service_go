package services

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

type OrmConnection struct {
	conn gorm.DB
}

func GetDBConnection() OrmConnection {
	db, err := gorm.Open(postgres.Open(PostgresConnectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	FailOnError(err, "Failed to connect to DB")
	return OrmConnection{conn: *db}
}

func Migrate() {
	db := GetDBConnection()
	dbInstance, _ := db.conn.DB()
	defer dbInstance.Close()

	err := db.conn.AutoMigrate(&Task{}, &BackpackTaskItem{}, &TaskUserSolution{}, &User{})
	FailOnError(err, "Failed to migrate")
}

func (db *OrmConnection) SaveNewTaskParts(task Task) Task {
	db.conn.Create(&task)

	return task
}

func (db *OrmConnection) RegisterNewUser(user User) error {
	existingUser := new(User)

	db.conn.Where("username = ?", user.Username).First(existingUser)
	if existingUser.Username == "" {
		db.conn.Create(&user)
	} else {
		return errors.New("user already exists")
	}

	return nil
}

func (db *OrmConnection) GetUserByUsername(user User) User {
	existingUser := new(User)

	db.conn.Where("username = ?", user.Username).First(existingUser)

	return *existingUser
}

func (db *OrmConnection) CheckIfUserAlreadyDidTheTask(user User, task Task) bool {
	user = db.GetUserByUsername(user)

	taskUserSolution := new(TaskUserSolution)

	err := db.conn.Where("user_id = ? AND task_id = ?", user.ID, task.ID).First(taskUserSolution).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (db *OrmConnection) GetUser(user User) User {
	existingUser := new(User)

	db.conn.Where("username = ? AND password = ?", user.Username, user.Password).First(existingUser)

	return *existingUser
}

func (db *OrmConnection) CreateNewTaskUserSolution(solution TaskUserSolution) TaskUserSolution {
	db.conn.Create(&solution)

	return solution
}
