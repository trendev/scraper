package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type HTTPClient struct {
	Name           string `json:"name"`
	MethodRe       string `json:"methodRe"`
	URLRe          string `json:"urlRe"`
	URLMethodIndex []int  `json:"urlMethodIndex"`
	methodRegexp   *regexp.Regexp
	urlRegexp      *regexp.Regexp
}

func LoadClients(configFile string) ([]HTTPClient, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var config struct {
		Clients []HTTPClient `json:"clients"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal config data: %w", err)
	}

	for i, client := range config.Clients {
		config.Clients[i].methodRegexp = regexp.MustCompile(client.MethodRe)
		config.Clients[i].urlRegexp = regexp.MustCompile(client.URLRe)
	}

	return config.Clients, nil
}
