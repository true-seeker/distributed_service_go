package services

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net"
	"server/backpackTaskGRPC"
	"sync"
)

type backpackTaskServer struct {
	backpackTaskGRPC.BackpackTaskServer
	mu sync.Mutex
}

// StartGRPCListener Старт gRPC сервера
func StartGRPCListener() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s",
		GetProperty("gRPC", "server_port")))

	FailOnError(err, "failed to listen gRPC")

	grpcServer := grpc.NewServer()
	backpackTaskGRPC.RegisterBackpackTaskServer(grpcServer, &backpackTaskServer{})
	fmt.Println("gRPC listener started")
	err = grpcServer.Serve(lis)
	FailOnError(err, "failed to start gRPC listener")

}

// Register gRPC метод - регистрация пользователя
func (s *backpackTaskServer) Register(ctx context.Context, user *backpackTaskGRPC.User) (*backpackTaskGRPC.Response, error) {
	newUser := User{
		Username: user.Username,
		Password: user.Password,
	}

	db := GetDBConnection()
	dbInstance, _ := db.conn.DB()
	defer dbInstance.Close()

	err := db.RegisterNewUser(newUser)
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

// GetTask gRPC метод - получение новой задачи
func (s *backpackTaskServer) GetTask(ctx context.Context, user *backpackTaskGRPC.User) (*backpackTaskGRPC.Task, error) {
	ormUser := User{
		Username: user.Username,
		Password: user.Password,
	}

	db := GetDBConnection()
	dbInstance, _ := db.conn.DB()
	defer dbInstance.Close()

	if !db.AuthenticateUser(ormUser) {
		return nil, errors.New("wrong credentials")
	}

	qc := GetQueueConnection()
	defer qc.conn.Close()
	defer qc.channel.Close()
	defer qc.ctxCancel()

	// Если сообщений в очереди нет, то сгенерируем новую задачу
	if qc.GetMessageCountFromChannel() == 0 {
		GenerateTask(DefaultTaskSize)
	}

	for i := 0; i < qc.GetMessageCountFromChannel(); i++ {
		// Получение задачи из очереди
		task := qc.GetTaskFromQueue()
		if task == nil {
			GenerateTask(DefaultTaskSize)
			task = qc.GetTaskFromQueue()
		}
		// Если эту задачу пользователь уже делал
		if db.CheckIfUserAlreadyDidTheTask(ormUser, *task) {
			qc.PutTaskInQueue(*task)
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
		return &grpcTask, nil
	}
	return nil, errors.New("no tasks are available")
}

// SendAnswer Отправить ответ на задачу
func (s *backpackTaskServer) SendAnswer(ctx context.Context, answer *backpackTaskGRPC.TaskAnswer) (*backpackTaskGRPC.Response, error) {
	user := User{
		Username: answer.User.Username,
		Password: answer.User.Password,
	}
	db := GetDBConnection()
	dbInstance, _ := db.conn.DB()
	defer dbInstance.Close()

	if !db.AuthenticateUser(user) {
		return &backpackTaskGRPC.Response{
			Code:    403,
			Message: "wrong credentials",
		}, nil
	}
	user = db.GetUserByUsername(user)
	task := Task{
		Model: gorm.Model{ID: uint(answer.TaskId)},
	}

	if db.CheckIfUserAlreadyDidTheTask(user, task) {
		return &backpackTaskGRPC.Response{
			Code:    400,
			Message: "Current user already did the task",
		}, nil
	}

	solution := TaskUserSolution{
		TaskId: uint(answer.TaskId),
		UserId: user.ID,
		Answer: answer.TotalPrice,
	}

	fmt.Println("solution.TaskId", solution.TaskId)
	solution = db.CreateNewTaskUserSolution(solution)

	return &backpackTaskGRPC.Response{
		Code:    200,
		Message: "ok",
	}, nil
}
