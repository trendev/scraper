package utils

import (
	"encoding/json"
	"io/ioutil"
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

func LoadClients(cfgFile string) ([]HTTPClient, error) {
	var clients []HTTPClient
	cfgData, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(cfgData, &clients)
	if err != nil {
		return nil, err
	}

	for i, client := range clients {
		clients[i].methodRegexp, err = regexp.Compile(client.MethodRe)
		if err != nil {
			return nil, err
		}
		clients[i].urlRegexp, err = regexp.Compile(client.URLRe)
		if err != nil {
			return nil, err
		}
	}

	return clients, nil
}
