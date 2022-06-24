package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hosseintrz/torob/auth/conf"
	pb "github.com/hosseintrz/torob/auth/pb/user"
	"github.com/hosseintrz/torob/auth/pkg/errors"
	"strings"
)

func GenerateToken(header string, payload interface{}, secret string) (string, error) {

	h := hmac.New(sha256.New, []byte(secret))
	encoder := base64.StdEncoding
	header64 := encoder.EncodeToString([]byte(header))
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("couldn't marshal paylaod!")
		return string(payloadStr), err
	}
	payload64 := encoder.EncodeToString(payloadStr)
	//fmt.Printf("payload64: %s\n", payload64)
	//fmt.Println("payload64 2 : ", payload64)
	token := header64 + "." + payload64
	unSignedStr := header + string(payloadStr)

	h.Write([]byte(unSignedStr))
	signature := encoder.EncodeToString(h.Sum(nil))

	token += "." + signature
	return token, nil
}

func ValidateToken(token string, secret string) (*pb.UserMsg, error) {
	subTokens := strings.Split(token, ".")
	h := hmac.New(sha256.New, []byte(secret))
	encoding := base64.StdEncoding
	if len(subTokens) != 3 {
		return nil, errors.ErrInvalidToken
	}
	header, err := encoding.DecodeString(subTokens[0])
	if err != nil {
		return nil, errors.ErrInvalidToken
	}
	payload, err := encoding.DecodeString(subTokens[1])
	if err != nil {
		return nil, errors.ErrInvalidToken
	}
	var unSignedStr string = string(header) + string(payload)
	h.Write([]byte(unSignedStr))
	newSignature := encoding.EncodeToString(h.Sum(nil))

	if newSignature == subTokens[2] {
		var userPayload pb.UserMsg
		err = json.Unmarshal(payload, &userPayload)
		fmt.Println("err: ", err)
		fmt.Println(userPayload.Username)
		return &userPayload, err
	}
	return nil, errors.ErrInvalidToken
}

func GetSignedToken(payload *pb.UserMsg) (string, error) {
	header := "HS256"
	secret, err := conf.GetEnv("JWT_SECRET")
	if err != nil {
		fmt.Println("error getting secret")
		return "", err
	}
	token, err := GenerateToken(header, payload, secret)
	return token, err
}
