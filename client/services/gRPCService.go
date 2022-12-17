package services

import (
	"client/backpackTaskGRPC"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"os"
	"strings"
	"time"
)

func getGrpcConnection() (*grpc.ClientConn, error) {
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

func gRPCRegister(username string, password string) {
	conn, err := getGrpcConnection()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
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
		fmt.Println("Registration successful")
	} else if answer.Code == 400 {
		fmt.Printf("User with username: %s already exists\n", username)
	}
}

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
			fmt.Println("Wrong credentials")
			os.Exit(0)
		}
		return nil, errors.New("no available tasks")
	}
	return task, nil
}

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
		fmt.Println("Solution submitted")
		return nil
	} else if response.Code == 403 {
		fmt.Println("Wrong credentials")
		os.Exit(0)
		return nil
	} else {
		return errors.New("cant submit task")
	}
}
