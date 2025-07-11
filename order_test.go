package economic

import (
	"net/http"
	"testing"
)

func TestCreateOrderDraft(t *testing.T) {
	c := &Customer{
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
	client := getTestClient()
	_, err := client.GetOrCreateCustomer(c, &contact, 1)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	defer client.DeleteCustomer(c)
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
	_, err = client.CreateInvoice(order)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetLayouts(t *testing.T) {
	resp := map[string]any{}
	client := getTestClient()
	err := client.callRestAPI("layouts", http.MethodGet, nil, &resp)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Fatalf("got %+v", resp)
}

func TestGetDrafts(t *testing.T) {
	resp := map[string]any{}
	client := getTestClient()
	err := client.callRestAPI("orders/drafts", http.MethodGet, nil, &resp)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Fatalf("got %+v", resp)
}

func TestGetProducts(t *testing.T) {
	resp := map[string]any{}
	client := getTestClient()
	err := client.callRestAPI("products", http.MethodGet, nil, &resp)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Fatalf("got %+v", resp)
}

func TestGetInvoice(t *testing.T) {
	c := &Customer{
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
	client := getTestClient()
	_, err := client.GetOrCreateCustomer(c, &contact, 1)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	ref := "Test1234567" // Must be unique
	order := &Order{
		Date:     "2023-10-01",
		Currency: "DKK",
		References: &References{
			Other: ref,
		},
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
		Lines: []OrderLine{
			{
				LineNumber:   1,
				Product:      &Product{ProductNumber: "1"},
				Quantity:     1,
				UnitNetPrice: 300,
			},
		}}
	invoice, err := client.CreateInvoice(order)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	found, err := client.GetInvoiceByRef(ref)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if found.DraftInvoiceNumber != invoice.DraftInvoiceNumber {
		t.Fatalf("Expected invoice number `%d`, got `%d`", invoice.DraftInvoiceNumber, found.DraftInvoiceNumber)
	}
	invoice, err = client.BookInvoice(invoice.DraftInvoiceNumber)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	found, err = client.GetInvoiceByRef(ref)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if found.BookedInvoiceNumber != invoice.BookedInvoiceNumber {
		t.Fatalf("Expected invoice number `%d`, got `%d`", invoice.BookedInvoiceNumber, found.BookedInvoiceNumber)
	}
}

func TestGetBooked(t *testing.T) {
	client := getTestClient()
	invoices, err := client.GetBookedInvoices()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Fatalf("got %+v", invoices)

}
