package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Get(url string) (*http.Response, error) {
	m.Called(url)
	return getDummyResponse(), nil
}

func TestGetCepTestFirstChanResponse(t *testing.T) {
	chan1 := make(chan string)
	chan2 := make(chan string)

	c := "37501049"

	httpClient := &MockHttpClient{}
	httpClient.On("Get", fmt.Sprintf(CEP_1_URL, c)).Return(getDummyResponse(), nil).Once()
	httpClient.On("Get", fmt.Sprintf(CEP_2_URL, c)).Return(getDummyResponse(), nil).Once().After(time.Second)

	cep := cep{
		HttpClient: httpClient,
		resp1:      chan1,
		resp2:      chan2,
	}

	dto := &InputDto{Cep: c}
	out, err := cep.GetCep(dto)

	assert.NoError(t, err)
	assert.Equal(t, "resposta", out.Response)
	assert.Equal(t, "https://cdn.apicep.com", out.Api)
}

func TestGetCepTestSecondChanResponse(t *testing.T) {
	chan1 := make(chan string)
	chan2 := make(chan string)
	c := "37501049"

	httpClient := &MockHttpClient{}
	httpClient.On("Get", fmt.Sprintf(CEP_1_URL, c)).Return(getDummyResponse(), nil).Once().After(time.Second)
	httpClient.On("Get", fmt.Sprintf(CEP_2_URL, c)).Return(getDummyResponse(), nil).Once()

	cep := cep{
		HttpClient: httpClient,
		resp1:      chan1,
		resp2:      chan2,
	}

	dto := &InputDto{Cep: c}
	out, err := cep.GetCep(dto)

	assert.NoError(t, err)
	assert.Equal(t, "resposta", out.Response)
	assert.Equal(t, "http://viacep.com.br", out.Api)
}

func TestGetCepTestTimeoutChanResponse(t *testing.T) {
	chan1 := make(chan string)
	chan2 := make(chan string)

	c := "37501049"

	httpClient := &MockHttpClient{}
	httpClient.On("Get", fmt.Sprintf(CEP_1_URL, c)).Return(getDummyResponse(), nil).Once().After(time.Second * 2)
	httpClient.On("Get", fmt.Sprintf(CEP_2_URL, c)).Return(getDummyResponse(), nil).Once().After(time.Second * 2)

	cep := cep{
		HttpClient: httpClient,
		resp1:      chan1,
		resp2:      chan2,
	}

	dto := &InputDto{Cep: c}
	out, err := cep.GetCep(dto)

	assert.NoError(t, err)
	assert.Equal(t, "timeout de 1 segundo atingido", out.Response)
	assert.Equal(t, "", out.Api)
}

func getDummyResponse() *http.Response {
	w := httptest.NewRecorder()
	w.WriteString("resposta")
	return w.Result()
}
