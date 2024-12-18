package economic

import (
	"fmt"
	"net/http"
	"strings"
)

const DIMENSIONAPI_BASE = "/dimensionsapi/v4.3.0"

func (client *Client) CreateDimensionValue(number, key int, name string) error {
	body := map[string]any{
		"active":          true,
		"dimensionNumber": number,
		"key":             key,
		"name":            name,
	}
	return client.callAPI(DIMENSIONAPI_BASE+"/values", http.MethodPost, nil, body, nil)
}

// Creates dimension value if doesn't exist.
// Returns true if the value was created, or false if it already exists.
// Name is not changed/updated if the value already exists.
func (client *Client) CreateDimensionValueIfItDoesNotExist(number, key int, name string) (bool, error) {
	err := client.callAPI(fmt.Sprintf(DIMENSIONAPI_BASE+"/values/%d/%d", number, key), http.MethodGet, nil, nil, nil)
	if err == nil {
		return false, nil
	} else if err != nil {
		// XXX TODO nicer 404 handling
		if !strings.Contains(err.Error(), "not found") {
			return false, err
		}
	}
	return true, client.CreateDimensionValue(number, key, name)
}

func (client *Client) AddDimensionValueToDraftEntry(dimensionNumber, dimensionKey, journalNumber, entryNumber int) error {
	body := map[string]any{
		"dimensionNumber": dimensionNumber,
		"dimensionKey":    dimensionKey,
		"journalNumber":   journalNumber,
		"entryNumber":     entryNumber,
	}
	return client.callAPI(fmt.Sprintf(DIMENSIONAPI_BASE+"/dimension-data/draft-entries"), http.MethodPost, nil, body, nil)
}
