package django

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func GeneratePassword(password string) string {
	// 生成随机的盐值
	BASE_STR := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	salt := ""
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 22; i++ {
		salt += string(BASE_STR[rand.Intn(len(BASE_STR))])
	}
	// 使用PBKDF2算法将密码哈希化
	hashedPassword := pbkdf2.Key([]byte(password), []byte(salt), 600000, 32, sha256.New)
	encodedHashedPassword := base64.StdEncoding.EncodeToString(hashedPassword)
	// 构建Django格式的密码字符串
	djangoPassword := fmt.Sprintf("pbkdf2_sha256$600000$%s$%s", salt, encodedHashedPassword)
	//
	return djangoPassword
}

func CheckPassword(password, djangoPassword string) bool {
	// 解析Django密码哈希字符串
	djangoParts := strings.Split(djangoPassword, "$")
	iterations := 0
	fmt.Sscanf(djangoParts[1], "%d", &iterations)
	salt := djangoParts[2]
	hashedPassword, _ := base64.StdEncoding.DecodeString(djangoParts[3])
	// 生成与Django密码哈希字符串相匹配的密码哈希值
	newHashedPassword := pbkdf2.Key([]byte(password), []byte(salt), 600000, len(hashedPassword), sha256.New)
	// fmt.Println("New hashed password: ", base64.StdEncoding.EncodeToString(newHashedPassword))
	return subtle.ConstantTimeCompare(hashedPassword, newHashedPassword) == 1
}
