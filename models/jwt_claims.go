package models

import "github.com/dgrijalva/jwt-go"

// JWTClaims JWT实体信息
type JWTClaims struct {
	StudentID  string `json:"student_id"`
	ClassName  string `json:"class_name"`
	Permission uint   `json:"permission"`
	jwt.StandardClaims
}
