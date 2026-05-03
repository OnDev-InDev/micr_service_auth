package pkg

import "github.com/google/uuid"

func GenerateID() string {
	newID := uuid.New().String()
	return newID
}
