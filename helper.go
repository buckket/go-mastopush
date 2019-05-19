package mastopush

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func ParseHeader(header *http.Header) (dh, salt, token []byte, err error) {
	if header.Get("Content-Encoding") != "aesgcm" {
		return nil, nil, nil, fmt.Errorf("unsupported Content-Encoding: %s", header.Get("Content-Encoding"))
	}

	dh, err = encodedValue(header, "Crypto-Key", "dh")
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = encodedValue(header, "Encryption", "salt")
	if err != nil {
		return nil, nil, nil, err
	}

	t := header.Get("Authorization")
	if len(t) > 0 {
		token = []byte(strings.TrimPrefix(t, "WebPush "))
	}
	return dh, salt, token, nil
}

func encodedValue(header *http.Header, name, key string) ([]byte, error) {
	keyValues := parseKeyValues(header.Get(name))
	value, exists := keyValues[key]
	if !exists {
		return nil, fmt.Errorf("value %s not found in header %s", key, name)
	}

	bytes, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func parseKeyValues(values string) map[string]string {
	f := func(c rune) bool {
		return c == ';'
	}

	entries := strings.FieldsFunc(values, f)

	m := make(map[string]string)
	for _, entry := range entries {
		parts := strings.Split(entry, "=")
		m[parts[0]] = parts[1]
	}

	return m
}
