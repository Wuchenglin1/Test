package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyClaims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.StandardClaims
}

func main() {
	mySigningKey := []byte("abcdefg")
	c := MyClaims{
		Username: "wuchenglin",
		IsAdmin:  true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 5).Unix(),
			Issuer:    "Admin", //jwt的签发者
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c) //创建一个token

	fmt.Println("t:", t)

	s, err := t.SignedString(mySigningKey) //将密钥输入进去，进行签名

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("s:", s)

	token, err2 := jwt.ParseWithClaims(s, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		fmt.Printf("%v %v", claims.Username, claims.StandardClaims.ExpiresAt)
	} else {
		fmt.Println(err2)
	}

}
