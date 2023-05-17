package service

import "net/http"

type RequestCepInterface interface {
	execute(cep, url string, ch chan string)
}

type HttpClientInterface interface {
	Get(url string) (*http.Response, error)
}
