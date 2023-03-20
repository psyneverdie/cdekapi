package main

import (
	"cdekapi/package"
	"encoding/json"
	"fmt"
)

func main() {
	grantType := "client_credentials"
	clientID := "EMscd6r9JnFiQ3bLoyjJY6eM78JrJceI"
	clientSecret := "PjLZkKBHEiLK3YsjtNrt3TGNG0ahs3kG"
	apiAuthURL := "https://api.edu.cdek.ru/v2/oauth/token?parameters"
	credentials := cdekapi.NewCDEKAuth(grantType, clientID, clientSecret, apiAuthURL) // получение токена
	token, err := credentials.GetToken()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	apiCalcURL := "https://api.edu.cdek.ru/v2/calculator/tarifflist"
	client := cdekapi.NewCDEKClient(token, true, apiCalcURL) // формирование данных
	size := cdekapi.Size{                                    // формирование данных посылки
		Height: 10,
		Length: 10,
		Weight: 4000,
		Width:  10,
	}

	addressFrom := "270"
	addressTo := "40"
	prices, err := client.Calculate(addressFrom, addressTo, size) // калькулятор запроса CDEK API
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result := struct {
		TariffCodes []cdekapi.PriceSending `json:"tariff_codes"` // формирование данных
	}{
		TariffCodes: prices,
	}

	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonOutput))
}
