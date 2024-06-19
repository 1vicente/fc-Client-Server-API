package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Dolar struct {
	DolarChamada string `json:"bid"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	ChecaErros(err)

	res, err := http.DefaultClient.Do(req)
	ChecaErros(err)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	ChecaErros(err)

	var d Dolar
	err = json.Unmarshal(body, &d)
	ChecaErros(err)

	file, err := os.Create("cotacao.txt")
	ChecaErros(err)

	file.WriteString("Dólar: " + d.DolarChamada)
	defer file.Close()

}

func ChecaErros(e error) {
	if e != nil {
		fmt.Println("Houve uma falha: %v", e)
	}

}
