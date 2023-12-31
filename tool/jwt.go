package tool

import (
	"ddlBackend/models"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWTToken 生成JWT令牌
func GenerateJWTToken(info models.UserInformation) (string, error) {
	// 设置token的有效时间为24小时
	expireTime := time.Now().Add(24 * time.Hour)

	// 设置token中的信息
	claims := models.JWTClaims{
		StudentID:  info.StudentID,
		ClassName:  info.ClassName,
		Permission: info.Permission,
		RegisteredClaims: jwt.RegisteredClaims{
			// token失效时间
			ExpiresAt: jwt.NewNumericDate(expireTime),
			// token签发人
			Issuer: "SquidWard",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte(Setting.JWTSecret))
}

// ParseJWTToken 解析令牌并返回其中的信息
func ParseJWTToken(token string) (*models.JWTClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Setting.JWTSecret), nil
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

// GetClaimsInContext 获得HTTP上下文中的JWT令牌信息
func GetClaimsInContext(context *gin.Context) (*models.JWTClaims, error) {
	value, exist := context.Get("Claims")
	if !exist {
		// 没有找到令牌信息
		return nil, errors.New("no JWT claims")
	}

	claims, ok := value.(models.JWTClaims)
	if !ok {
		return nil, errors.New("can not read claims")
	}

	return &claims, nil
}

// CheckPermission 验证请求者的权限
func CheckPermission(context *gin.Context, permission uint) (bool, error) {
	claims, err := GetClaimsInContext(context)
	if err != nil {
		return false, err
	}

	if claims.Permission >= permission {
		return true, nil
	} else {
		return false, nil
	}
}
