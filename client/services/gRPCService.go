package services

import (
	"client/backpackTaskGRPC"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func gRPCRegister(username string, password string) {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s",
		GetProperty("gRPC", "server_address"),
		GetProperty("gRPC", "server_port")),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	FailOnError(err, "failed to dial")
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
