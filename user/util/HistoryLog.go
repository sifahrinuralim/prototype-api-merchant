package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Payload struct {
	Action      string `json:"action"`
	Description string `json:"description"`
}

func HistoryLog(action, description string) error {
	url := "http://localhost:8200/v1/history/log"

	payload := Payload{
		Action:      action,
		Description: description,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	fmt.Println("Request berhasil dengan status:", resp.Status)
	return nil
}
