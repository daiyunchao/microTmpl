package lib

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
)

type AESEnCode struct {
}

func (encode *AESEnCode) Encode(src []byte, key string) ([]byte, error) {
	return aesEncrypt(src, []byte(key))
}

func (encode *AESEnCode) Decode(src []byte, key string) ([]byte, error) {
	return aesDecrypt(src, []byte(key))
}

// AesEncrypt 加密
func aesEncrypt(data []byte, key []byte) ([]byte, error) {
	// 密钥和待加密数据转成[]byte
	block, _ := aes.NewCipher(key)
	// 获取密钥长度
	blockSize := block.BlockSize()
	// 补码
	data = PKCS7Padding(data, blockSize)
	// 创建保存加密变量
	encryptResult := make([]byte, len(data))
	// CEB是把整个明文分成若干段相同的小段，然后对每一小段进行加密
	for bs, be := 0, blockSize; bs < len(data); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(encryptResult[bs:be], data[bs:be])
	}
	return []byte(base64.StdEncoding.EncodeToString(encryptResult)), nil
}

// AesDecrypt 解密
func aesDecrypt(data []byte, key []byte) ([]byte, error) {
	// 反解密码base64
	originByte, _ := base64.StdEncoding.DecodeString(string(data))
	// 密钥和待加密数据转成[]byte
	keyByte := []byte(key)
	// 创建密码组，长度只能是16、24、32字节
	block, _ := aes.NewCipher(keyByte)
	// 获取密钥长度
	blockSize := block.BlockSize()
	// 创建保存解密变量
	decrypted := make([]byte, len(originByte))
	for bs, be := 0, blockSize; bs < len(originByte); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(decrypted[bs:be], originByte[bs:be])
	}
	// 解码
	return PKCS7UNPadding(decrypted), nil
}

// PKCS7Padding 补码
func PKCS7Padding(originByte []byte, blockSize int) []byte {
	// 计算补码长度
	padding := blockSize - len(originByte)%blockSize
	// 生成补码
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	// 追加补码
	return append(originByte, padText...)
}

// 解码
func PKCS7UNPadding(originDataByte []byte) []byte {
	length := len(originDataByte)
	unpadding := int(originDataByte[length-1])
	return originDataByte[:(length - unpadding)]
}
