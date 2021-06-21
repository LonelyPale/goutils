package aes

import (
	"encoding/base64"
	"github.com/LonelyPale/goutils/crypto"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

const sessionKey = "tiihtNczf5v6AKRyjwEUhQ=="
const iv = "r7BXXKkLb8qrSNn05n0qiA=="
const encryptedData = "CiyLU1Aw2KjvrjMdj8YKliAjtP4gsMZMQmRzooG2xrDcvSnxIMXFufNstNGTyaGS9uT5geRa0W4oTOb1WT7fJlAC+oNPdbB+3hVbJSRgv+4lGOETKUQz6OYStslQ142dNCuabNPGBzlooOmB231qMM85d2/fV6ChevvXvQP8Hkue1poOFtnEtpyxVLW1zAo6/1Xx1COxFvrc2d7UL/lmHInNlxuacJXwu0fjpXfz/YqYzBIBzD6WUfTIF9GRHpOn/Hz7saL8xz+W//FRAUid1OksQaQx4CMs8LOddcQhULW4ucetDf96JcR3g0gfRK4PC7E/r7Z6xNrXd2UIeorGj5Ef7b1pJAYB6Y5anaHqZ9J6nKEBvB4DnNLIVWSgARns/8wR2SiRS7MNACwTyrGvt9ts8p12PKFdlqYTopNHR1Vf7XjfhQlVsAJdNiKdYmYVoKlaRv85IfVunYzO0IKXsyl7JCUjCpoG20f0a04COwfneQAGGwd5oa+T8yO5hzuyDb/XcxxmK01EpqOyuxINew=="
const result = `{"openId":"oGZUI0egBJY1zhBYw2KhdUfwVJJE","nickName":"Band","gender":1,"language":"zh_CN","city":"Guangzhou","province":"Guangdong","country":"CN","avatarUrl":"http://wx.qlogo.cn/mmopen/vi_32/aSKcBBPpibyKNicHNTMM0qJVh8Kjgiak2AHWr8MHM4WgMEm7GFhsf8OYrySdbvAMvTsw3mo8ibKicsnfN5pRjl1p8HQ/0","unionId":"ocMvos6NjeKLIBqg5Mr9QjxrP1FA","watermark":{"timestamp":1477314187,"appid":"wx4f4bc4dec97d474b"}}`

func Test(t *testing.T) {
	text := []byte("hello world! 你好，中国！")
	key := []byte("1234567890123456")

	bs, err := Encrypt(text, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bs)

	ts, err := Decrypt(bs, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(ts))
}

func TestAESCBCDecrypt(t *testing.T) {
	sessionKeyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		log.Fatalln(err)
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		log.Fatalln(err)
	}

	encryptedDataBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		log.Fatalln(err)
	}

	plainText, err := CBCDecrypt(encryptedDataBytes, sessionKeyBytes, ivBytes)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("解密后为：", string(plainText))
	if string(plainText) == result {
		log.Println("ok")
	} else {
		log.Fatal("fail")
	}
}

func TestAESCBCEncrypt(t *testing.T) {
	sessionKeyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("sessionKey:", sessionKey, string(sessionKeyBytes))

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("iv:", iv, string(ivBytes))

	encryptedDataBytes, err := CBCEncrypt([]byte(result), sessionKeyBytes, ivBytes)
	if err != nil {
		log.Fatalln(err)
	}

	encrypted := base64.StdEncoding.EncodeToString(encryptedDataBytes)
	if encrypted == encryptedData {
		log.Println("ok")
	} else {
		log.Fatal("fail")
	}

	log.Println("encryptedData:", encryptedDataBytes)
	log.Println("encryptedData+base64:", encrypted)
}

func TestBigFile(t *testing.T) {
	t1 := time.Now()
	bs, err := ioutil.ReadFile("/Users/wyb/backup/software/os/ubuntu-20.04.2-live-server-amd64.iso")
	if err != nil {
		t.Fatal(err)
	}
	elapsed := time.Since(t1)
	t.Log(len(bs), elapsed)

	key, err := crypto.GenerateSecretKey()
	if err != nil {
		t.Fatal(err)
	}

	t1 = time.Now()
	bs2, err := Encrypt(bs, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(bs2), time.Since(t1))

	t1 = time.Now()
	ts, err := Decrypt(bs2, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(ts), time.Since(t1))
}
