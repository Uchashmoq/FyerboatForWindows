package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"hash"
)

const (
	KeyLen = 32
)

type Cipher struct {
	n0   string
	key  []byte
	hash hash.Hash32
}

// AesEncrypt 原数据	32字节秘钥，16字节初始化向量
func AesEncrypt(plaintext, key, iv []byte) []byte {
	// 创建一个AES加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// 创建一个填充块
	padding := block.BlockSize() - len(plaintext)%block.BlockSize()
	paddedPlaintext := append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...)
	// 创建一个CBC模式的加密器
	mode := cipher.NewCBCEncrypter(block, iv)
	// 加密数据
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)
	return ciphertext
}

func AesDecrypt(ciphertext, key, iv []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	decrypter := cipher.NewCBCDecrypter(block, iv)
	// 解密数据
	decryptedText := make([]byte, len(ciphertext))
	decrypter.CryptBlocks(decryptedText, ciphertext)
	// 去除填充
	decryptedText = removePadding(decryptedText)
	return decryptedText
}

func removePadding(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
