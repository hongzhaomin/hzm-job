package tool

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hongzhaomin/hzm-job/admin/internal/consts"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"time"
)

type Claims struct {
	UserId  int64 `json:"userId"`
	Version int64 `json:"version"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userId int64, version int64) (string, error) {
	claims := Claims{
		UserId:  userId,
		Version: version,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.JwtTokenExpiresDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    consts.JwtIssuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(consts.JwtSecret)
}

// ParseToken 验证JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return consts.JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}

// GetUserId 获取当前登录userId
func GetUserId(ctx *gin.Context) int64 {
	return *GetLoginUser(ctx).Id
}

// GetLoginUser 获取当前登录用户
func GetLoginUser(ctx *gin.Context) *vo.User {
	if user, ok := ctx.Get("loginUser"); ok {
		return user.(*vo.User)
	}
	return &vo.User{}
}
