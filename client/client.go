package main

import (
	"client/services"
	"flag"
	"fmt"
)

func main() {
	register := flag.Bool("r", false, "register new account")
	username := flag.String("u", "", "username")
	password := flag.String("p", "", "password")
	flag.Parse()

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s: %s\n", f.Name, f.Value)
	})
	fmt.Println()

	if *register {
		services.RegisterNewUser()
	} else if *username == "" || *password == "" {
		fmt.Println("Please, provide username and password\n-h for help")
	} else {
		services.TaskLoop(services.User{Username: *username, Password: *password})
	}
	//services.GetAvailableServices()
}
