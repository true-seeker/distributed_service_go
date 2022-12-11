package services

func AuthenticateUser(user User) bool {
	existingUser := GetUser(user)
	return existingUser.Username != ""
}
