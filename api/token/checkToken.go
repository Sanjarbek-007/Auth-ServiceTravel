package token

import (
	"errors"

	"github.com/form3tech-oss/jwt-go"
	_ "github.com/form3tech-oss/jwt-go"
)

var secret_key = "salom"

func ExtractClaim(tokenStr string) (jwt.MapClaims, error) {

	var (
		token *jwt.Token
		err   error
	)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
	}
	token, err = jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
