package services

import (
	"fmt"
)

type User struct {
	Username string
	Password string
}

// RegisterNewUser Регистрация пользователя
func RegisterNewUser() {
	var username string
	var password string
	fmt.Print("Введите имя пользователя: ")
	fmt.Scan(&username)

	fmt.Print("Введите пароль: ")
	fmt.Scan(&password)

	gRPCRegister(username, password)

}
