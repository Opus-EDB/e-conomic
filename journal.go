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
	IsCredit            bool        `json:"isCredit,omitempty"`
}

// Create a draft of a cash payment.
// If the entry is created successfully, the EntryNumber field will be set.
// Set IsCredit to true if the amount should be negative.
func (client *Client) CreateJournalEntry(j *JournalEntry) error {
	resp := map[string]any{}
	err := client.callAPI("/journalsapi/v6.0.0/draft-entries/", http.MethodPost, nil, j, &resp)
	if err == nil {
		entryNumber := resp["entryNumber"]
		log.Printf("entryNumber: %#v", entryNumber)
		if entryNumber != nil {
			j.EntryNumber = int(entryNumber.(float64))
		}
	}
	return err
}

func (client *Client) DeleteJournalEntry(j *JournalEntry) error {
	return client.callAPI(fmt.Sprintf("/journalsapi/v6.0.0/draft-entries/%d", j.EntryNumber), http.MethodDelete, nil, nil, nil)
}

func (client *Client) GetDraftEntriesCount() (int, error) {
	var count int
	err := client.callAPI("/journalsapi/v6.0.0/draft-entries/count", http.MethodGet, nil, nil, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (client *Client) GetCashPaymentById(id int) (JournalEntry, error) {
	je := JournalEntry{}
	resp := ItemsReponse[JournalEntry]{}
	params := url.Values{
		"filter": {fmt.Sprintf("voucherNumber$eq:%d", id)},
	}
	err := client.callAPI("/journalsapi/v6.0.0/draft-entries", http.MethodGet, params, nil, &resp)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	if len(resp.Items) == 0 {
		return je, fmt.Errorf("no payment with id %d", id)
	}
	je = resp.Items[0]
	return je, err
}

// Get a booked cash payment by its voucher number.
// There are basically two fields that may be interesting to you:
// - VoucherNumber: The voucher number of the payment.
// - Amount: The amount of the payment.
//
// If you need to credit the payment fill in the remaining fields and set IsCredit to true.
// The amount will always be negative when IsCredit=true.
func (client *Client) GetBookedCashPaymentById(id int) (JournalEntry, error) {
	je := JournalEntry{}
	resp := ItemsReponse[JournalEntry]{}
	params := url.Values{
		"filter": {fmt.Sprintf("voucherNumber$eq:%d", id)},
	}
	err := client.callAPI("/bookedEntriesapi/v2.0.0/booked-entries", http.MethodGet, params, nil, &resp)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	if len(resp.Items) == 0 {
		return je, fmt.Errorf("no payment with id %d", id)
	}
	je = resp.Items[0]
	return je, err
}

func (client *Client) BookAllEntries(journalNumber int) error {
	return client.callAPI(fmt.Sprintf("/journalsapi/v6.0.0/journals/%d/book", journalNumber), http.MethodPost, nil, nil, nil)
}
