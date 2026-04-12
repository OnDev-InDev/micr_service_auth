package service

type AuthService struct {
	userRepo UserRepository
}

// идентификация  - должен вернуть пароль мне например и я его уже здесь в бизнес логите сравню с тем что ввел пользователь
func (s *AuthService) AuthenticateUser(username, password string) bool {
	if s.userRepo.GetUser(username) {
		return true
	}

	return false
}
