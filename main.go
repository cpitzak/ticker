package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	ticker := strings.ToUpper(args[0])
	// Fetch stock quote
	stockQuote, err := getStockQuote(ticker)
	if err != nil {
		fmt.Println("Error fetching stock quote:", err)
		os.Exit(1)
	}
	fmt.Println(stockQuote)
}

func getStockQuote(ticker string) (string, error) {
	var client = &http.Client{Timeout: 10 * time.Second}
	url := "https://query2.finance.yahoo.com/v6/finance/quoteSummary/" + ticker + "?modules=financialData"
	r, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body) // response body is []byte
	if err != nil {
		return "", err
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		return "", err
	}

	// Access the nested value
	quoteSummary, ok := jsonMap["quoteSummary"].(map[string]interface{})
	if !ok {
		return "", errors.New("quoteSummary not found or not a map")
	}

	resultArray, ok := quoteSummary["result"].([]interface{})
	if !ok || len(resultArray) == 0 {
		return "", errors.New("result not found or not an array or is empty")
	}

	result0, ok := resultArray[0].(map[string]interface{})
	if !ok {
		return "", errors.New("result[0] not found or not a map")
	}

	financialData, ok := result0["financialData"].(map[string]interface{})
	if !ok {
		return "", errors.New("financialData not found or not a map")
	}

	currentPrice, ok := financialData["currentPrice"].(map[string]interface{})
	if !ok {
		return "", errors.New("currentPrice not found or not a map")
	}

	raw, ok := currentPrice["raw"].(float64)
	if !ok {
		return "", errors.New("raw not found or not a float64")
	}
	price := fmt.Sprintf("%.2f", raw)
	return ticker + ": $" + price, nil
}
