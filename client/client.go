package main

import (
	"client/services"
	"flag"
	"fmt"
)

func main() {
	register := flag.Bool("r", false, "register new account")
	flag.Parse()

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s: %s\n", f.Name, f.Value)
	})

	if *register {
		services.RegisterNewUser()
	}
	//services.GetAvailableServices()
}
