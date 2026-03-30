package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type CotacaoAPI struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	db, err := sql.Open("sqlite", "cotacoes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable(db)

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		handleCotacao(w, r, db)
	})

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		valor TEXT,
		data DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Erro criando tabela:", err)
	}
}

func handleCotacao(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	ctxAPI, cancelAPI := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancelAPI()

	req, err := http.NewRequestWithContext(
		ctxAPI,
		"GET",
		"https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil,
	)
	if err != nil {
		log.Println("Erro criando request API:", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao chamar API:", err)
		http.Error(w, "Timeout API externa", http.StatusGatewayTimeout)
		return
	}
	defer resp.Body.Close()

	var cotacao CotacaoAPI
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		log.Println("Erro decode API:", err)
		http.Error(w, "Erro processando API", http.StatusInternalServerError)
		return
	}

	valor := cotacao.USDBRL.Bid

	ctxDB, cancelDB := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancelDB()

	stmt, err := db.PrepareContext(ctxDB, "INSERT INTO cotacoes(valor) VALUES(?)")
	if err != nil {
		log.Println("Erro prepare DB:", err)
	} else {
		defer stmt.Close()

		_, err = stmt.ExecContext(ctxDB, valor)
		if err != nil {
			log.Println("Erro insert DB:", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"bid": valor,
	})
}
