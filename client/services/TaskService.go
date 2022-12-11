package services

func TaskLoop(user User) {
	for {
		task := GetTask(user)
		if task == nil {
			return
		}
		SolveBackPackTask(task)
		//return
	}
}
