package middleware

import (
	"fmt"
	"ino-chat/res"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const SIGNKEY = "BA5ktbKaV47uOcQpnuUT76GvBRYpMdHX"

type CustomerInfo struct {
	Name string
	Uid  string
}
type CustomClaims struct {
	*jwt.StandardClaims
	*CustomerInfo
}

func CreateToken(uid string, name string) (string, error) {
	key := []byte(SIGNKEY)
	c := CustomClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    name,
		},
		&CustomerInfo{
			Uid:  uid,
			Name: name,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	sign, err := token.SignedString(key)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	return sign, nil
}

func ParseToken(tokenStr string) (*CustomerInfo, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SIGNKEY), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.CustomerInfo, nil
	} else {
		return nil, errors.Wrap(err, "")
	}
}

func AuthMiddleware(c *gin.Context) {
	var token string
	token = c.Request.Header.Get("Authorization")
	if token == "" {
		token, _ = c.Cookie("Authorization")
	}
	if token == "" {
		c.JSON(401, res.Unauthrized("Unauthrized."))
		c.Redirect(301, "/")
		c.Abort()
		return
	}
	user, err := ParseToken(token)
	if err != nil {
		logrus.Warnf("jwt解析失败：%v", err)
		c.JSON(401, res.Unauthrized("Unauthrized."))
		c.Redirect(301, "/")
		c.Abort()
		return
	}
	c.Request.Header.Add("USER_ID", user.Uid)
	c.Request.Header.Add("USER_NAME", user.Name)
	c.Next()
}
