package postgres

// что здесь писать - модель для юсера ?
type PostgresAuthRepo struct{}

func (p *PostgresAuthRepo) IdentifyRepo(username, password string) bool {
	userInit := map[string]string{"Alexey": "1234"}
	userPassword := userInit[username]

	if userPassword == password {
		return true
	} else {
		return false
	}
}
