# Scraper Project

This project contains a Go-based scraper to extract URLs and HTTP methods from JavaScript files on a website.

## Prerequisites

- Go (version 1.16 or later)

## Setup

1. **Install Go**

    Follow the instructions on the [official Go website](https://golang.org/dl/) to install Go on your system.

2. **Build and Run the Scraper**

    To build and run the program, use the following commands:

    ```sh
    go build -o scraper main.go
    ./scraper -url https://poln.org
    ```

## Usage

To use the scraper, pass the URL of the website you want to scrape as a command-line argument using the `-url` flag. For example:

```sh
./scraper -url https://poln.org
```

### Command-Line Options

- **-url**: The main URL of the website to analyze. (Required)
- **-config**: The configuration file for HTTP clients. (Optional, default is `config.json`)

## Project Structure

`main.go`: The main entry point of the application.

`utils/`: Directory containing utility functions for fetching HTML, parsing scripts, and extracting URLs and methods.

## Configuration

The configuration for the HTTP clients (`fetch` and `axios`) is defined in a JSON file specified by the `-config` flag. Adjust the regex patterns as necessary to match the JavaScript syntax used on the target website.
