package services

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"server/backpackTaskGRPC"
	"sync"
)

type backpackTaskServer struct {
	backpackTaskGRPC.BackpackTaskServer
	mu sync.Mutex
}

func StartGRPCListener() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s",
		GetProperty("gRPC", "server_port")))

	FailOnError(err, "failed to listen")

	grpcServer := grpc.NewServer()
	backpackTaskGRPC.RegisterBackpackTaskServer(grpcServer, &backpackTaskServer{})
	fmt.Println("Listener started")
	grpcServer.Serve(lis)
}

func (s *backpackTaskServer) Register(ctx context.Context, user *backpackTaskGRPC.User) (*backpackTaskGRPC.Response, error) {
	newUser := User{
		Username: user.Username,
		Password: user.Password,
	}

	err := RegisterNewUser(newUser)
	if err != nil {
		return &backpackTaskGRPC.Response{
			Code:    400,
			Message: "User already exists",
		}, nil
	}
	return &backpackTaskGRPC.Response{
		Code:    200,
		Message: "ok",
	}, nil
}

func (s *backpackTaskServer) GetTask(ctx context.Context, user *backpackTaskGRPC.User) (*backpackTaskGRPC.Task, error) {
	ormUser := User{
		Username: user.Username,
		Password: user.Password,
	}

	if !AuthenticateUser(ormUser) {
		return nil, errors.New("wrong credentials")
	}

	fmt.Println("GetTask", GetMessageCountFromChannel())
	for i := 0; i < GetMessageCountFromChannel(); i++ {
		task := GetTaskPartFromQueue()
		if task == nil {
			return nil, errors.New("no tasks are available")
		}
		fmt.Println("task.Id=", task.ID)
		if CheckIfUserAlreadyDidTheTask(ormUser, *task) {
			fmt.Println("Already did")
			PutTaskInQueue(*task, queueConnection{})
			continue
		}
		var grpcItems []*backpackTaskGRPC.Item
		for _, item := range task.Items {
			grpcItems = append(grpcItems, &backpackTaskGRPC.Item{
				Id:     int32(item.ID),
				Weight: item.Weight,
				Price:  item.Price,
			})
		}
		grpcTask := backpackTaskGRPC.Task{
			Id:               int32(task.ID),
			Items:            grpcItems,
			BackpackCapacity: task.BackpackCapacity,
		}
		fmt.Println(grpcTask)
		return &grpcTask, nil
	}
	return nil, errors.New("no tasks are available")
}
