package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type jwtCustomClaims struct {
	jwt.StandardClaims

	// 追加自己需要的信息
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
}

func GenJwt(userid int64, username, secret string) string {
	claims := &jwtCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Hour * 2).Unix()),
			Issuer:    username,
		},
		userid,
		username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtStr, _ := token.SignedString([]byte(secret))
	return jwtStr
}

func CheckJwt(jwtStr, secret string) (int64, error) {
	args := strings.Split(jwtStr, " ")
	if len(args) == 2 && args[0] == "token" {
		jwtStr = args[1]
	} else {
		return -1, errors.New("Wrong type of jwt token!")
	}
	claims := &jwtCustomClaims{}
	token, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return -1, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return -1, err
	}
	return claims.Uid, token.Claims.Valid()
}
