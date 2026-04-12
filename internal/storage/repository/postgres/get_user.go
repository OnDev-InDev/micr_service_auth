package postgres

// что здесь писать - модель для юсера ?
type PostgresAuthRepo struct{
	db *sql.DB
}

func (p *PostgresAuthRepo) GetUser(username string) bool {
	userInit := map[string]string{"Alexey": "1234"}
	userPassword := userInit[username]

	// if userPassword == password {
	// 	return true
	// } else {
	// 	return false
	// }
}
