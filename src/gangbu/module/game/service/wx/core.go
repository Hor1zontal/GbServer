package wx

import (
	"aliens/log"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"gangbu/module/game/service/wx/model"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var wxConfig model.Config

var responseMapping map[string]string = make(map[string]string)

func Init(config model.Config) {
	wxConfig = config
	wxConfig.AesKey = EncodingAESKey2AESKey(wxConfig.EncodingAESKey)
	RefreshToken()

	responseMapping["DAY_GIFT"] = "有没有礼包~进游戏就知道了！\n👇👇👇👇"
}

func GetTextResponse(key string) string {
	return responseMapping[key]
}

func EncodingAESKey2AESKey(encodingKey string) []byte {
	data, _ := base64.StdEncoding.DecodeString(encodingKey + "=")
	return data
}


func ValidateUrl(timestamp, nonce, signatureIn string) bool {
	signatureGen := MakeSignature(timestamp, nonce)
	if signatureGen != signatureIn {
		return false
	}
	return true
}

func ValidateMsg(timestamp, nonce, msgEncrypt, msgSignatureIn string) bool {
	msgSignatureGen := MakeMsgSignature(timestamp, nonce, msgEncrypt)
	if msgSignatureGen != msgSignatureIn {
		return false
	}
	return true
}

func ParseEncryptRequestBody(r *http.Request) *model.EncryptRequestBody {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	//  mlog.AppendObj(nil, "Wechat Message Service: RequestBody--", body)
	requestBody := &model.EncryptRequestBody{}
	xml.Unmarshal(body, requestBody)
	return requestBody
}

func ParseTextRequestBody(r *http.Request) *model.TextRequestBody {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Error("%v", err)
		return nil
	}
	requestBody := &model.TextRequestBody{}
	xml.Unmarshal(body, requestBody)
	return requestBody
}

func MakeSignature(timestamp, nonce string) string {
	sl := []string{wxConfig.Token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func MakeMsgSignature(timestamp, nonce, msg_encrypt string) string {
	sl := []string{wxConfig.Token, timestamp, nonce, msg_encrypt}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func Value2CDATA(v string) model.CDATAText {
	//return CDATAText{[]byte("<![CDATA[" + v + "]]>")}
	return model.CDATAText{"<![CDATA[" + v + "]]>"}
}

func MakeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &model.TextResponseBody{}
	textResponseBody.FromUserName = Value2CDATA(fromUserName)
	textResponseBody.ToUserName = Value2CDATA(toUserName)
	textResponseBody.MsgType = Value2CDATA("text")
	textResponseBody.Content = Value2CDATA(content)
	textResponseBody.CreateTime = strconv.Itoa(int(time.Duration(time.Now().Unix())))
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}

func MakeEncryptResponseBody(fromUserName, toUserName, content, nonce, timestamp string) ([]byte, error) {
	encryptBody := &model.EncryptResponseBody{}
	encryptXmlData, _ := MakeEncryptXmlData(fromUserName, toUserName, timestamp, content)
	encryptBody.Encrypt = Value2CDATA(encryptXmlData)
	encryptBody.MsgSignature = Value2CDATA(MakeMsgSignature(timestamp, nonce, encryptXmlData))
	encryptBody.TimeStamp = timestamp
	encryptBody.Nonce = Value2CDATA(nonce)
	return xml.MarshalIndent(encryptBody, " ", "  ")
}

func MakeEncryptXmlData(fromUserName, toUserName, timestamp, content string) (string, error) {
	textResponseBody := &model.TextResponseBody{}
	textResponseBody.FromUserName = Value2CDATA(fromUserName)
	textResponseBody.ToUserName = Value2CDATA(toUserName)
	textResponseBody.MsgType = Value2CDATA("text")
	textResponseBody.Content = Value2CDATA(content)
	textResponseBody.CreateTime = timestamp
	body, err := xml.MarshalIndent(textResponseBody, " ", "  ")
	if err != nil {
		return "", errors.New("xml marshal error")
	}
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, int32(len(body)))
	if err != nil {
		log.Error("Binary write err:", err)
	}
	bodyLength := buf.Bytes()
	randomBytes := []byte("abcdefghijklmnop")
	plainData := bytes.Join([][]byte{randomBytes, bodyLength, body, []byte(wxConfig.AppID)}, nil)
	cipherData, err := AesEncrypt(plainData, wxConfig.AesKey)
	if err != nil {
		return "", errors.New("AesEncrypt error")
	}
	return base64.StdEncoding.EncodeToString(cipherData), nil
}

// PadLength calculates padding length, from github.com/vgorin/cryptogo
func PadLength(slice_length, blocksize int) (padlen int) {
	padlen = blocksize - slice_length%blocksize
	if padlen == 0 {
		padlen = blocksize
	}
	return padlen
}

//from github.com/vgorin/cryptogo
func PKCS7Pad(message []byte, blocksize int) (padded []byte) {
	// block size must be bigger or equal 2
	if blocksize < 1<<1 {
		panic("block size is too small (minimum is 2 bytes)")
	}
	// block size up to 255 requires 1 byte padding
	if blocksize < 1<<8 {
		// calculate padding length
		padlen := PadLength(len(message), blocksize)


		// define PKCS7 padding block
		padding := bytes.Repeat([]byte{byte(padlen)}, padlen)


		// apply padding
		padded = append(message, padding...)
		return padded
	}
	// block size bigger or equal 256 is not currently supported
	panic("unsupported block size")
}

func AesEncrypt(plainData []byte, aesKey []byte) ([]byte, error) {
	k := len(aesKey)
	if len(plainData)%k != 0 {
		plainData = PKCS7Pad(plainData, k)
	}


	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}


	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}


	cipherData := make([]byte, len(plainData))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherData, plainData)


	return cipherData, nil
}

func AesDecrypt(content string) ([]byte, error) {
	aesKey := wxConfig.AesKey
	// Decode base64
	cipherData, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}
	// AES Decrypt


	k := len(aesKey) //PKCS#7
	if len(cipherData)%k != 0 {
		return nil, errors.New("crypto/cipher: ciphertext size is not multiple of aes key length")
	}


	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}


	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainData := make([]byte, len(cipherData))
	blockMode.CryptBlocks(plainData, cipherData)
	return plainData, nil
}

func ValidateAppId(id []byte) bool {
	if string(id) == wxConfig.AppID {
		return true
	}
	return false
}


func ParseEncryptTextRequestBody(plainText []byte) *model.TextRequestBody {
	// Read length
	buf := bytes.NewBuffer(plainText[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)


	// appID validation
	appIDstart := 20 + length
	id := plainText[appIDstart : int(appIDstart)+len(wxConfig.AppID)]
	if !ValidateAppId(id) {
		log.Debug("parse encrypt text body error : Appid is invalid")
		//mlog.AppendObj(nil, "Wechat Message Service: appid is invalid!")
		return nil
	}
	//mlog.AppendObj(nil, "Wechat Message Service: appid validation is ok!")
	textRequestBody := &model.TextRequestBody{}
	xml.Unmarshal(plainText[20:20+length], textRequestBody)
	return textRequestBody
}