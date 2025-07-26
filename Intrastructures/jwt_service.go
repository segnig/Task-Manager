package Intrastructures

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	domain "github.com/segnig/task-manager/Domains"
)

type UserToken struct {
	SECTRET_KEY string
}

func NeWUserToken(secretKey string) domain.IUserToken {
	return &UserToken{
		SECTRET_KEY: secretKey,
	}
}

func (ut *UserToken) GenerateAllTokens(username, uid, userType string) (signedToken, signedRefreshToken string, err error) {
	claims := &domain.SignedDetails{
		Username: username,
		Uid:      uid,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	refreshClaims := &domain.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(200 * time.Hour).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(ut.SECTRET_KEY))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(ut.SECTRET_KEY))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

func (ut *UserToken) ValidateToken(signedToken string) (claims *domain.SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&domain.SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(ut.SECTRET_KEY), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.SignedDetails)
	if !ok {
		return nil, fmt.Errorf("the token is in valid")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token is expired")
	}
	return claims, err
}
