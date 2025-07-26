package Intrastructures

import (
	"log"

	domain "github.com/segnig/task-manager/Domains"
	"golang.org/x/crypto/bcrypt"
)

type PasswordProvider struct {
	cost int
}

func (pp *PasswordProvider) HashPassword(userPassword string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), pp.cost)
	if err != nil {
		log.Println("error hashing password:", err)
		return ""
	}
	return string(hashedPassword)
}

func (pp *PasswordProvider) VerifyPassword(hashedPwd, plainPwd string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return false, "username or password is incorrect"
	}
	return true, ""
}

func NewPasswordProvider(cost int) domain.PasswordServiceProvider {
	return &PasswordProvider{
		cost: cost,
	}
}
