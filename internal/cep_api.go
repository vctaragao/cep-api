package internal

import "github.com/vctaragao/cep-api/internal/service"

func GetCep(cep string) (*service.OutputDto, error) {
	dto := &service.InputDto{Cep: cep}
	return service.NewCepService().GetCep(dto)
}
