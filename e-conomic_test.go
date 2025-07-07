package economic

import (
	"log"
	"testing"
)

func TestEconomicCustomer(t *testing.T) {
	c := &Customer{
		Address: "Testvej 1",
		City:    "Testby",
		Name:    "Test Testesen",
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
	}
	client := getTestClient()
	c, err := client.CreateCustomer(c, nil)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	otherC, err := client.GetCustomerByNumber(c.ID())
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if c.Address != otherC.Address {
		t.Fatalf("Expected %s, got %s", c.Address, otherC.Address)
	}
	if c.City != otherC.City {
		t.Fatalf("Expected %s, got %s", c.City, otherC.City)
	}
	if c.Name != otherC.Name {
		t.Fatalf("Expected %s, got %s", c.Name, otherC.Name)
	}
	if c.Zip != otherC.Zip {
		t.Fatalf("Expected %s, got %s", c.Zip, otherC.Zip)
	}
	err = client.DeleteCustomer(c)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}

func TestPaymentTerms(t *testing.T) {
	client := getTestClient()
	terms, err := client.GetPaymentTerms()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if len(terms) == 0 {
		t.Fatalf("No payment terms")
	}
	for _, term := range terms {
		log.Printf("%+v", term)
	}
}
