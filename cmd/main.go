package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
	{
  "cep": "76530-000",
  "logradouro": "",
  "complemento": "",
  "bairro": "",
  "localidade": "Mundo Novo",
  "uf": "GO",
  "ibge": "5214051",
  "gia": "",
  "ddd": "62",
  "siafi": "9651"
}
*/
/*

{
  "code": "76530-000",
  "state": "GO",
  "city": "Mundo Novo",
  "district": "",
  "address": "",
  "status": 200,
  "ok": true,
  "statusText": "ok"
}
*/

type Cep struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	ch1 := make(chan ViaCep)
	ch2 := make(chan Cep)

	reader := bufio.NewReader(os.Stdin)

	line, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}
	go requestByApiViaCep(strings.TrimSpace(line), ch1)
	go requestByApiCep(strings.TrimSpace(line), ch2)

	select {
	case cep := <-ch1:
		{
			fmt.Println("Via cep foi mais rápido")
			fmt.Println(cep)

		}
	case cep := <-ch2:
		{
			fmt.Println("Cep foi mais rápido")
			fmt.Println(cep)
		}

	}

}

func requestByApiViaCep(cep string, ch chan<- ViaCep) {

	now := time.Now()
	str := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, str, nil)

	if err != nil {
		log.Fatal(err)
	}

	data, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	bt, err := io.ReadAll(data.Body)

	if err != nil {
		log.Fatal(err)
	}

	var cepCep ViaCep

	if err = json.Unmarshal(bt, &cepCep); err != nil {
		log.Fatal(err)
	}

	diff := time.Since(now)

	fmt.Printf("Via cep Time taken: %s\n", diff)
	ch <- cepCep
}

func requestByApiCep(cep string, ch chan<- Cep) {

	now := time.Now()
	str := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, str, nil)

	if err != nil {
		log.Fatal(err)
	}

	data, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	bt, err := io.ReadAll(data.Body)

	if err != nil {
		log.Fatal(err)
	}

	var cepCep Cep

	if err = json.Unmarshal(bt, &cepCep); err != nil {
		log.Fatal(err)
	}

	diff := time.Since(now)

	fmt.Printf("Api Cep Time taken: %s\n", diff)

	ch <- cepCep
}
