package repository



type AuthRepository interface {
	IdentificationRepo(username, password string) bool
}