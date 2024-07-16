package token

import (
	"errors"
	"time"

	pb "Auth-Service/genproto"

	"github.com/form3tech-oss/jwt-go"
	"go.uber.org/zap"
)

var logger *zap.Logger

var secret_key = "my_secret_key"

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

func GENERATEJWTToken(user *pb.LoginResponse) (*pb.Token, error) {
	AccessToken := jwt.New(jwt.SigningMethodHS256)

	claims := AccessToken.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["full_name"] = user.FullName
	claims["email"] = user.Email
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(3 * time.Hour).Unix()

	access, err := AccessToken.SignedString([]byte(secret_key))
	if err != nil {
		logger.Error("error signing access token", zap.Error(err))
		return nil, err
	}

	RefreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaim := RefreshToken.Claims.(jwt.MapClaims)
	refreshClaim["username"] = user.Username
	refreshClaim["full_name"] = user.FullName
	refreshClaim["email"] = user.Email
	refreshClaim["iat"] = time.Now().Unix()
	refreshClaim["exp"] = time.Now().Add(24 * time.Hour).Unix()

	refresh, err := RefreshToken.SignedString([]byte(secret_key))
	if err != nil {
		logger.Error("error signing refresh token", zap.Error(err))
		return nil, err
	}

	return &pb.Token{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
