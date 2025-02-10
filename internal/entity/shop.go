package entity

import (
	"fmt"
	"regexp"
)

type UserInfo struct {
}

type PurchaseInfo struct {
}

type TransferInfo struct {
}

type UserCredentials struct {
	identifier string
	password   string
}

func (ru UserCredentials) Identifier() string {
	return ru.identifier
}

func (ru UserCredentials) Password() string {
	return ru.password
}

func NewUserCredentials(identifier string, password string) (*UserCredentials, error) {
	ru := &UserCredentials{identifier: identifier, password: password}
	err := ru.validate()

	if err != nil {
		return nil, fmt.Errorf("-> ru.validate%v", err)
	}

	return &UserCredentials{identifier: identifier, password: password}, nil
}

func (ru UserCredentials) validate() error {
	err := validateIdentifier(ru.identifier)
	if err != nil {
		return fmt.Errorf("-> validateIdentifier%v", err)
	}

	err = validatePassword(ru.password)
	if err != nil {
		return fmt.Errorf("-> validatePassword%v", err)
	}

	return nil
}

func validateIdentifier(identifier string) error {
	// от 4 до 50 символов: буквы, цифры, точки, дефисы и симовлы подчеркивания
	identifierRegex := `^[a-zA-Z0-9._-]{4,50}$`

	matched, err := regexp.MatchString(identifierRegex, identifier)
	if err != nil {
		return fmt.Errorf("-> regexp.MatchString%v", err)
	}

	if !matched {
		return fmt.Errorf(": идентификатор должен быть длинной 4-50 и содержать только буквы, цифры, точки, дефисы и симовлы подчеркивания")
	}

	return nil
}

func validatePassword(password string) error {
	// от 8 до 64 символов: хотя бы одну букву, одну цифру и один спецсимвол

	passwordRegex := `^(?=.*[a-zA-Z])(?=.*\d)(?=.*[@$!%*?&])[a-zA-Z\d@$!%*?&]{8,64}$`

	matched, err := regexp.MatchString(passwordRegex, password)
	if err != nil {
		return fmt.Errorf("-> regexp.MatchString%v", err)
	}

	if !matched {
		return fmt.Errorf(": идентификатор должен быть длинной 8-64 и содержать хотя бы одну букву, одну цифру и один спецсимвол")
	}

	return nil

}
