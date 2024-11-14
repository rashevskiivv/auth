package usecase

import (
	"net/mail"
	env "tax-auth/internal"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func ValidateEmail(email string) (bool, error) {
	address, err := mail.ParseAddress(email)
	if err != nil {
		return false, err
	}
	return address.Address != "", nil
}

func GenerateHashedPassword(password string) (*string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}
	hashedPassword := string(bytes)
	return &hashedPassword, nil
}

func GetJWTToken(email string) (*string, error) {
	secretKey, err := env.GetJWTSecretKey()
	if err != nil {
		return nil, err
	}
	payload := jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(secretKey))
	return &t, err
}
