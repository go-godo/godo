package godo

import (
	"io/ioutil"
	"net/http"
)

func HTTPGetString(url string) (string, error) {
	resp, err := http.Get(url)
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

func MustHTTPGetString(url string) string {
	body, err := HTTPGetString(url)
	if err != nil {
		panic(&mustPanic{
			err: err,
		})
	}
	return body
}
