package utils

import (
	"fmt"
	"github.com/google/uuid"
)

func GenerateUUIDV7() (string, error) {
	uuidNew, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("-> uuid.NewV7:%s", err)
	}
	return uuidNew.String(), err
}

func Validate(UUID string) error {
	return uuid.Validate(UUID)
}
