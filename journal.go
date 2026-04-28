package economic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const journalApiVersion = "v14.0.1"
const journalDraftEntryBaseUrl = "/journalsapi/" + journalApiVersion + "/draft-entries"
const bookedEntriesApiBaseUrl = "/bookedEntriesapi/v4.0.0/booked-entries"

type TimeWindow struct {
	From time.Time
	To   time.Time
}

func YesterdayWindow() TimeWindow {
	loc, _ := time.LoadLocation("Europe/Copenhagen")
	yesterday := time.Now().In(loc).AddDate(0, 0, -1)
	return TimeWindow{
		From: time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, loc),
		To:   time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, loc),
	}
}

func (client *Client) GetJournalEntries(journalNumber int, windows ...TimeWindow) ([]JournalEntry, error) {
	window := YesterdayWindow()
	if len(windows) > 0 {
		window = windows[0]
	}
	filter := fmt.Sprintf("date$gte:%s$and:date$lte:%s",
		window.From.Format(time.RFC3339),
		window.To.Format(time.RFC3339),
	)
	if journalNumber != 0 {
		filter += fmt.Sprintf("$and:journalNumber$eq:%d", journalNumber)
	}
	params := url.Values{"filter": {filter}}

	draft, err := getAllPaged[JournalEntry](client, journalDraftEntryBaseUrl+"/paged", params)
	if err != nil {
		return nil, err
	}
	if len(draft) > 0 {
		log.Printf("GetJournalEntries: %d draft entries found", len(draft))
	}
	booked, err := getAllPaged[JournalEntry](client, bookedEntriesApiBaseUrl+"/paged", params)
	if err != nil {
		return nil, err
	}
	if len(booked) > 0 {
		log.Printf("GetJournalEntries: %d booked entries found", len(booked))
	}
	return append(draft, booked...), nil
}

type JournalEntry struct {
	EntryTypeNumber     int         `json:"entryTypeNumber"`
	VoucherNumber       int         `json:"voucherNumber"`
	JournalNumber       int         `json:"journalNumber"`
	Date                string      `json:"date"`
	Amount              json.Number `json:"amount"`
	Currency            string      `json:"currency"`
	EntryNumber         int         `json:"entryNumber,omitempty"`
	AccountNumber       int         `json:"accountNumber"`
	ContraAccountNumber int         `json:"contraAccountNumber,omitempty"`
	Text                string      `json:"text,omitempty"`
	VatCode             string      `json:"vatCode,omitempty"`
	ContraVatCode       string      `json:"contraVatCode,omitempty"`
}

func truncateEntryText(j *JournalEntry) {
	const maxRunes = 240
	runes := []rune(j.Text)
	if len(runes) > maxRunes {
		j.Text = string(runes[:maxRunes])
	}
}

// Create a draft of a cash payment.
// If the entry is created successfully, the EntryNumber field will be set.
// Credits use negative amounts.
func (client *Client) CreateJournalEntry(j *JournalEntry) error {
	resp := map[string]any{}
	truncateEntryText(j)
	err := client.callAPI(journalDraftEntryBaseUrl, http.MethodPost, nil, j, &resp)
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
	return client.callAPI(fmt.Sprintf("%s/%d", journalDraftEntryBaseUrl, j.EntryNumber), http.MethodDelete, nil, nil, nil)
}

func (client *Client) GetDraftEntriesCount() (int, error) {
	var count int
	err := client.callAPI(journalDraftEntryBaseUrl+"/count", http.MethodGet, nil, nil, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (client *Client) GetCashPaymentById(id int) (JournalEntry, error) {
	je := JournalEntry{}
	jes, err := client.GetCashPaymentsById(id)
	if err != nil {
		return je, err
	}
	return jes[0], err
}

func (client *Client) GetCashPaymentsById(id int) ([]JournalEntry, error) {
	jes := []JournalEntry{}
	resp := ItemsReponse[JournalEntry]{}
	params := url.Values{
		"filter": {fmt.Sprintf("voucherNumber$eq:%d", id)},
	}
	err := client.callAPI(journalDraftEntryBaseUrl, http.MethodGet, params, nil, &resp)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	if len(resp.Items) == 0 {
		return jes, fmt.Errorf("no payment with id %d", id)
	}
	return resp.Items, err
}

// Get a booked cash payment by its voucher number.
// There are basically two fields that may be interesting to you:
// - VoucherNumber: The voucher number of the payment.
// - Amount: The amount of the payment.
//
// If you need to credit the payment fill in the remaining fields and use a negative amount.
func (client *Client) GetBookedCashPaymentById(id int) (JournalEntry, error) {
	je := JournalEntry{}
	jes, err := client.GetBookedCashPaymentsById(id)
	if err != nil {
		return je, err
	}
	return jes[0], err
}

func (client *Client) GetBookedCashPaymentsById(id int) ([]JournalEntry, error) {
	jes := []JournalEntry{}
	resp := ItemsReponse[JournalEntry]{}
	params := url.Values{
		"filter": {fmt.Sprintf("voucherNumber$eq:%d", id)},
	}
	err := client.callAPI(bookedEntriesApiBaseUrl, http.MethodGet, params, nil, &resp)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	if len(resp.Items) == 0 {
		return jes, fmt.Errorf("no payment with id %d", id)
	}
	return resp.Items, err
}

// GetDraftEntriesByVoucherNumber returns all draft entries with the given voucher number.
// Returns an empty slice and no error if none are found.
func (client *Client) GetDraftEntriesByVoucherNumber(voucherNumber int) ([]JournalEntry, error) {
	resp := ItemsReponse[JournalEntry]{}
	params := url.Values{
		"filter": {fmt.Sprintf("voucherNumber$eq:%d", voucherNumber)},
	}
	err := client.callAPI(journalDraftEntryBaseUrl, http.MethodGet, params, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// UpdateJournalEntry updates an existing draft entry using PUT. Needs an entryNumber (returned from GetDraftEntriesByVoucherNumber).
func (client *Client) UpdateJournalEntry(j *JournalEntry) error {
	truncateEntryText(j)
	return client.callAPI(journalDraftEntryBaseUrl, http.MethodPut, nil, j, nil)
}

func (client *Client) BookAllEntries(journalNumber int) error {
	return client.callAPI(fmt.Sprintf("/journalsapi/%s/journals/%d/book", journalApiVersion, journalNumber), http.MethodPost, nil, nil, nil)
}

func (client *Client) GetJournalBalanceById(id int) (float64, error) {
	resp := ItemsReponse[JournalEntry]{}
	params := url.Values{
		"filter": {fmt.Sprintf("voucherNumber$eq:%d", id)},
	}
	err := client.callAPI(journalDraftEntryBaseUrl, http.MethodGet, params, nil, &resp)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	balance := 0.0
	for _, item := range resp.Items {
		a, err := item.Amount.Float64()
		if err != nil {
			return 0, nil
		}
		balance += a
	}
	return balance, nil
}
