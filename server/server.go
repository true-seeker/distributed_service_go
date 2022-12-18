package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"server/services"
	"syscall"
	"time"
)

func main() {
	//Задаем сид генерации случайных чисел
	rand.Seed(time.Now().UnixNano())

	//При появлении SIGTERM уведомлять канал
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	//Накатывание миграций БД
	services.Migrate()

	//Регистрация сервиса в консуле
	service := services.RegisterService()
	defer service.DeregisterServices()
	fmt.Printf("Successfully registered Consul service with name %s", service.Name)

	//Старт gRPC и веб серверов
	go services.StartGRPCListener()
	go services.StartWebServerListener()

	//Дерегистрация сервисов, при завершении программы
	select {
	case <-c:
		service.DeregisterServices()
	}
}
