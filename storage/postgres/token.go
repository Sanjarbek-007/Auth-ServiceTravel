package postgres

import (
	pb "Auth-Service/genproto"
	"time"

	_ "github.com/form3tech-oss/jwt-go"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

var secret_key = "salom"

type Token struct {
	AccessToken  string
	RefreshToken string
}

var logge *zap.Logger

func (repo *UserRepository) GENERATEJWTToken(user *pb.LoginResponse) (*pb.Token, error) {
	AccessToken := jwt.New(jwt.SigningMethodHS256)

	claims := AccessToken.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["full_name"] = user.FullName
	claims["email"] = user.Email
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(3 * time.Hour).Unix()

	access, err := AccessToken.SignedString([]byte(secret_key))
	if err != nil {
		logge.Error("error signing access token", zap.Error(err))
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
		logge.Error("error signing refresh token", zap.Error(err))
		return nil, err
	}

	_, err = repo.Db.Exec("UPDATE users SET token=$1 WHERE email=$2 AND username=$3", refresh, user.Email, user.Username)
	if err != nil {
		logge.Error("error updating user with new refresh token", zap.Error(err))
		return nil, err
	}

	return &pb.Token{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
