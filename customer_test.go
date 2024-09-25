package economic

import (
	"log"
	"testing"
)

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

	c.Create()
	defer c.Delete()
	found := FindCustomerByOrgNumber("28971958")
	if len(found) == 0 {
		t.Fatalf("Expected to find customer")
	}
	if found[0].Name != c.Name {
		t.Fatalf("Expected %s, got %s", c.Name, found[0].Name)
	}
}

func TestGetOrCreateCustomer(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	c := Customer{
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
	_, err := GetOrCreateCustomer(c, contact)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	found := FindCustomerByOrgNumber("66666666")
	if len(found) != 1 {
		t.Fatalf("Expected to find customer only one customer")
	}
	contact2 := CustomerContact{
		Name:  "Employee 2 Testesen",
		Email: "employee2@abe.com",
	}
	got, err := GetOrCreateCustomer(c, contact2)
	defer got.Delete()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	found = FindCustomerByOrgNumber("66666666")
	if len(found) != 1 {
		t.Fatalf("Expected to find customer only one customer")
	}
}
