package services

import (
	"encoding/json"
	"log"
	"net/http"
)

type CsrvClient struct {
	Secret string
}

func NewCsrvClient(secret string) *CsrvClient {
	return &CsrvClient{Secret: secret}
}

type VoucherResponse struct {
	Code string `json:"code"`
}

func (c *CsrvClient) GetCSRVCode() (string, error) {
	req, err := http.NewRequest("POST", "https://craftserve.pl/api/generate_voucher", nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("csrvbot", c.Secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("getCSRVCode http.DefaultClient.Do(req) " + err.Error())
		return "", err
	}

	var code VoucherResponse
	err = json.NewDecoder(resp.Body).Decode(&code)
	if err != nil {
		return "", err
	}
	err = resp.Body.Close()
	if err != nil {
		log.Println("getCSRVCode resp.Body.Close() ", err)
	}
	return code.Code, nil
}
