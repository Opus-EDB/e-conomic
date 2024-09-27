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
	err := j.CreateEntry()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	defer j.Delete()
	found, err := GetCashPaymentById(j.VoucherNumber)
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
	err := j.CreateEntry()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	err = BookAllEntries(6)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	found, err := GetBookedCashPaymentById(j.VoucherNumber)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if found.VoucherNumber != j.VoucherNumber {
		t.Fatalf("Expected %d, got %d", j.VoucherNumber, found.VoucherNumber)
	}
}

func TestCreditBookedCashPayment(t *testing.T) {
	paymentId := 50160
	je, err := GetCashPaymentById(paymentId)
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
	err = j.CreateEntry()

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

}
