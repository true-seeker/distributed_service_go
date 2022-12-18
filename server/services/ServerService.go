package services

var TaskIterationCount = 3
var DefaultTaskSize = 5000

func GenerateTask(taskSize int) Task {
	db := GetDBConnection()
	dbInstance, _ := db.conn.DB()
	defer dbInstance.Close()

	qc := getQueueConnection()
	defer qc.conn.Close()
	defer qc.channel.Close()
	defer qc.ctxCancel()

	newTask := GenerateRandomTask(taskSize)
	task := db.SaveNewTaskParts(newTask)
	qc.PutNewTasksInQueue(task)

	return task
}
