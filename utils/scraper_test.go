package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchHTML(t *testing.T) {
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Final destination"))
	})
	finalServer := httptest.NewServer(finalHandler)
	defer finalServer.Close()

	redirectHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, finalServer.URL, http.StatusMovedPermanently)
	})
	redirectServer := httptest.NewServer(redirectHandler)
	defer redirectServer.Close()

	url := redirectServer.URL
	expected := "Final destination"

	result, err := FetchHTML(url)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != expected {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}

func TestParseScripts(t *testing.T) {
	htmlContent := `
        <html>
            <head>
                <script src="test1.js"></script>
                <script src="test2.js"></script>
            </head>
        </html>`
	expected := []string{"test1.js", "test2.js"}

	result := ParseScripts(htmlContent, "http://example.com")
	if len(result) != len(expected) {
		t.Fatalf("Expected %d URLs, got %d", len(expected), len(result))
	}
	for i, url := range result {
		if url != expected[i] {
			t.Fatalf("Expected %v, got %v", expected[i], url)
		}
	}
}

func TestExtractURLsAndMethods(t *testing.T) {
	jsContent := `
        fetch("https://example.com/api", {method: "POST"});
        axios.get("https://example.com/axios");
        axios.post("https://example.com/axios-post");
    `
	clients := loadTestConfig()
	result := ExtractURLs(jsContent, clients)

	expected := map[string]map[string]bool{
		"https://example.com/api":        {"POST": true},
		"https://example.com/axios":      {"GET": true},
		"https://example.com/axios-post": {"POST": true},
	}

	for url, methods := range expected {
		if _, exists := result[url]; !exists {
			t.Fatalf("Expected URL %v to exist", url)
		}
		for method := range methods {
			if _, exists := result[url][method]; !exists {
				t.Fatalf("Expected method %v for URL %v to exist", method, url)
			}
		}
	}
}
