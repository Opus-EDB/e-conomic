package economic

import (
	"fmt"
	"net/http"
	"strings"
)

const DIMENSIONAPI_BASE = "/dimensionsapi/v5.0.0"

type dimension struct {
	Active          bool   `json:"active"`
	DimensionNumber int    `json:"dimensionNumber"`
	Key             int    `json:"key"`
	Name            string `json:"name"`
	ObjectVersion   string `json:"objectVersion,omitempty"`
}

func (client *Client) CreateDimensionValue(number, key int, name string) error {
	body := dimension{
		Active:          true,
		DimensionNumber: number,
		Key:             key,
		Name:            name,
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
	}
	// XXX TODO nicer 404 handling
	if !strings.Contains(err.Error(), "not found") {
		return false, err
	}
	return true, client.CreateDimensionValue(number, key, name)
}

// Updates or creates a dimension value.
// Returns true if the value was created, or false if it already exists.
func (client *Client) CreateOrUpdateDimensionValue(number, key int, name string) (bool, error) {
	var existingDimension dimension
	err := client.callAPI(fmt.Sprintf(DIMENSIONAPI_BASE+"/values/%d/%d", number, key), http.MethodGet, nil, nil, &existingDimension)
	if err == nil {
		return client.UpdateDimensionValue(number, key, name, existingDimension.ObjectVersion)
	}
	// XXX TODO nicer 404 handling
	if !strings.Contains(err.Error(), "not found") {
		return false, err
	}
	return true, client.CreateDimensionValue(number, key, name)
}

func (client *Client) UpdateDimensionValue(number, key int, name, objectVersion string) (bool, error) {
	body := dimension{
		Active:          true,
		DimensionNumber: number,
		Key:             key,
		Name:            name,
		ObjectVersion:   objectVersion,
	}
	return false, client.callAPI(DIMENSIONAPI_BASE+"/values", http.MethodPut, nil, body, nil)
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
