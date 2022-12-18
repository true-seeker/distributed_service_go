package services

import (
	"client/backpackTaskGRPC"
	"fmt"
	"time"
)

// TaskLoop Цикл решения задач
func TaskLoop(user User) {
	fmt.Println("Начали цикл решения задач. Для выхода: Ctrl+C")
	isTaskExists := true
	noTaskMessage := true
	for {
		//Проверяем есть ли кэширование сообщение
		cashedAnswer, err := GetCashedAnswer()

		//Если есть кэшированная задача, то отправим ее на сервер
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

		// Получение очередной задачи
		task, err := GetTask(user)
		// Если задачу не нашли
		if task == nil {
			if err != nil && noTaskMessage && isTaskExists {
				fmt.Println("Нет доступных задач. Ждем появления новой")
				noTaskMessage = false
			} else if isTaskExists && noTaskMessage {
				fmt.Println("Не можем подключиться к серверу. Переподключаемся")
				isTaskExists = false
			}
			time.Sleep(5 * time.Second)
			GetAvailableServices()
			continue
		}
		isTaskExists = true
		noTaskMessage = true

		// Решение полученной задачи
		ans := SolveBackPackTask(task)
		taskAnswer := backpackTaskGRPC.TaskAnswer{
			TotalPrice: ans,
			User: &backpackTaskGRPC.User{
				Username: user.Username,
				Password: user.Password,
			},
			TaskId: task.Id,
		}

		// Отправляем ответ
		err = SendAnswer(&taskAnswer)
		// Если при отправке ответа возникла ошибка, то кэшируем ответ
		if err != nil {
			SaveSolvedTask(&taskAnswer)
		}
	}
}
