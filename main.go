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

	fmt.Println("Resolving URL...")
	resolvedURL, err := utils.ResolveURL(webURL)
	if err != nil {
		fmt.Println("Error resolving URL:", err)
		return
	}
	fmt.Println("Resolved URL:", resolvedURL)

	fmt.Println("Loading HTTP clients configuration...")
	clients, err := utils.LoadClients(cfgFile)
	if err != nil {
		fmt.Println("Error loading HTTP clients:", err)
		return
	}
	fmt.Println("Loaded HTTP clients configuration.")

	fmt.Println("Fetching HTML content...")
	htmlContent, err := utils.FetchHTML(resolvedURL)
	if err != nil {
		fmt.Println("Error fetching HTML content:", err)
		return
	}
	fmt.Println("Fetched HTML content.")

	fmt.Println("Parsing script tags...")
	scriptURLs := utils.ParseScripts(htmlContent, resolvedURL)
	if len(scriptURLs) == 0 {
		fmt.Println("No script URLs found.")
		return
	}
	fmt.Println("Found script URLs:", scriptURLs)

	fmt.Println("Extracting URLs and methods from JavaScript files...")
	allExtracted := utils.ExtractAllURLs(scriptURLs, clients, resolvedURL)
	if len(allExtracted) == 0 {
		fmt.Println("No URLs and methods extracted.")
		return
	}

	fmt.Println("Extracted URLs and Methods:")
	for u, methods := range allExtracted {
		fmt.Printf("URL: %s\nMethods: ", u)
		for m := range methods {
			fmt.Printf("%s ", m)
		}
		fmt.Println()
	}
}
