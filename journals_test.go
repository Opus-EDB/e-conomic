package economic

import (
	"encoding/json"
	"testing"
)

func TestCreateJournal(t *testing.T) {
	j := JournalEntry{
		EntryTypeNumber:     5,
		VoucherNumber:       50160,
		JournalNumber:       6,
		Date:                "2024-09-26",
		Amount:              json.Number("500"),
		Currency:            "DKK",
		AccountNumber:       4610,
		ContraAccountNumber: 4630,
		ContraVatCode:       "U25",
		VatCode:             "U25",
	}
	client := getTestClient()
	err := j.CreateEntry(client)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	defer j.Delete(client)
	found, err := client.GetCashPaymentById(j.VoucherNumber)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if found.VoucherNumber != j.VoucherNumber {
		t.Fatalf("Expected %d, got %d", j.VoucherNumber, found.VoucherNumber)
	}
	if found.VatCode != j.VatCode {
		t.Fatalf("Expected %s, got %s", j.VatCode, found.VatCode)
	}
}

func TestGetBookedCashPaymentById(t *testing.T) {
	j := JournalEntry{
		EntryTypeNumber:     5,
		VoucherNumber:       50160,
		JournalNumber:       6,
		Date:                "2024-09-26",
		Amount:              json.Number("500"),
		Currency:            "DKK",
		AccountNumber:       4610,
		ContraAccountNumber: 4630,
		ContraVatCode:       "U25",
		VatCode:             "U25",
	}
	client := getTestClient()
	err := j.CreateEntry(client)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	err = client.BookAllEntries(6)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	found, err := client.GetBookedCashPaymentById(j.VoucherNumber)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if found.VoucherNumber != j.VoucherNumber {
		t.Fatalf("Expected %d, got %d", j.VoucherNumber, found.VoucherNumber)
	}
}

func TestCreditBookedCashPayment(t *testing.T) {
	paymentId := 50160
	client := getTestClient()
	je, err := client.GetCashPaymentById(paymentId)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	j := JournalEntry{
		EntryTypeNumber:     5,
		VoucherNumber:       50160,
		JournalNumber:       6,
		Date:                "2024-09-26",
		Amount:              je.Amount,
		Currency:            "DKK",
		AccountNumber:       4610,
		ContraAccountNumber: 4630,
		ContraVatCode:       "U25",
		VatCode:             "U25",
		IsCredit:            true,
	}
	err = j.CreateEntry(client)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

}
