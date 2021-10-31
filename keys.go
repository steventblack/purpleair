package purpleair

import (
	"encoding/json"
	"net/http"
)

// Package variables for retaining the read & write keys.
// These are set by successful calls to SetAPIKey.
var (
	apiReadKey  string
	apiWriteKey string
)

// SetAPIKey calls CheckAPIKey to validate the key and, if valid, retains the
// read and write keys for further calls. Only one read key and write key will
// be retained with additional (valid) calls resulting in replacement.
func SetAPIKey(k string) (KeyType, error) {
	kt, err := CheckAPIKey(k)
	if err != nil {
		return KeyUnknown, err
	}

	switch kt {
	case KeyRead:
		apiReadKey = k
	case KeyWrite:
		apiWriteKey = k
	default:
	}

	return kt, nil
}

// CheckAPIKey checks the validity and permissions of the specified key.
// It does not save the key for further calls. Use SetAPIKey to retain key
// values if desired.
func CheckAPIKey(k string) (KeyType, error) {
	req, err := http.NewRequest(http.MethodGet, urlKeys, nil)
	if err != nil {
		return KeyUnknown, err
	}
	req.Header.Add(keyHeader, k)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return KeyUnknown, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return KeyUnknown, paError(resp)
	}

	payload := struct {
		K KeyType `json:"api_key_type"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		return KeyUnknown, err
	}

	return payload.K, nil
}
