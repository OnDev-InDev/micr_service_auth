package service



type AuthService struct {
	authRepo AuthRepository
}

// идентификация
func (s *AuthService) AuthenticateUser(username, password string) bool {
	if s.authRepo.IdentifyRepo(username, password) {
		return true
	}

	return false
}
