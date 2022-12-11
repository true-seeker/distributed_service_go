package services

import (
	"client/backpackTaskGRPC"
)

func TaskLoop(user User) {
	for {
		task := GetTask(user)
		if task == nil {
			return
		}

		ans := SolveBackPackTask(task)
		taskAnswer := backpackTaskGRPC.TaskAnswer{
			TotalPrice: ans,
			User: &backpackTaskGRPC.User{
				Username: user.Username,
				Password: user.Password,
			},
			TaskId: task.Id,
		}
		SendAnswer(&taskAnswer)
	}
}
