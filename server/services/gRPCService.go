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

func (s *backpackTaskServer) GetTask(ctx context.Context, user *backpackTaskGRPC.User) (*backpackTaskGRPC.TaskPart, error) {
	ormUser := User{
		Username: user.Username,
		Password: user.Password,
	}

	if !AuthenticateUser(ormUser) {
		return nil, errors.New("wrong credentials")
	}

	fmt.Println("GetTask", GetMessageCountFromChannel())
	for i := 0; i < GetMessageCountFromChannel(); i++ {
		taskPart := GetTaskPartFromQueue()
		if taskPart == nil {
			return nil, errors.New("no tasks are available")
		}
		task := GetTaskById(taskPart.TaskId)
		fmt.Println("taskPart.Id=", taskPart.ID)
		if CheckIfUserAlreadyDidTheTask(ormUser, *taskPart) {
			fmt.Println("Already did")
			PutTaskPartInQueue(*taskPart, queueConnection{})
			continue
		}
		var grpcItems []*backpackTaskGRPC.Item

		for _, item := range taskPart.Items {
			grpcItems = append(grpcItems, &backpackTaskGRPC.Item{
				Id:      int32(item.ID),
				Weight:  item.Weight,
				Price:   item.Price,
				IsFixed: item.IsFixed,
			})
		}
		grpcTaskPart := backpackTaskGRPC.TaskPart{
			Id:               int32(taskPart.ID),
			TaskId:           int32(taskPart.TaskId),
			Items:            grpcItems,
			BackpackCapacity: task.BackpackCapacity,
		}
		fmt.Println(grpcTaskPart)
		return &grpcTaskPart, nil
	}
	return nil, errors.New("no tasks are available")
}
