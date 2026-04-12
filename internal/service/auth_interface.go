package service

import "context"

type AuthRepository interface {
	IdentifyRepo(username, password string) bool
}	




