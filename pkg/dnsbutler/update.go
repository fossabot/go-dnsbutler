package dnsbutler

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type response struct {
	URL        string
	StatusCode int
	Body       []byte
}

func updateTarget(urls, ip string) (*response, error) {
	response := &response{
		URL: "",
	}

	u, _ := url.Parse(urls)
	if u != nil {
		q := u.Query()
		q.Del("password")
		q.Del("secret")
		u.RawQuery = q.Encode()
		u.User = nil
		response.URL = u.String()
	}

	resp, err := http.Get(urls)
	if err != nil {
		return response, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	response.Body = body
	response.StatusCode = resp.StatusCode

	return response, nil
}

func updateTargets(targets []string, newIP string, logger *log.Logger, done chan<- bool) {
	var waitGroup sync.WaitGroup
	for _, t := range targets {
		waitGroup.Add(1)

		go func(t string) {
			defer waitGroup.Done()

			r, err := updateTarget(t, newIP)
			if err != nil {
				log.Printf("Received err '%v' for url '%s'", err, r.URL)
				return
			}
			if r.StatusCode != http.StatusOK {
				log.Printf("Received StatusCode '%d' for url '%s'", r.StatusCode, r.URL)
				return
			}

			log.Printf("IP for url '%s' updated", r.URL)
		}(t)
	}

	waitGroup.Wait()

	done <- true
}
