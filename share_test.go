package gox_test

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"github.com/infavor/gox"
	"github.com/infavor/gox/convert"
	"github.com/infavor/gox/logger"
	"testing"
	"time"
)

func init() {
	logger.Init(nil)
}

func TestNetwork(t *testing.T) {
	logger.Info(gox.GetMyAddress("vEthernet", "192.168.0"))
}

func DesEncryption(key, iv, plainText []byte) ([]byte, error) {

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData := PKCS5Padding(plainText, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return cryted, nil
}

func DesDecryption(key, iv, cipherText []byte) ([]byte, error) {

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cipherText)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func TestDES(t *testing.T) {
	originalText := convert.Int64ToStr(time.Now().UnixNano()) + "|123456789012345678901234567890df"
	fmt.Println(originalText)
	mytext := []byte(originalText)

	key := []byte{0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC}
	iv := []byte{0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC}

	cryptoText, _ := DesEncryption(key, iv, mytext)
	base64String := base64.StdEncoding.EncodeToString(cryptoText)
	//																// fDE1NjExMDAwOTY5NTM4MzI2MDD10pSRpEgO4DZti3M2w/YkYNKl0TvWxyQ=
	//base64String := base64.StdEncoding.EncodeToString(cryptoText) // 9dKUkaRIDuA2bYtzNsP2JGDSpdE71sck
	fmt.Println(base64String)
	bs, _ := base64.StdEncoding.DecodeString(base64String)
	decryptedText, _ := DesDecryption(key, iv, bs)
	fmt.Println(string(decryptedText))
}

func TestGetHumanReadableDuration(t *testing.T) {
	createTime := gox.CreateTime(gox.GetTimestamp(time.Now()) - 420000)
	fmt.Println(gox.GetHumanReadableDuration(createTime, time.Now()))
	fmt.Println(gox.GetLongHumanReadableDuration(createTime, time.Now()))
}
