package ip

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Get() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get IP: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	ip, ok := result["ip"]
	if !ok {
		return "", fmt.Errorf("ip address not found in response")
	}

	return ip, nil
}
