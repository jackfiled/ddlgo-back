package tool

import (
	"ddlBackend/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// GenerateJWTToken 生成JWT令牌
func GenerateJWTToken(info models.UserInformation) (string, error) {
	// 设置token的有效时间为24小时
	expireTime := time.Now().Add(24 * time.Hour)

	// 设置token中的信息
	claims := models.JWTClaims{
		Username:   info.Username,
		Classname:  info.Classname,
		Permission: info.Permission,
		StandardClaims: jwt.StandardClaims{
			// token失效时间
			ExpiresAt: expireTime.Unix(),
			// token签发人
			Issuer: "SquidWard",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(Setting.JWTSecret)
}

// ParseJWTToken 解析令牌并返回其中的信息
func ParseJWTToken(token string) (*models.JWTClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Setting.JWTSecret, nil
	})

	if tokenClaims != nil {
		// 尝试判断令牌中的信息是否正确
		claims, ok := tokenClaims.Claims.(*models.JWTClaims)
		// 如果信息正确且令牌有效
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
