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

	// Test with a malformed URL
	_, err = FetchHTML("http://%41:8080/")
	if err == nil {
		t.Fatalf("Expected error for malformed URL, got none")
	}

	// Test with a non-existent URL
	_, err = FetchHTML("http://localhost:9999")
	if err == nil {
		t.Fatalf("Expected error for non-existent URL, got none")
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
	expected := []string{"http://example.com/test1.js", "http://example.com/test2.js"}

	result := ParseScripts(htmlContent, "http://example.com")
	if len(result) != len(expected) {
		t.Fatalf("Expected %d URLs, got %d", len(expected), len(result))
	}
	for i, url := range result {
		if url != expected[i] {
			t.Fatalf("Expected %v, got %v", expected[i], url)
		}
	}

	// Test with no scripts
	htmlContentNoScripts := `<html><head></head></html>`
	resultNoScripts := ParseScripts(htmlContentNoScripts, "http://example.com")
	if len(resultNoScripts) != 0 {
		t.Fatalf("Expected 0 URLs, got %d", len(resultNoScripts))
	}

	// Test with malformed HTML
	htmlContentMalformed := `<html><head><script src="test1.js"></script><script src="test2.js"></head></html`
	resultMalformed := ParseScripts(htmlContentMalformed, "http://example.com")
	if len(resultMalformed) != len(expected) {
		t.Fatalf("Expected %d URLs, got %d", len(expected), len(resultMalformed))
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

	// Test with malformed JS content
	jsContentMalformed := `
        fetch("https://example.com/api", {method: "POST");
        axios.get("https://example.com/axios"
        axios.post("https://example.com/axios-post";
    `
	resultMalformed := ExtractURLs(jsContentMalformed, clients)
	if len(resultMalformed) != 0 {
		t.Fatalf("Expected 0 URLs, got %d", len(resultMalformed))
	}
}
