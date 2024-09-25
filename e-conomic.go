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

type Layout struct {
	LayoutNumber int    `json:"layoutNumber"`   //A unique identifier of the layout."`
	Self         string `json:"self,omitempty"` //A unique reference to the layout resource."`
}

// VatZone represents a VAT zone.
type VatZone struct {
	VatZoneNumber int    `json:"vatZoneNumber"`  // The unique identifier of the VAT-zone.
	Self          string `json:"self,omitempty"` // A unique link reference to the VAT-zone item.
}

// PaymentTerms represents the default payment terms for the customer.
type PaymentTerms struct {
	PaymentTermsNumber int    `json:"paymentTermsNumber"` // The unique identifier of the payment terms.
	Self               string `json:"self,omitempty"`     // A unique link reference to the payment terms item.
}

// SalesPerson represents the employee responsible for contact with the customer.
type SalesPerson struct {
	EmployeeNumber int    `json:"employeeNumber"` // The unique identifier of the employee.
	Self           string `json:"self,omitempty"` // A unique link reference to the employee resource.
}

func callRestAPI(endpoint, method string, request, response any) error {
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
	grant := c.AgreementGrant
	secret := c.AppSecretToken
	req.Header.Set("X-AppSecretToken", secret)
	req.Header.Set("X-AgreementGrantToken", grant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.DefaultClient
	res, err := client.Do(req)
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

// This is used by for collections in the Rest API
type CollectionReponse[T any] struct {
	Collection []T        `json:"collection"`
	MetaData   any        `json:"metaData"`
	Pagination Pagination `json:"pagination"`
	Self       string     `json:"self"`
}

// This is used by for collections in the Regular API
type ItemsReponse[T any] struct {
	Items      []T        `json:"items"`
	MetaData   any        `json:"metaData"`
	Pagination Pagination `json:"pagination"`
	Self       string     `json:"self"`
}

type Pagination struct {
	FirstPage            string `json:"firstPage"`
	LastPage             string `json:"lastPage"`
	MaxPageSizeAllowed   int    `json:"maxPageSizeAllowed"`
	PageSize             int    `json:"pageSize"`
	Results              int    `json:"results"`
	ResultsWithoutFilter int    `json:"resultsWithoutFilter"`
	SkipPages            int    `json:"skipPages"`
}

func callAPI(endpoint string, method string, params url.Values, body interface{}, response interface{}) error {
	client := http.DefaultClient
	grant := c.AgreementGrant
	secret := c.AppSecretToken
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
			return err
		}
		req.Body = io.NopCloser(bytes.NewReader(jsonBody))
	}
	resp, err := client.Do(req)
	if err != nil {
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

type ErrorResp struct {
	Errors []struct {
		Property  string `json:"property"`
		Message   string `json:"message"`
		ErrorCode string `json:"errorCode"`
	} `json:"errors"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	Status       int    `json:"status"`
	Instance     string `json:"instance"`
	TraceID      string `json:"traceId"`
	TraceTimeUtc string `json:"traceTimeUtc"`
	ErrorCode    string `json:"errorCode"`
}
