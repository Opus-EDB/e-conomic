package economic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// {
//   "entryTypeNumber": 0,
//   "voucherNumber": 0,
//   "journalNumber": 0,
//   "date": "2019-08-24T14:15:22Z",
//   "amount": 0,
//   "currency": "string"
// }

type JournalEntry struct {
	EntryTypeNumber     int         `json:"entryTypeNumber,omitempty"`
	VoucherNumber       int         `json:"voucherNumber"`
	JournalNumber       int         `json:"journalNumber"`
	Date                string      `json:"date"`
	Amount              json.Number `json:"amount"`
	Currency            string      `json:"currency"`
	EntryNumber         int         `json:"entryNumber,omitempty"`
	AccountNumber       int         `json:"accountNumber,omitempty"`
	ContraAccountNumber int         `json:"contraAccountNumber,omitempty"`
	Text                string      `json:"text,omitempty"`
	VatCode             string      `json:"vatCode,omitempty"`
	ContraVatCode       string      `json:"contraVatCode,omitempty"`
}

func (j *JournalEntry) CreateEntry() error {
	resp := map[string]any{}
	err := callAPI("/journalsapi/v6.0.0/draft-entries/", http.MethodPost, url.Values{}, j, &resp)
	if err == nil {
		entryNumber := resp["entryNumber"]
		log.Printf("entryNumber: %#v", entryNumber)
		if entryNumber != nil {
			j.EntryNumber = int(entryNumber.(float64))
		}
	}
	return err
}

func (j *JournalEntry) Delete() error {
	return callAPI(fmt.Sprintf("/journalsapi/v6.0.0/draft-entries/%d", j.EntryNumber), http.MethodDelete, url.Values{}, nil, nil)
}

func GetDraftEntriesCount() (int, error) {
	var count int
	err := callAPI("/journalsapi/v6.0.0/draft-entries/count", http.MethodGet, nil, nil, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetCashPaymentById(id int) (JournalEntry, error) {
	je := JournalEntry{}
	resp := ItemsReponse[JournalEntry]{}
	params := url.Values{
		"filter": {fmt.Sprintf("voucherNumber$eq:%d", id)},
	}
	err := callAPI("/journalsapi/v6.0.0/draft-entries", http.MethodGet, params, nil, &resp)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	if len(resp.Items) == 0 {
		return je, fmt.Errorf("no payment with id %d", id)
	}
	je = resp.Items[0]
	return je, err
}
