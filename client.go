package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"http://localhost:8080/cotacao",
		nil,
	)
	if err != nil {
		log.Println("Erro criando request:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Timeout cliente ou erro na requisição:", err)
		return
	}
	defer resp.Body.Close()

	var cotacao Cotacao
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		log.Println("Erro decode resposta:", err)
		return
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println("Erro criando arquivo:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("Dólar: " + cotacao.Bid)
	if err != nil {
		log.Println("Erro escrevendo arquivo:", err)
		return
	}

	log.Println("Cotação salva com sucesso!")
}
