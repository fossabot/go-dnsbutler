package dnsbutler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var ipRegex = regexp.MustCompile(`(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)

func getIPFrom(url string, logger *log.Logger) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := string(ipRegex.Find(body))
	if ip == "" {
		return "", fmt.Errorf("Regex did not match an IP address in body")
	}

	logger.Printf("Received IP '%s' from provider '%s'\n", ip, url)

	return ip, nil
}
