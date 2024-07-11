package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Proxy struct {
	IP   string
	Port string
}

// GetResponseWithProxy get response with proxy (in some cases we need to use proxy to get data from OpenWeatherMap API)
func GetResponseWithProxy(requestURL string, proxy Proxy) ([]byte, error) {
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%s", proxy.IP, proxy.Port))
	if err != nil {
		log.Println("Wrong proxy:", proxy, "Error", err)
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
