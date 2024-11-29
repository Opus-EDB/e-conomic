package economic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	AgreementGrant string `json:"agreement_grant"`
	AppSecretToken string `json:"app_secret"`
}

func (client *Client) assertClientIsConfigured() {
	if len(client.AgreementGrant) == 0 || len(client.AppSecretToken) == 0 {
		panic(fmt.Sprintf("client is not configured: %#v", client))
	}
}

func (client *Client) callRestAPI(endpoint, method string, request, response any) error {
	client.assertClientIsConfigured()
	url := fmt.Sprintf("https://restapi.e-conomic.com/%s", endpoint)
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Printf("error in marshalling request: %s", err)
		return err
	}
	if request == nil {
		jsonRequest = []byte{}
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonRequest))
	if err != nil {
		return err
	}
	grant := client.AgreementGrant
	secret := client.AppSecretToken
	req.Header.Set("X-AppSecretToken", secret)
	req.Header.Set("X-AgreementGrantToken", grant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("error in calling e-conomic (%s %s) err: %s", url, method, err)
		return err
	}
	defer res.Body.Close()
	body := new(bytes.Buffer)
	body.ReadFrom(res.Body)
	log.Printf("status code from e-conomic: %d", res.StatusCode)
	if res.StatusCode >= 400 {
		log.Printf("error in calling e-conomic (%s %s) err: %s", url, method, body.String())
		return fmt.Errorf("error in calling e-conomic (%s %s) err: %s", url, method, body.String())
	}
	if response == nil {
		return nil
	}
	err = json.Unmarshal(body.Bytes(), response)
	return err
}

func (client *Client) callAPI(endpoint string, method string, params url.Values, body interface{}, response interface{}) error {
	client.assertClientIsConfigured()
	grant := client.AgreementGrant
	secret := client.AppSecretToken
	if grant == "" || secret == "" {
		panic("missing agreement grant or app secret")
	}
	req := &http.Request{
		Method: method,
		URL: &url.URL{Path: endpoint, RawQuery: params.Encode(),
			Scheme: "https", Host: "apis.e-conomic.com"},
		Header: make(http.Header),
	}
	req.Header.Set("X-AppSecretToken", secret)
	req.Header.Set("X-AgreementGrantToken", grant)
	req.Header.Set("Accept", "application/json")
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		jsonBody, err := json.Marshal(body)
		if err != nil {
			log.Printf("error in marshalling request: %s", err)
			return err
		}
		req.Body = io.NopCloser(bytes.NewReader(jsonBody))
	}
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("error in calling e-conomic (%s %s) err: %s", endpoint, method, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		errResp := ErrorResp{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return fmt.Errorf("error in calling e-conomic (%s %s) err: %s", endpoint, method, err)
		}
		errs := errResp.Title + " "
		for _, e := range errResp.Errors {
			errs += fmt.Sprintf("%s: %s\n", e.Property, e.Message)
		}
		return fmt.Errorf("error in calling e-conomic (%s %s) %s", endpoint, method, errs)
	}
	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(response)
	}
	log.Printf("status code from e-conomic: %d", resp.StatusCode)
	return err
}
