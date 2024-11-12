package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestCommandPingLocalhost(t *testing.T) {
	command := struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
	}{
		Type:    "ping",
		Payload: "localhost",
	}

	jsonData, err := json.Marshal(command)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	resp, err := http.Post(
		"http://localhost:8080/execute",
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
			Successful bool `json:"successful"`
			Time       int  `json:"time"`
		} `json:"data"`
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.StatusCode)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if !response.Success || !response.Data.Successful {
		t.Errorf("Expected successful ping, got response: %s", string(body))
	}

	if response.Data.Time <= 0 {
		t.Errorf("Expected positive ping time, got: %d", response.Data.Time)
	}
}
