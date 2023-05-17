package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	CEP_1_URL = "https://cdn.apicep.com/file/apicep/%s.json"
	CEP_2_URL = "http://viacep.com.br/ws/%s/json/"
)

type cep struct {
	HttpClient HttpClientInterface
	resp1      chan string
	resp2      chan string
}

func NewCepService() *cep {
	return &cep{
		HttpClient: http.DefaultClient,
		resp1:      make(chan string),
		resp2:      make(chan string),
	}
}

func (c *cep) GetCep(dto *InputDto) (*OutputDto, error) {
	go c.makeRequest(dto.Cep, CEP_1_URL, c.resp1)
	go c.makeRequest(dto.Cep, CEP_2_URL, c.resp2)

	var body string
	var api string
	select {
	case body = <-c.resp1:
		api = "https://cdn.apicep.com"
	case body = <-c.resp2:
		api = "http://viacep.com.br"
	case <-time.After(time.Second):
		body = "timeout de 1 segundo atingido"
	}

	return &OutputDto{
		Api:      api,
		Response: body,
	}, nil
}

func (c *cep) makeRequest(cep, url string, ch chan string) {
	resp, err := c.HttpClient.Get(fmt.Sprintf(url, cep))
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ch <- string(body)
}
