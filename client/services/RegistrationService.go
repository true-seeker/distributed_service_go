package services

import (
	"fmt"
)

type User struct {
	Username string
	Password string
}

func RegisterNewUser() {
	var username string
	var password string
	fmt.Print("Enter username: ")
	fmt.Scan(&username)

	fmt.Print("Enter password: ")
	fmt.Scan(&password)

	gRPCRegister(username, password)

}
