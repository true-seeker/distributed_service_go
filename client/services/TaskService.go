package services

func TaskLoop(user User) {
	for {
		if GetTask(user) == nil {
			break
		}
	}
}
