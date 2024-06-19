package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dolar struct {
	DolarChamadaID int
	DolarChamada   DolarChamada `json:"USDBRL"`
}

type DolarChamada struct {
	gorm.Model
	ID          int    `gorm:"primaryKey"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	High        string `json:"high"`
	Codein      string `json:"codein"`
	Low         string `json:"low"`
	Varbid      string `json:"varbid"`
	PctChange   string `json:"pctchange"`
	Bid         string `json:"bid"`
	Ask         string `json:"ask"`
	Timestamp   string `json:"timestamp"`
	Create_date string `json:"create_date"`
}

type DolarDTO struct {
	BidDTO string `json:"bid"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /cotacao", DefaultHandler)
	http.ListenAndServe(":8080", mux)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	retorno, err := PrecoDolar()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println(retorno)
	var b DolarDTO
	b.BidDTO = retorno.Bid

	json.NewEncoder(w).Encode(b)
}

func PrecoDolar() (*DolarChamada, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, error := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if error != nil {
		return nil, error
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, error
	}
	defer res.Body.Close()
	body, error := io.ReadAll(res.Body)
	if error != nil {
		return nil, error
	}

	var d Dolar
	error = json.Unmarshal(body, &d)
	if error != nil {
		return nil, error
	}
	dTemp := d.DolarChamada
	SalvarDb(&dTemp)

	return &dTemp, nil
}

func SalvarDb(dolar *DolarChamada) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&Dolar{}, &DolarChamada{})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	db.WithContext(ctx).Create(&dolar)
}
