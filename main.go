package main

import (
	"flag"
	"fmt"

	"github.com/trendev/scraper/utils"
)

func main() {
	var webURL, cfgFile string

	flag.StringVar(&webURL, "url", "", "The main URL of the website to analyze")
	flag.StringVar(&cfgFile, "config", "config.json", "The configuration file for HTTP clients")
	flag.Parse()

	if webURL == "" {
		fmt.Println("Please provide the main URL using the -url flag.")
		return
	}

	resolvedURL, err := utils.ResolveURL(webURL)
	if err != nil {
		fmt.Println("Error resolving URL:", err)
		return
	}

	clients, err := utils.LoadClients(cfgFile)
	if err != nil {
		fmt.Println("Error loading HTTP clients:", err)
		return
	}

	htmlContent, err := utils.FetchHTML(resolvedURL)
	if err != nil {
		fmt.Println("Error fetching HTML content:", err)
		return
	}

	scriptURLs := utils.ParseScripts(htmlContent, resolvedURL)
	if len(scriptURLs) == 0 {
		fmt.Println("No script URLs found.")
		return
	}

	allExtracted := utils.ExtractAllURLs(scriptURLs, clients, resolvedURL)
	fmt.Println("Extracted URLs and Methods:")
	for u, methods := range allExtracted {
		fmt.Printf("URL: %s\\nMethods: ", u)
		for m := range methods {
			fmt.Printf("%s ", m)
		}
		fmt.Println()
	}
}
