package tool

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256PasswordWithSalt 密码加盐hash10000次函数
func Sha256PasswordWithSalt(password string) string {
	salt := []byte(Setting.PasswordSalt)
	result := sha256.Sum256(append([]byte(password), salt...))

	for i := 0; i < 9999; i++ {
		result = sha256.Sum256(append(result[:], salt...))
	}

	return hex.EncodeToString(result[:])
}
