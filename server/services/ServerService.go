package services

var TaskIterationCount = 3
var DefaultTaskSize = 5000

// GenerateTask Сгенерировать задачу
func GenerateTask(taskSize int) Task {
	db := GetDBConnection()
	dbInstance, _ := db.conn.DB()
	defer dbInstance.Close()

	qc := GetQueueConnection()
	defer qc.conn.Close()
	defer qc.channel.Close()
	defer qc.ctxCancel()

	newTask := GenerateRandomTask(taskSize)
	task := db.SaveNewTask(newTask)
	qc.PutNewTasksInQueue(task)

	return task
}
