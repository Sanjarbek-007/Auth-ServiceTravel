package postgres

import (
	"fmt"
	pb "Auth-Servicegenproto"
	_ "github.com/form3tech-oss/jwt-go"
	"time"
)

var secret_key = "salom"

type Token struct {
	AccessToken  string
	RefreshToken string
}

func (repo *UserRepository) GENERATEJWTToken(user *pb.LoginResponse) (*Token, error) {
	AccessToken := jwt.New(jwt.SigningMethodHS256)

	claims := AccessToken.Claims.(jwt.MapClaims)
	claims["user_name"] = user.UserName
	claims["password"] = user.Password
	claims["email"] = user.Email
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(3 * time.Hour).Unix()

	access, err := AccessToken.SignedString([]byte(secret_key))
	if err != nil {
		fmt.Println("error  access singed my_secret key")
	}
	RefreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaim := RefreshToken.Claims.(jwt.MapClaims)
	refreshClaim["user_name"] = user.UserName
	refreshClaim["password"] = user.Password
	refreshClaim["email"] = user.Email
	refreshClaim["iat"] = time.Now().Unix()
	refreshClaim["exp"] = time.Now().Add(24 * time.Hour)
	refresh, err := RefreshToken.SignedString([]byte(secret_key))
	if err != nil {
		return nil, err
	}
	_, err = repo.Db.Exec("update users set  token=$1 where email =$2", refresh, user.Email)
	if err != nil {
		return nil, err
	}
	return &Token{
		AccessToken: access,
	}, nil

}
