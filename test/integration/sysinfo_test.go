package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestCommandSysinfo(t *testing.T) {
	command := struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
	}{
		Type: "sysinfo",
	}

	jsonData, err := json.Marshal(command)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/execute", testPort),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Hostname  string `json:"hostname"`
			IPAddress string `json:"ip_address"`
		} `json:"data"`
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.StatusCode)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected successful request, got response: %s", string(body))
	}

	if len(response.Data.Hostname) <= 0 {
		t.Errorf("Expected nonempty hostname")
	}

	if len(response.Data.IPAddress) <= 0 {
		t.Errorf("Expected nonempty hostname")
	}
}
