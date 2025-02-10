package bytesutil

import "encoding/json"

// JSONEscape escapes the given string such that it is a valid JSON string object.
func JSONEscape(s string) (escapedString string, err error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	// Trim the beginning and trailing " character
	return string(b[1 : len(b)-1]), err
}
