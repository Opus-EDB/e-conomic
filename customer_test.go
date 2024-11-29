package economic

import (
	"log"
	"net/http"
	"testing"
)

func TestWhoami(t *testing.T) {
	resp := map[string]any{}
	client := getTestClient()
	err := client.callRestAPI("/self", http.MethodGet, nil, &resp)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	log.Printf("agreementNumber: %d", int(resp["agreementNumber"].(float64)))
	log.Printf("application.AppNumber: %d", int(resp["application"].(map[string]any)["appNumber"].(float64)))
	t.Fatalf("%v", resp)
}

func TestFindCustomerByName(t *testing.T) {
	c := Customer{
		Address: "Testvej 1",
		City:    "Testby",
		Name:    "Abe's Company",
		Zip:     "1234",
		Email:   "info@abe.com",
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
		CorporateIdentificationNumber: "28971958",
	}

	client := getTestClient()
	client.CreateCustomer(&c)
	defer client.DeleteCustomer(&c)
	found := client.FindCustomerByOrgNumber("28971958")
	if len(found) == 0 {
		t.Fatalf("Expected to find customer")
	}
	if found[0].Name != c.Name {
		t.Fatalf("Expected %s, got %s", c.Name, found[0].Name)
	}
}

func TestGetOrCreateCustomer(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	c := &Customer{
		Address: "Testvej 1",
		City:    "Testby",
		Name:    "Abe Testesen",
		Email:   "corporate@abe.com",
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
	err := client.GetOrCreateCustomer(c, contact)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	found := client.FindCustomerByOrgNumber("66666666")
	if len(found) != 1 {
		t.Fatalf("Expected to find customer only one customer")
	}
	contact2 := CustomerContact{
		Name:  "Employee 2 Testesen",
		Email: "employee2@abe.com",
	}
	err = client.GetOrCreateCustomer(c, contact2)
	defer client.DeleteCustomer(c)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	found = client.FindCustomerByOrgNumber("66666666")
	if len(found) != 1 {
		t.Fatalf("Expected to find customer only one customer")
	}
}
