package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	value, ok := h[key]
	return value, ok
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	indexCRLF := bytes.Index(data, []byte(crlf))
	if indexCRLF == -1 {
		return 0, false, nil
	}

	if indexCRLF == 0 {
		return len(crlf), true, nil
	}
	fieldLine := data[:indexCRLF]

	indexColon := bytes.Index(fieldLine, []byte(":"))
	fieldKey := strings.TrimLeft(string(data[:indexColon]), " ")
	fieldKey = strings.ToLower(fieldKey)

	if isKeyValid := validateFieldKey(fieldKey); !isKeyValid {
		return 0, false, fmt.Errorf("error: field key has invalid characters: %s", fieldKey)
	}

	fieldValue := strings.TrimSpace(string(fieldLine[indexColon+1:]))

	value, ok := h[fieldKey]
	if ok {
		h[fieldKey] = strings.Join([]string{value, fieldValue}, ", ")
	} else {
		h[fieldKey] = fieldValue
	}

	return indexCRLF + len(crlf), false, nil
}

func (h Headers) Override(key, value string) error {
	key = strings.ToLower(key)
	if _, ok := h[key]; !ok {
		return fmt.Errorf("key does not exist in headers: %s", key)
	}
	h[key] = value
	return nil
}

func validateFieldKey(key string) bool {
	validChars := map[string]bool{
		"!": true, "#": true, "$": true, "%": true, "&": true,
		"'": true, "*": true, "+": true, "-": true, ".": true,
		"^": true, "_": true, "`": true, "|": true, "~": true,
		"a": true, "b": true, "c": true, "d": true, "e": true,
		"f": true, "g": true, "h": true, "i": true, "j": true,
		"k": true, "l": true, "m": true, "n": true, "o": true,
		"p": true, "q": true, "r": true, "s": true, "t": true,
		"u": true, "v": true, "w": true, "x": true, "y": true, "z": true,
		"0": true, "1": true, "2": true, "3": true, "4": true,
		"5": true, "6": true, "7": true, "8": true, "9": true,
	}

	for _, char := range key {
		if _, ok := validChars[string(char)]; !ok {
			return false
		}
	}

	return true
}
