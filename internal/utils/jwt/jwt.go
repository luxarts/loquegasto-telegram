package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"loquegasto-telegram/internal/defines"
	"loquegasto-telegram/internal/utils/base64"
	"os"
	"strings"
)

type Header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type Payload struct {
	Subject int64 `json:"sub"`
}

func GenerateToken(header *Header, payload *Payload) string {
	if header == nil {
		header = &Header{
			Algorithm: "HS256",
			Type:      "JWT",
		}
	}
	headerBytes, _ := json.Marshal(header)
	headerEncoded := base64.Encode(headerBytes)

	payloadBytes, _ := json.Marshal(payload)
	payloadEncoded := base64.Encode(payloadBytes)

	signature := generateSignature(headerEncoded, payloadEncoded)

	return headerEncoded + "." + payloadEncoded + "." + signature
}

func generateSignature(header string, payload string) string {
	secret := base64.DecodeString(os.Getenv(defines.EnvJWTSecret))

	h := hmac.New(sha256.New, []byte(secret))

	h.Write([]byte(header + "." + payload))

	signature := base64.Encode(h.Sum(nil))
	return signature
}

func Verify(header string, payload string, signature string) bool {
	signatureGenerated := generateSignature(header, payload)

	return signatureGenerated == signature
}

func GetSubject(token string) (int64, error) {
	payload := strings.Split(token, ".")[1]
	payloadDecoded := base64.Decode(payload)

	var p Payload
	err := json.Unmarshal(payloadDecoded, &p)
	if err != nil {
		return 0, err
	}

	return p.Subject, nil
}
