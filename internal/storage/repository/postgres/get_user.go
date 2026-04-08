package postgres




type PostgresAuthRepo struct {

}




func (p *PostgresAuthRepo)IdentificationRepo(username, password string) bool {
	userInit := map[string]string{"Alexey": "1234"}
	userPassword := userInit[username]

	if userPassword == password {
		return true
	} else {
		return false
	}
}
