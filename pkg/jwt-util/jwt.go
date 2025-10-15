package jwt_util

import (
	"errors"
	"time"

	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/dgrijalva/jwt-go"
)

const (
	// jwt subject
	USER = "user"

	UserTokenTimeout      = int64(60 * 60 * 24) // 1天
	BufferTime            = int64(60 * 60 * 2)
	UserLoginTokenTimeout = int64(60 * 5) // 5 min
)

var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("that's not even a token")
	ErrTokenInvalid     = errors.New("couldn't handle this token")
)

var (
	userSecretKey string
)

type CustomClaims struct {
	UserID     string `json:"userId"` // 用户ID
	BufferTime int64  `json:"bufferTime"`
	jwt.StandardClaims
}

func InitUserJWT(key string) {
	if userSecretKey != "" {
		log.Panicf("jwt already init")
	}
	if key == "" {
		log.Panicf("jwt secret key empty")
	}
	userSecretKey = key
}

func GenerateToken(userID string, timeout int64) (*CustomClaims, string, error) {
	return generateToken(userID, timeout, userSecretKey)
}

func ParseToken(token string) (*CustomClaims, error) {
	return parseToken(token, userSecretKey)
}

func generateToken(id string, timeout int64, secretKey string) (*CustomClaims, string, error) {
	if secretKey == "" {
		return nil, "", errors.New("jwt secret key empty")
	}
	nowTime := time.Now().Unix()
	claims := &CustomClaims{
		UserID:     id,
		BufferTime: nowTime + BufferTime, // 缓冲时间，当nowTime大于等于BufferTime and nowTime小于ExpiresAt是获得新的token
		StandardClaims: jwt.StandardClaims{
			Issuer:    "wanwu",
			Subject:   USER,              // 用途，目前固定user
			NotBefore: nowTime,           // 生效时间
			ExpiresAt: nowTime + timeout, // 过期时间
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return nil, "", err
	}
	return claims, token, err
}

func parseToken(token, secretKey string) (*CustomClaims, error) {
	if secretKey == "" {
		return nil, errors.New("jwt secret key empty")
	}
	tokenClaims, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*CustomClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
		return nil, ErrTokenInvalid

	} else {
		return nil, ErrTokenInvalid
	}
}
