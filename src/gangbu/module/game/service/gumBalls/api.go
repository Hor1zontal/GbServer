package gumBalls

import (
	"aliens/log"
	"bytes"
	"crypto/des"
	"encoding/hex"
	"encoding/json"
	"gangbu/exception"
	"io/ioutil"
	"net/http"
	"strings"
)

type DrawResp struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

type DrawGiftReq struct {
	Act string `json:"act"`
	Sign string `json:"sign"`
}

func DrawGift(openId string) (string, string) {
	//var sign string
	//sign, err := DesEncrypt([]byte(openId), []byte("gB1903Gm"))
	//cipher.ecb


	//sign := make([]byte, len([]byte(openId)))
	//block, err := des.NewCipher([]byte("gB1903Gm"))
	//block.Encrypt(sign, []byte(openId))

	//log.Debug("sign:%s", string(sign))

	sign := EncryptDesECB([]byte(openId), []byte("gB1903Gm"))
	url := "http://activity.leiting.com/gumballs/201903/game/index.php"
	data := "act=gift&sign=" + sign

	resp, err := http.Post(url, "application/x-www-form-urlencoded;charset=UTF-8", 	strings.NewReader(data))
	if err != nil || resp.StatusCode != http.StatusOK {
		exception.GameExceptionCustom("DrawGumBallsGift",exception.GumBallsGiftError, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exception.GameExceptionCustom("DrawGumBallsGift",exception.GumBallsGiftError, err)
	}
	atr := &DrawResp{}
	json.Unmarshal(body, &atr)
	return atr.Status, atr.Message
}

func EncryptDesECB(data, key []byte) string {
	if len(key) > 8 {
		key = key[:8]
	}
	block, err := des.NewCipher(key)
	if err != nil {
		log.Error("EntryptDesECB newCipher error[%v]", err)
		return ""
	}
	bs := block.BlockSize()
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		log.Error("EntryptDesECB Need a multiple of the blocksize")
		return ""
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return 	hex.EncodeToString(out)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}


//func DesEncrypt(origData, key []byte) ([]byte, error) {
//	block, err := des.NewCipher(key)
//	if err != nil {
//		return nil, err
//	}
//	origData = PKCS5Padding(origData, block.BlockSize())
//	// origData = ZeroPadding(origData, block.BlockSize())
//	//cipher1.ECBEncrypt()
//	//todo
//	blockMode := cipher.NewCBCEncrypter(block, key)
//	crypted := make([]byte, len(origData))
//	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
//	// crypted := origData
//	blockMode.CryptBlocks(crypted, origData)
//	return crypted, nil
//
//}

//func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
//	padding := blockSize - len(ciphertext)%blockSize
//	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
//	return append(ciphertext, padtext...)
//}
//
//func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
//	padding := blockSize - len(ciphertext)%blockSize
//	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
//	return append(ciphertext, padtext...)
//}
