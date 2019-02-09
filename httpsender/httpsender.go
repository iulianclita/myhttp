package httpsender

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Sender defines the contract to send HTTP requests
type Sender interface {
	Do(req *http.Request) (*http.Response, error)
}

// Make sends a request using the provided sender client
func Make(s Sender, req *http.Request) (string, error) {
	res, err := s.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(body)), nil
}
