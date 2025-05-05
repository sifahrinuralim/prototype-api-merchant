package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"history/config"
	"log"
	"strings"
)

func ValidateSignature(rawBody []byte, xTimestamp, signature string) bool {
	var buf bytes.Buffer
	if err := json.Compact(&buf, rawBody); err != nil {
		log.Printf("failed to minify JSON: %v", err)
		return false
	}
	minifiedJSON := buf.String()

	message := minifiedJSON + xTimestamp
	log.Printf("raw message after minify: %s", message)

	key := []byte(config.SignatureKey)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(message))
	expectedMAC := hex.EncodeToString(h.Sum(nil))

	finalAuth := strings.EqualFold(expectedMAC, signature)
	log.Printf("expectedMAC: %s", expectedMAC)
	log.Printf("finalAuth: %v", finalAuth)

	return finalAuth
}
