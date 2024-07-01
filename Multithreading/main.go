package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BrasilCep struct {
	Cep          string
	State        string
	City         string
	Neighborhood string
	Street       string
	Service      string
}

type ViaCep struct {
	Cep         string
	Logradouro  string
	Complemento string
	Unidade     string
	Bairro      string
	Localidade  string
	Uf          string
	Ibge        string
	Gia         string
	Ddd         string
	Siafi       string
}

type Request struct {
	Tipo string
}

func main() {
	ch1 := make(chan BrasilCep)
	ch2 := make(chan ViaCep)

	go ConsultaBrasilCep("88115710", ch1)
	go ConsultaViaCep("88115710", ch2)

	select {
	case retornoBrasilCep := <-ch1:
		fmt.Println(retornoBrasilCep, "Brasil CEP")
	case retornoViaCep := <-ch2:
		fmt.Println(retornoViaCep, "Via CEP")
	case <-time.After(time.Second * 1):
		fmt.Println("Timeout. Try again")

	}

}

func ConsultaBrasilCep(reqcep string, ch chan BrasilCep) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://brasilapi.com.br/api/cep/v1/"+reqcep, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var cep BrasilCep
	err = json.Unmarshal(body, &cep)

	ch <- cep
}

func ConsultaViaCep(reqcep string, ch chan ViaCep) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/"+reqcep+"/json/", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var cep ViaCep
	err = json.Unmarshal(body, &cep)

	ch <- cep

}
