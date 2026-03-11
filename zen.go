package main

import (
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// zenMsg carries the fetched zen quote (or error).
type zenMsg struct {
	quote string
	err   error
}

func fetchZen() tea.Msg {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://api.github.com/zen", nil)
	if err != nil {
		return zenMsg{err: err}
	}

	// try to use gh auth token if available
	if token, tokenErr := ghAuthToken(); tokenErr == nil && token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return zenMsg{err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return zenMsg{err: err}
	}
	quote := strings.TrimSpace(string(body))
	quote = strings.Trim(quote, "\"")
	return zenMsg{quote: quote}
}

func ghAuthToken() (string, error) {
	out, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
