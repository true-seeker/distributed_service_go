package services

import (
	"fmt"
)

func RegisterNewUser() {
	var username string
	var password string
	fmt.Print("Enter username: ")
	fmt.Scan(&username)

	fmt.Print("Enter password: ")
	fmt.Scan(&password)

	gRPCRegister(username, password)

}
