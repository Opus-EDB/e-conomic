package economic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type APIError struct {
	StatusCode int
	body       string
}

func (e *APIError) Error() string {
	return e.body
}

type Client struct {
	AgreementGrant string `json:"agreement_grant"`
	AppSecretToken string `json:"app_secret"`
}

func (client *Client) assertClientIsConfigured() {
	if len(client.AgreementGrant) == 0 || len(client.AppSecretToken) == 0 {
		panic(fmt.Sprintf("client is not configured: %#v", client))
	}
}

const (
	DEFAULT_PAGE_SIZE = 100
	maxRetries        = 4
	baseDelay         = time.Second
)

// isIdempotent reports whether repeating a request is safe when we cannot tell
// whether the first one took effect. GET, HEAD, PUT and DELETE are idempotent by
// definition; POST and PATCH are not, and repeating one that was already applied
// can create a second record.
//
// PATCH is excluded on principle rather than on evidence: a patch document is
// free to append or increment, so HTTP does not call it idempotent, even though
// the only one we send today (SetEInvoicing, a "replace" of a single field) would
// be safe to repeat. The cost of that caution is one failed SetEInvoicing on a
// 504; the cost of the opposite mistake is duplicated data nobody notices.
func isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete, http.MethodOptions:
		return true
	}
	return false
}

// isRetryableStatus says whether a failed attempt is worth repeating.
//
// What decides it is whether E-conomic itself received the request. 408 is their
// "DownstreamServiceTimeout" and 504 is the gateway's "the origin web server did
// not respond within the allowed time" — both mean the request reached them and
// the outcome is unknown: it may well have been applied. Repeating a
// non-idempotent request in that state duplicates it, which is exactly how order
// 1185688 ended up with two identical Billetgebyr journal entries on 2026-09-08,
// the second posted by this retry five minutes after the first.
//
// Everything else stays retryable for every method. A 429 was refused before any
// work was done, and the connection-level failures that dominate here ("upstream
// connect error or disconnect/reset before headers") never reach the origin —
// one of those retried a journal entry on 2026-08-23 and correctly created
// nothing twice. 502 and 503 are left in for the same reason, though they are
// the least certain of the set: widen this if a duplicate ever turns up behind
// one.
func isRetryableStatus(method string, code int) bool {
	switch code {
	case http.StatusRequestTimeout, http.StatusGatewayTimeout: // 408, 504
		return isIdempotent(method)
	}
	return code == http.StatusTooManyRequests || code >= 500
}

// backoffDelay returns how long to wait before the next attempt. For 429
// responses it respects the Retry-After header when present; otherwise it
// uses plain exponential backoff (1s, 2s, 4s, 8s).
func backoffDelay(attempt int, res *http.Response) time.Duration {
	if res != nil && res.StatusCode == http.StatusTooManyRequests {
		if s := res.Header.Get("Retry-After"); s != "" {
			if secs, err := strconv.Atoi(s); err == nil {
				return time.Duration(secs) * time.Second
			}
		}
	}
	return baseDelay * (1 << attempt) // left bit shift
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

	var lastErr error
	var lastRes *http.Response
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := backoffDelay(attempt-1, lastRes)
			log.Printf("retrying e-conomic/REST %s %s (attempt %d/%d) after %s", method, endpoint, attempt, maxRetries, delay)
			time.Sleep(delay)
			lastRes = nil
		}

		req, err := http.NewRequest(method, url, bytes.NewReader(jsonRequest))
		if err != nil {
			return err
		}
		req.Header.Set("X-AppSecretToken", client.AppSecretToken)
		req.Header.Set("X-AgreementGrantToken", client.AgreementGrant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("error in calling e-conomic (%s %s) err: %s", url, method, err)
			lastErr = err
			continue
		}

		body := new(bytes.Buffer)
		body.ReadFrom(res.Body)
		res.Body.Close()
		log.Printf("e-conomic/REST %s %s => %d", method, endpoint, res.StatusCode)

		if isRetryableStatus(method, res.StatusCode) {
			lastErr = &APIError{
				StatusCode: res.StatusCode,
				body:       fmt.Sprintf("error calling e-conomic (%s %s) err: %s", url, method, body.String()),
			}
			log.Printf("will retry e-conomic/REST %s %s (attempt %d/%d): %s", method, endpoint, attempt+1, maxRetries, body.String())
			lastRes = res
			continue
		}

		if res.StatusCode >= 400 {
			log.Printf("error calling e-conomic (%s %s) err: %s", url, method, body.String())
			return &APIError{
				StatusCode: res.StatusCode,
				body:       fmt.Sprintf("error calling e-conomic (%s %s) err: %s", url, method, body.String()),
			}
		}

		if response == nil {
			return nil
		}
		return json.Unmarshal(body.Bytes(), response)
	}
	return lastErr
}

