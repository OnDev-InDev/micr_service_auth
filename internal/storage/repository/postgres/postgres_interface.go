package postgres



type AuthRepository interface {
	IdentificationRepo(username, password string) bool
}