package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
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
}

type EconomiaUsdBrl struct {
	UsdBrl Cotacao `json:"USDBRL"`
}

func main() {
	bid, err := getQuotation()
	if err != nil {
		log.Println(err)
		return
	}
	err = writeQuotation(bid)
	if err != nil {
		log.Println(err)
		return
	}
}

func getQuotation() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var data EconomiaUsdBrl
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return "", err
	}

	return data.UsdBrl.Bid, nil
}

func writeQuotation(bid string) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s\n", bid))
	if err != nil {
		return err
	}
	log.Println("Cotação escrita com sucesso!")
	return nil
}
