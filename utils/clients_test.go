package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
)

func loadTestConfig() []HTTPClient {
	config := `{
        "clients": [
            {
                "name": "fetch",
                "methodRe": "fetch\\(\\\"([^\\\"]+)\\\",\\s*\\{[^}]*method:\\s*\\\"([^\\\"]+)\\\"\\s*\\}",
                "urlRe": "fetch\\(\\\"([^\\\"]+)\\\"\\)",
                "urlMethodIndex": [1, 2]
            },
            {
                "name": "axios",
                "methodRe": "axios\\.([a-z]+)\\(\\\"([^\\\"]+)\\\"\\s*(,\\s*\\{[^}]*\\})?\\)",
                "urlRe": "axios\\.([a-z]+)\\(\\\"([^\\\"]+)\\\"\\s*(,\\s*\\{[^}]*\\})?\\)",
                "urlMethodIndex": [2, 1]
            }
        ]
    }`
	var configData struct {
		Clients []HTTPClient `json:"clients"`
	}
	err := json.Unmarshal([]byte(config), &configData)
	if err != nil {
		panic(err)
	}
	clients := configData.Clients
	for i, client := range clients {
		clients[i].methodRegexp = regexp.MustCompile(client.MethodRe)
		clients[i].urlRegexp = regexp.MustCompile(client.URLRe)
	}
	return clients
}

func TestLoadClients(t *testing.T) {
	clients := loadTestConfig()
	fmt.Printf("Loaded clients: %+v\n", clients) // Debug statement
	if len(clients) != 2 {
		t.Fatalf("Expected 2 clients, got %d", len(clients))
	}

	fetchClient := clients[0]
	if fetchClient.Name != "fetch" {
		t.Fatalf("Expected client name to be 'fetch', got %s", fetchClient.Name)
	}

	axiosClient := clients[1]
	if axiosClient.Name != "axios" {
		t.Fatalf("Expected client name to be 'axios', got %s", axiosClient.Name)
	}
}
