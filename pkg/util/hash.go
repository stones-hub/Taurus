package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// @function: MD5V
// @description: md5加密
// @param: encrypt string 待加密的字符串
// @param: b ...byte 可变参数，用于hash附加的值
// @return: string
func MD5V(encrypt string, b ...byte) string {
	h := md5.New()
	h.Write([]byte(encrypt))
	return hex.EncodeToString(h.Sum(b))
}

// MD5 计算字节流的MD5值
func MD5(b []byte) string {
	h := md5.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// MD5String 计算字符串的MD5值
func MD5String(s string) string {
	return MD5([]byte(s))
}

// SHA1 计算字节流的SHA1值
func SHA1(b []byte) string {
	h := sha1.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA1String 计算字符串的SHA1值
func SHA1String(s string) string {
	return SHA1([]byte(s))
}

// Byte2Hex 将字节流转换为十六进制字符串
func Byte2Hex(b []byte) string {
	return hex.EncodeToString(b)
}

// Hex2Byte 将十六进制字符串转换为字节流
func Hex2Byte(h string) ([]byte, error) {
	return hex.DecodeString(h)
}

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和 存储的 bcrypt 哈希值 是否一致
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
