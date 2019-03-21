package dnsbutler

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type response struct {
	Err        error
	URL        string
	StatusCode int
	Body       []byte
}

func updateTarget(urls, ip string, ch chan<- *response) {
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
		response.Err = err
		ch <- response
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.Err = err
		ch <- response
		return
	}

	response.Body = body
	response.StatusCode = resp.StatusCode

	ch <- response
}

func updateTargets(targets []string, newIP string, logger *log.Logger) {
	ch := make(chan *response)

	for _, t := range targets {
		go updateTarget(t, newIP, ch)
	}

	for range targets {
		r := <-ch
		if r.Err != nil {
			log.Printf("Received err '%v' for url '%s'", r.Err, r.URL)
			return
		}

		if r.StatusCode != http.StatusOK {
			log.Printf("Received StatusCode '%d' for url '%s'", r.StatusCode, r.URL)
			return
		}

		log.Printf("IP for url '%s' updated", r.URL)
	}
}
