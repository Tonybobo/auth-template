package utils

import (
	"encoding/base64"
	"fmt"
)

func Encode(s string) string {
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return string(data)
}

func Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("decoding error , %w", err)
	}

	return string(data), nil
}
