package token

import (
	pb "Auth-Service/genproto/users"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	newsigningkey = "my_secret_key"
)

func GeneratedRefreshJWTToken(req *pb.RegisterResponse, tok *pb.Token) error {
	token := *jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = req.Id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	newToken, err := token.SignedString([]byte(newsigningkey))
	if err != nil {
		return err
	}

	tok.RefreshToken = newToken
	return nil
}

func ValidateRefreshToken(tokenStr string) (bool, error) {
	_, err := ExtractRefreshClaim(tokenStr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractRefreshClaim(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(newsigningkey), nil
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

func GetUserIdFromRefreshToken(accessTokenString string) (string, error) {
	refreshToken, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (interface{}, error) { return []byte(newsigningkey), nil })
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
