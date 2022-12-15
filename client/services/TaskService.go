package services

import (
	"client/backpackTaskGRPC"
	"fmt"
	"time"
)

func TaskLoop(user User) {
	fmt.Println("Started task loop. To exit press ctrl+C")
	isTaskExists := true
	noTaskMessage := true
	for {
		cashedAnswer, err := GetCashedAnswer()

		if err == nil {
			cashedAnswer.User = &backpackTaskGRPC.User{
				Username: user.Username,
				Password: user.Password,
			}
			err = SendAnswer(cashedAnswer)
			if err != nil {
				SaveSolvedTask(cashedAnswer)
			}
		}
		task, err := GetTask(user)
		if task == nil {
			if err != nil && noTaskMessage && isTaskExists {
				fmt.Println("No available tasks. Waiting for new task to appear")
				noTaskMessage = false
			} else if isTaskExists && noTaskMessage {
				fmt.Println("Cant connect to server. Trying to reconnect")
				isTaskExists = false
			}
			time.Sleep(5 * time.Second)
			GetAvailableServices()
			continue
		}
		isTaskExists = true
		noTaskMessage = true

		ans := SolveBackPackTask(task)
		taskAnswer := backpackTaskGRPC.TaskAnswer{
			TotalPrice: ans,
			User: &backpackTaskGRPC.User{
				Username: user.Username,
				Password: user.Password,
			},
			TaskId: task.Id,
		}
		err = SendAnswer(&taskAnswer)
		if err != nil {
			SaveSolvedTask(&taskAnswer)
		}
	}
}
