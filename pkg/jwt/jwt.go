package jwt

import (
	"health_checker/pkg/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(1 * time.Hour)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	})
	token, err := tokenClaims.SignedString([]byte(config.Conf.Server.Secret))
	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.Server.Secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
