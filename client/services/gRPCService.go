package services

import (
	"client/backpackTaskGRPC"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"strings"
	"time"
)

// getGrpcConnection Подключение к серверу gRPC
func getGrpcConnection() (*grpc.ClientConn, error) {
	isFirstMessage := true
	for len(AvailableServices) == 0 {
		if isFirstMessage {
			fmt.Println("Нет доступных серверов. Ждём")
			isFirstMessage = false
		}
		GetAvailableServices()
		time.Sleep(5 * time.Second)
	}

	rand.Shuffle(len(AvailableServices), func(i, j int) {
		AvailableServices[i], AvailableServices[j] = AvailableServices[j], AvailableServices[i]
	})

	for _, AvailableService := range AvailableServices {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d",
			AvailableService.Address,
			AvailableService.Port),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithTimeout(5*time.Second))
		if err == nil {
			return conn, nil
		}
	}
	return nil, errors.New("Failed to connect to grpc server. Try again later")
}

// gRPCRegister gRPC метод - регистрация пользователя
func gRPCRegister(username string, password string) {
	conn, err := getGrpcConnection()
	if err != nil {
		log.Fatalln(err)
	}

	client := backpackTaskGRPC.NewBackpackTaskClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUser := backpackTaskGRPC.User{
		Username: username,
		Password: password,
	}

	answer, err := client.Register(ctx, &newUser)
	FailOnError(err, "failed to send register request")

	if answer.Code == 200 {
		fmt.Println("Регистрация успешна")
	} else if answer.Code == 400 {
		fmt.Printf("Пользователь с именем: %s уже сушествует\n", username)
	}
}

// GetTask gRPC метод - получение задачи
func GetTask(user User) (*backpackTaskGRPC.Task, error) {
	conn, err := getGrpcConnection()
	if err != nil {
		return nil, errors.New("cant connect to server")
	}

	client := backpackTaskGRPC.NewBackpackTaskClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcUser := backpackTaskGRPC.User{
		Username: user.Username,
		Password: user.Password,
	}

	task, err := client.GetTask(ctx, &grpcUser)
	if err != nil {
		if strings.Contains(err.Error(), "wrong credentials") {
			log.Fatalln("Неверный логин или пароль")
		}
		return nil, errors.New("no available tasks")
	}
	return task, nil
}

// SendAnswer gRPC метод - отправка ответа
func SendAnswer(answer *backpackTaskGRPC.TaskAnswer) error {
	conn, err := getGrpcConnection()

	if err != nil {
		return errors.New("cant connect to grpc")
	}

	client := backpackTaskGRPC.NewBackpackTaskClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := client.SendAnswer(ctx, answer)
	FailOnError(err, "Failed To send answer")

	if response.Code == 200 {
		fmt.Println("Ответ принят")
		return nil
	} else if response.Code == 403 {
		log.Fatalln("Неверный логин или пароль")
		return nil
	} else {
		return errors.New("cant submit task")
	}
}
