package cdekapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type CDEKClient struct {
	Token      string
	TestMode   bool
	APIAddress string
}

type Size struct {
	Height int
	Length int
	Weight int
	Width  int
}

type PriceSending struct {
	TariffCode   int     `json:"tariff_code"`
	TariffName   string  `json:"tariff_name"`
	DeliveryMode int     `json:"delivery_mode"`
	DeliverySum  float64 `json:"delivery_sum"`
	PeriodMin    int     `json:"period_min"`
	PeriodMax    int     `json:"period_max"`
	CalendarMin  int     `json:"calendar_min"`
	CalendarMax  int     `json:"calendar_max"`
}

type requestData struct {
	Type         int      `json:"type"`
	Date         string   `json:"date"`
	Currency     int      `json:"currency"`
	Lang         string   `json:"lang"`
	FromLocation location `json:"from_location"`
	ToLocation   location `json:"to_location"`
	Packages     []Size   `json:"packages"`
}

type CDEKAuth struct {
	GrantType    string
	ClientID     string
	ClientSecret string
	TestMode     bool
	APIURL       string
}

type CDEKTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func NewCDEKAuth(grantType, clientID, clientSecret string, apiURL string) *CDEKAuth {
	return &CDEKAuth{
		GrantType:    grantType,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		APIURL:       apiURL,
	}
}

func (c *CDEKAuth) GetToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", c.GrantType)
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)

	req, err := http.NewRequest("POST", c.APIURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get token")
	}

	var tokenResp CDEKTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

type location struct {
	Code int `json:"code"`
}

type responseData struct {
	TariffCodes []PriceSending `json:"tariff_codes"`
}

func NewCDEKClient(token string, testMode bool, apiAddress string) *CDEKClient {
	return &CDEKClient{
		Token:      token,
		TestMode:   testMode,
		APIAddress: apiAddress,
	}
}

func (c *CDEKClient) Calculate(addressFrom string, addressTo string, size Size) ([]PriceSending, error) {
	addressFromInt, _ := strconv.Atoi(addressFrom)
	addressToInt, _ := strconv.Atoi(addressTo)
	request := requestData{
		Type:     1,
		Date:     time.Now().Format("2006-01-02T15:04:05-0700"),
		Currency: 1,
		Lang:     "rus",
		FromLocation: location{
			Code: addressFromInt,
		},
		ToLocation: location{
			Code: addressToInt,
		},
		Packages: []Size{size},
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.APIAddress, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error: non-200 status code")
	}

	var response responseData
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.TariffCodes, nil
}
