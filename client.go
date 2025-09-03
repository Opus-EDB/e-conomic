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

func (client *Client) callRestAPI(endpoint, method string, request, response any) error { // also return status code?
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
	log.Printf("e-conomic/REST %s %s => %d", method, endpoint, res.StatusCode)
	if res.StatusCode >= 400 {
		log.Printf("error calling e-conomic (%s %s) err: %s", url, method, body.String())
		return fmt.Errorf("error calling e-conomic (%s %s) err: %s", url, method, body.String())
	}
	if response == nil {
		return nil
	}
	err = json.Unmarshal(body.Bytes(), response)
	return err
}

func (client *Client) callAPI(endpoint string, method string, params url.Values, body interface{}, response interface{}) error {
	if params == nil {
		params = url.Values{}
	}
	client.assertClientIsConfigured()
	grant := client.AgreementGrant
	secret := client.AppSecretToken
	if grant == "" || secret == "" {
		panic("missing agreement grant or app secret")
	}
	req := &http.Request{
		Method: method,
		URL: &url.URL{
			Path:     endpoint,
			RawQuery: params.Encode(),
			Scheme:   "https",
			Host:     "apis.e-conomic.com",
		},
		Header: make(http.Header),
	}
	req.Header.Set("X-AppSecretToken", secret)
	req.Header.Set("X-AgreementGrantToken", grant)
	req.Header.Set("Accept", "application/json")
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
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("error in calling e-conomic (%s %s) err: %s", endpoint, method, err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body (internal error?) when calling e-conomic (%s %s => %d)", method, endpoint, res.StatusCode)
		}
		return fmt.Errorf("error in calling e-conomic (%s %s => %d) %s", method, endpoint, res.StatusCode, string(body))
	}
	if response != nil {
		err = json.NewDecoder(res.Body).Decode(response)
	}
	log.Printf("e-conomic/OpenAPI %s %s => %d", method, endpoint, res.StatusCode)
	return err
}

func (tc *TypedClient[T]) getEntities(baseUrl string, pageSize int) (entities []T, err error) { // generalize more by adding filter param?
	client := tc.client
	results := CollectionReponse[T]{}
	filter := &Filter{}
	filter.AndCondition("pagesize", FilterOperatorEquals, pageSize)
	err = client.callRestAPI(fmt.Sprintf(baseUrl+"?"+filter.filterStr), http.MethodGet, nil, &results)
	if err != nil {
		log.Printf("ERROR: %#v", err)
		return
	}
	entities = results.Collection
	numberOfResults := results.Pagination.Results
	numberOfPages := (numberOfResults / pageSize) + 1 // integer division (disregarding the remainder)
	for i := range numberOfPages {
		url := fmt.Sprintf("%s?skippages=%d&pagesize=%d", baseUrl, i, pageSize)
		fmt.Printf("URL: %s\n", url)
		err = client.callRestAPI(url, http.MethodGet, nil, &results)
		if err != nil {
			log.Printf("ERROR: %#v", err)
			return
		}
		entities = append(entities, results.Collection...)
	}
	return
}
