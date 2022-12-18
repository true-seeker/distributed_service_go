package services

import (
	"client/backpackTaskGRPC"
	"encoding/json"
	"errors"
	"os"
)

type AnswerDTO struct {
	TotalPrice uint32
	TaskId     int32
}

// SaveSolvedTask сохранение кэшированного ответа
func SaveSolvedTask(answer *backpackTaskGRPC.TaskAnswer) {
	file, err := os.Create("cashed_tasks.json")
	FailOnError(err, "Cant create cash file")
	file, err = os.OpenFile("cashed_tasks.json", os.O_WRONLY, os.ModeAppend)
	FailOnError(err, "Cant open cash file")
	defer file.Close()
	answerDTO := AnswerDTO{
		TotalPrice: answer.TotalPrice,
		TaskId:     answer.TaskId,
	}
	byteData, _ := json.Marshal(answerDTO)
	_, err = file.Write(byteData)
	FailOnError(err, "Cant write cash file")
}

// GetCashedAnswer Получение кэшированного ответа
func GetCashedAnswer() (*backpackTaskGRPC.TaskAnswer, error) {
	file, err := os.OpenFile("cashed_tasks.json", os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return nil, errors.New("cant open cash file")
	}
	defer file.Close()

	data, err := os.ReadFile("cashed_tasks.json")
	if err != nil {
		return nil, errors.New("cant read cash file")
	}

	answerDto := AnswerDTO{}
	err = json.Unmarshal(data, &answerDto)
	if err != nil {
		return nil, errors.New("cant unmarshal cash file")
	}

	_, err = os.Create("cashed_tasks.json")
	FailOnError(err, "Cant recreate cash file")

	return &backpackTaskGRPC.TaskAnswer{
		TaskId:     answerDto.TaskId,
		TotalPrice: answerDto.TotalPrice,
	}, nil
}
