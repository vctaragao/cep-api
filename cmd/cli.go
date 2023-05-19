package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vctaragao/cep-api/internal"
)

func main() {
	if len(os.Args[1:]) == 0 {
		log.Fatalln("Error: um cep deve ser passado")
	}
	cep := os.Args[1]

	if !strings.Contains("-", cep) {
		cep = fmt.Sprintf("%s-%s", cep[:len(cep)-3], cep[len(cep)-3:])
	}

	out, err := internal.GetCep(cep)
	if err != nil {
		log.Fatalf("buscando cep: %s", err)
	}

	fmt.Printf("Resposta da API %s\n", out.Api)
	fmt.Println(out.Response)
}
