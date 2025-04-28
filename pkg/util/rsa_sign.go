package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// ParsePriKEY 解析私钥
func ParsePriKEY(priKeyStr string) (*rsa.PrivateKey, error) {

	block, _ := pem.Decode([]byte(priKeyStr))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	// 解析私钥内容
	if privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if iprivKey, err := x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		} else {
			return iprivKey.(*rsa.PrivateKey), nil
		}
	} else {
		return privKey, nil
	}
}

// ParsePubKEY 解析公钥
func ParsePubKEY(pubKeyStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubKeyStr))
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return parsedKey.(*rsa.PublicKey), nil
}

// SignWithRSA  待签名字符串， 私钥， 签名算法  => 返回签名后的数据 .
// 签名算法有 : crypto.MD5  | crypto.SHA256 | crypto.SHA1 | crypto.SHA512 很多，根据需要选择即可
func SignWithRSA(message string, privateKey *rsa.PrivateKey, hashAlg crypto.Hash) ([]byte, error) {
	h := hashAlg.New()
	h.Write([]byte(message))
	hashed := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, hashAlg, hashed)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %v", err)
	}

	return signature, nil
}

// VerifyWithRSA  待验签字符串，公钥，签名后字符串，签名算法.
// 签名算法有 : crypto.MD5  | crypto.SHA256 | crypto.SHA1 | crypto.SHA512 很多，根据需要选择即可
func VerifyWithRSA(originalMessage string, publicKey *rsa.PublicKey, signature []byte, hashAlg crypto.Hash) error {
	h := hashAlg.New()
	h.Write([]byte(originalMessage))
	hashed := h.Sum(nil)

	err := rsa.VerifyPKCS1v15(publicKey, hashAlg, hashed, signature)
	if err != nil {
		return fmt.Errorf("failed to verify signature: %v", err)
	}

	return nil
}
