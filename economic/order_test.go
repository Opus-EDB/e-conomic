package economic

import (
	"net/http"
	"testing"
)

func TestCreateOrderDraft(t *testing.T) {
	c := Customer{
		Address: "Testvej 1",
		City:    "Testby",
		Name:    "Abe Testesen",
		Email:   "corporate@email.com",
		Zip:     "1234",
		PaymentTerms: PaymentTerms{
			PaymentTermsNumber: 10,
		},
		Currency: "DKK",
		VatZone: VatZone{
			VatZoneNumber: 1,
		},
		CustomerGroup: CustomerGroup{
			CustomerGroupNumber: 1,
		},
		CorporateIdentificationNumber: "66666666",
	}
	contact := CustomerContact{
		Name:  "Abe Testesen",
		Email: "jungle@abe.com",
	}
	c, err := GetOrCreateCustomer(c, contact)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	// defer c.Delete()
	order := &Order{
		Date:     "2023-10-01",
		Currency: "DKK",
		Layout: Layout{
			LayoutNumber: 19,
		},
		PaymentTerms: PaymentTerms{
			PaymentTermsNumber: 10,
		},
		Customer: CustomerID{CustomerNumber: c.CustomerNumber},
		Recipient: Recipient{
			Name:    c.Name,
			Address: c.Address,
			City:    c.City,
			Zip:     c.Zip,
			VatZone: VatZone{VatZoneNumber: 1},
		},
		Lines: []OrderLine{},
	}
	err = CreateInvoice(order)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetLayouts(t *testing.T) {
	resp := map[string]any{}
	err := callRestAPI("layouts", http.MethodGet, nil, &resp)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Fatalf("got %+v", resp)
}

func TestGetDrafts(t *testing.T) {
	resp := map[string]any{}
	err := callRestAPI("orders/drafts", http.MethodGet, nil, &resp)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Fatalf("got %+v", resp)
}

func TestGetProducts(t *testing.T) {
	resp := map[string]any{}
	err := callRestAPI("products", http.MethodGet, nil, &resp)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Fatalf("got %+v", resp)
}
