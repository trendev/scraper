package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func FetchHTML(u string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	resp, err := client.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func ParseScripts(htmlContent, baseURL string) []string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil
	}

	var scriptURLs []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					scriptURL := attr.Val
					if !strings.HasPrefix(scriptURL, "http") {
						base, err := url.Parse(baseURL)
						if err != nil {
							continue
						}
						scriptURL = base.ResolveReference(&url.URL{Path: scriptURL}).String()
					}
					scriptURLs = append(scriptURLs, scriptURL)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return scriptURLs
}

func FetchJS(u string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	resp, err := client.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func ExtractURLs(jsContent string, clients []HTTPClient) map[string]map[string]bool {
	results := make(map[string]map[string]bool)

	for _, client := range clients {
		methodMatches := client.methodRegexp.FindAllStringSubmatch(jsContent, -1)
		fmt.Printf("Method Matches for client %s: %+v\n", client.Name, methodMatches) // Debug statement

		if len(methodMatches) == 0 {
			fmt.Printf("No method matches for client %s\n", client.Name)
			continue
		}

		for _, match := range methodMatches {
			if len(client.URLMethodIndex) < 2 {
				fmt.Printf("Skipping client %s due to insufficient URLMethodIndex\n", client.Name)
				continue
			}
			if len(match) <= client.URLMethodIndex[0] || len(match) <= client.URLMethodIndex[1] {
				fmt.Printf("Skipping match due to insufficient length: %+v\n", match)
				continue
			}
			url := match[client.URLMethodIndex[0]]
			method := match[client.URLMethodIndex[1]]
			if _, exists := results[url]; !exists {
				results[url] = make(map[string]bool)
			}
			results[url][strings.ToUpper(method)] = true
		}
	}

	return results
}

func ExtractAllURLs(scriptURLs []string, clients []HTTPClient, baseURL string) map[string]map[string]bool {
	allExtracted := make(map[string]map[string]bool)
	for _, scriptURL := range scriptURLs {
		if !strings.HasPrefix(scriptURL, "http") {
			scriptURL = baseURL + "/" + scriptURL
		}
		jsContent, err := FetchJS(scriptURL)
		if err != nil {
			continue
		}
		extracted := ExtractURLs(jsContent, clients)
		for url, methods := range extracted {
			if _, exists := allExtracted[url]; !exists {
				allExtracted[url] = make(map[string]bool)
			}
			for method := range methods {
				allExtracted[url][method] = true
			}
		}
	}
	return allExtracted
}

func ResolveURL(u string) (string, error) {
	resolvedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			resolvedURL = req.URL
			return nil
		},
	}
	_, err = client.Get(u)
	if err != nil {
		return "", err
	}
	return resolvedURL.String(), nil
}