func (client *Client) callAPI(endpoint string, method string, params url.Values, body any, response any) error {
	if params == nil {
		params = url.Values{}
	}
	client.assertClientIsConfigured()
	grant := client.AgreementGrant
	secret := client.AppSecretToken
	if grant == "" || secret == "" {
		panic("missing agreement grant or app secret")
	}

	var jsonBody []byte
	if body != nil {
		var err error
		jsonBody, err = json.Marshal(body)
		if err != nil {
			log.Printf("error in marshalling request: %s", err)
			return err
		}
	}

	reqURL := &url.URL{
		Path:     endpoint,
		RawQuery: params.Encode(),
		Scheme:   "https",
		Host:     "apis.e-conomic.com",
	}

	var lastErr error
	var lastRes *http.Response
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := backoffDelay(attempt-1, lastRes)
			log.Printf("retrying e-conomic/OpenAPI %s %s (attempt %d/%d) after %s", method, endpoint, attempt, maxRetries, delay)
			time.Sleep(delay)
			lastRes = nil
		}

		req := &http.Request{
			Method: method,
			URL:    reqURL,
			Header: make(http.Header),
		}
		req.Header.Set("X-AppSecretToken", secret)
		req.Header.Set("X-AgreementGrantToken", grant)
		req.Header.Set("Accept", "application/json")
		if jsonBody != nil {
			req.Header.Set("Content-Type", "application/json")
			req.Body = io.NopCloser(bytes.NewReader(jsonBody))
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("error in calling e-conomic (%s %s) err: %s", endpoint, method, err)
			lastErr = err
			continue
		}

		if isRetryableStatus(method, res.StatusCode) {
			resBody, _ := io.ReadAll(res.Body)
			res.Body.Close()
			lastErr = fmt.Errorf("error in calling e-conomic (%s %s => %d) %s", method, endpoint, res.StatusCode, string(resBody))
			log.Printf("will retry e-conomic/OpenAPI %s %s (attempt %d/%d): %s", method, endpoint, attempt+1, maxRetries, string(resBody))
			lastRes = res
			continue
		}

		if res.StatusCode >= 400 {
			resBody, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return fmt.Errorf("failed to read response body (internal error?) when calling e-conomic (%s %s => %d)", method, endpoint, res.StatusCode)
			}
			return fmt.Errorf("error in calling e-conomic (%s %s => %d) %s", method, endpoint, res.StatusCode, string(resBody))
		}

		if response != nil {
			err = json.NewDecoder(res.Body).Decode(response)
		}
		res.Body.Close()
		log.Printf("e-conomic/OpenAPI %s %s => %d", method, endpoint, res.StatusCode)
		return err
	}
	return lastErr
}

// getEntities uses the REST API (callRestAPI) and CollectionReponse. Determines remaining pages from Pagination.Results in the first response.
func (tc *TypedClient[T]) getEntities(baseUrl string, pageSize int, filter string) (entities []T, err error) {
	client := tc.client
	results := CollectionReponse[T]{}
	url := baseUrl + fmt.Sprintf("?filter=%s&pagesize=%d", filter, pageSize)
	err = client.callRestAPI(url, http.MethodGet, nil, &results)
	if err != nil {
		log.Printf("ERROR: %#v", err)
		return
	}
	entities = results.Collection
	numberOfResults := results.Pagination.Results
	numberOfPages := (numberOfResults / pageSize) + 1 // integer division (disregarding the remainder)
	if numberOfPages > 1 {
		for i := 1; i < numberOfPages; i++ {
			url := fmt.Sprintf("%s?skippages=%d&pagesize=%d", baseUrl, i, pageSize)
			fmt.Printf("URL: %s\n", url)
			err = client.callRestAPI(url, http.MethodGet, nil, &results)
			if err != nil {
				log.Printf("ERROR: %#v", err)
				return
			}
			entities = append(entities, results.Collection...)
		}
	}
	return
}

// getAllCursor fetches all items from a cursor-based pagination endpoint.
// The cursor is absent in the response when there are no more items.
func getAllCursor[T any](client *Client, baseURL string, params url.Values) ([]T, error) {
	var all []T
	cursor := ""
	for {
		p := url.Values{}
		for k, v := range params {
			p[k] = v
		}
		if cursor != "" {
			p.Set("cursor", cursor)
		}
		resp := CursorResponse[T]{}
		if err := client.callAPI(baseURL, http.MethodGet, p, nil, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Items...)
		if resp.Cursor == "" {
			break
		}
		cursor = resp.Cursor
	}
	return all, nil
}

// function to get the last entity (e.g. customer contact)?
