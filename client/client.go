package main

import (
	"client/services"
	"flag"
	"fmt"
)

func main() {
	register := flag.Bool("r", false, "Зарегистрировать новый аккаунт")
	username := flag.String("u", "", "Имя пользователя")
	password := flag.String("p", "", "Пароль")
	flag.Parse()

	// Получение доступных сервисов
	services.GetAvailableServices()
	// Регистрация
	if *register {
		services.RegisterNewUser()
	} else if *username == "" || *password == "" {
		fmt.Println("Не указан логин или пароль\n\n-h для просмотра справки")
	} else {
		// Старт цикла решения задач
		services.TaskLoop(services.User{Username: *username, Password: *password})
	}
}
