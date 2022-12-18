package services

// AuthenticateUser Аутентификация пользователя
func (db *OrmConnection) AuthenticateUser(user User) bool {
	existingUser := db.GetUser(user)
	return existingUser.Username != ""
}
