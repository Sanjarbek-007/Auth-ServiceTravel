package token

import (
	user "Auth-Service/genproto/users"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	signingkey = "my_secret_key"
)

func GeneratedAccessJWTToken(req *user.RegisterResponse, tok *user.Token) error {

	token := *jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = req.Id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(30 * time.Minute).Unix()

	newToken, err := token.SignedString([]byte(signingkey))
	if err != nil {
		log.Println(err)
		return err
	}

	tok.AccessToken = newToken
	return nil
}

func ValidateAccessToken(tokenStr string) (bool, error) {
	_, err := ExtractAccessClaim(tokenStr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractAccessClaim(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(signingkey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, err
	}

	return &claims, nil
}

func GetUserIdFromAccessToken(accessTokenString string) (string, error) {
	refreshToken, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (interface{}, error) { return []byte(signingkey), nil })
	if err != nil || !refreshToken.Valid {
		return "", err
	}
	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}
	userID := claims["user_id"].(string)

	return userID, nil
}


