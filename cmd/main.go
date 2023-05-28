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
	chTm1 := make(chan  time.Duration)
	chTm2 := make(chan  time.Duration)

	reader := bufio.NewReader(os.Stdin)

	line, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}
	go requestByApiViaCep(strings.TrimSpace(line), ch1,chTm1)
	go requestByApiCep(strings.TrimSpace(line), ch2,chTm2)

	select {
	case cep := <-ch2:
		{
			fmt.Println("ApiCep foi mais rápido")
			fmt.Printf("ApiCep Time taken: %v\n", <-chTm2)
			fmt.Println(cep)
		}
	case cep := <-ch1:
		{
			fmt.Println("Via cep foi mais rápido")
			fmt.Printf("ApiCep Time taken: %v\n", <-chTm1)
			fmt.Println(cep)
		}
	}

}

func requestByApiViaCep(cep string, ch chan<- ViaCep, chT chan<- time.Duration) {

	now := time.Now()
	str := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, str, nil)

	if err != nil {
		log.Println(err)
	}

	data, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	bt, err := io.ReadAll(data.Body)

	if err != nil {
		log.Println(err)
	}

	var cepCep ViaCep

	if err = json.Unmarshal(bt, &cepCep); err != nil {
		log.Fatal(err)
	}

	diff := time.Since(now)

	ch <- cepCep
	chT <- diff
}

func requestByApiCep(cep string, ch chan<- Cep,chT chan<- time.Duration) {

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

	ch <- cepCep
	chT <- diff
}
