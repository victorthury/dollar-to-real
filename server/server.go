package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const API = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type EconomiaUsdBrl struct {
	UsdBrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func getDollarToRealCotation() (*EconomiaUsdBrl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", API, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cotacao EconomiaUsdBrl
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cotacao, err := getDollarToRealCotation()
	if err != nil {
		fmt.Println(err)
	}

	if cotacao == nil {
		log.Println("Tempo excedeu os 200ms")
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacao)
	w.WriteHeader(http.StatusOK)
}

func main() {
	log.Println("Server has started on localhost:8080")
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}
