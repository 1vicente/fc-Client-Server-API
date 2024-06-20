package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	var b DolarDTO
	b.BidDTO = retorno.Bid

	json.NewEncoder(w).Encode(b)
}

func PrecoDolar() (*DolarChamada, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, error := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if error != nil {
		return nil, errors.New("falha ao montar a requisição")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("falha ao realziar a requisição")
	}
	defer res.Body.Close()
	body, error := io.ReadAll(res.Body)
	if error != nil {
		return nil, errors.New("falha ao ler o body")
	}

	var d Dolar
	error = json.Unmarshal(body, &d)
	if error != nil {
		return nil, errors.New("falha ao realizar unmarshal")
	}
	dTemp := d.DolarChamada
	salvou := SalvarDb(&dTemp)
	if !salvou {
		return &dTemp, errors.New("falha ao salvar dados")
	}

	return &dTemp, nil
}

func SalvarDb(dolar *DolarChamada) bool {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&Dolar{}, &DolarChamada{})
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	save := db.WithContext(ctx).Create(&dolar)
	if save.Error != nil {
		return false
	}
	return true
}
