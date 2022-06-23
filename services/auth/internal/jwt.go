package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hosseintrz/torob/auth/conf"
	"strings"
	"time"
)

func GenerateToken(header string, payload map[string]string, secret string) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))
	encoder := base64.StdEncoding
	header64 := encoder.EncodeToString([]byte(header))
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("couldn't marshal paylaod!")
		return string(payloadStr), err
	}
	payload64 := encoder.EncodeToString(payloadStr)

	token := header64 + "." + payload64
	unSignedStr := header + string(payloadStr)

	h.Write([]byte(unSignedStr))
	signature := encoder.EncodeToString(h.Sum(nil))

	token += "." + signature
	return token, nil
}

func ValidateToken(token string, secret string) (bool, error) {
	subTokens := strings.Split(token, ".")
	h := hmac.New(sha256.New, []byte(secret))
	encoding := base64.StdEncoding
	if len(subTokens) != 3 {
		return false, nil
	}
	header, err := encoding.DecodeString(subTokens[0])
	if err != nil {
		return false, err
	}
	payload, err := encoding.DecodeString(subTokens[1])
	if err != nil {
		return false, err
	}
	var unSignedStr string = string(header) + string(payload)
	h.Write([]byte(unSignedStr))
	newSignature := encoding.EncodeToString(h.Sum(nil))

	return newSignature == subTokens[2], nil
}

func GetSignedToken() (string, error) {
	claims := map[string]string{
		"aud": "frontend.user",
		"iss": "frontend",
		"exp": fmt.Sprint(time.Now().Add(1 * time.Minute).Unix()),
	}
	header := "HS256"
	secret, err := conf.GetEnv("JWT_SECRET")
	if err != nil {
		fmt.Println("error getting secret")
		return "", err
	}
	token, err := GenerateToken(header, claims, secret)
	return token, err
}
