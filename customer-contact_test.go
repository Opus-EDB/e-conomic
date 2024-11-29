package economic

import "testing"

func TestUpdateOrCreateContact(t *testing.T) {
	client := getTestClient()
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
		CorporateIdentificationNumber: "66666667",
	}
	_, err := client.CreateCustomer(&c)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	defer client.DeleteCustomer(&c)
	contact := CustomerContact{
		Name:  "Abe Testesen",
		Email: "wrongemail@abe.com",
		Phone: "12345678",
	}
	err = client.UpdateOrCreateContact(c, contact)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	contacts, err := client.getCustomerContacts(c.CustomerNumber)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if len(contacts) != 1 {
		t.Fatalf("Expected 1 contact, got %d", len(contacts))
	}
	if contacts[0].Email != contact.Email {
		t.Fatalf("Expected %s, got %s", contact.Email, contacts[0].Email)
	}
	if contacts[0].Phone != contact.Phone {
		t.Fatalf("Expected %s, got %s", contact.Phone, contacts[0].Phone)
	}
	if contacts[0].Name != contact.Name {
		t.Fatalf("Expected %s, got %s", contact.Name, contacts[0].Name)
	}
	contact.Email = "correctEmail@abe.com"
	err = client.UpdateOrCreateContact(c, contact)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	contacts, err = client.getCustomerContacts(c.CustomerNumber)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if len(contacts) != 1 {
		t.Fatalf("Expected 1 contact, got %d", len(contacts))
	}
	if contacts[0].Email != contact.Email {
		t.Fatalf("Expected %s, got %s", contact.Email, contacts[0].Email)
	}
	if contacts[0].Phone != contact.Phone {
		t.Fatalf("Expected %s, got %s", contact.Phone, contacts[0].Phone)
	}
	if contacts[0].Name != contact.Name {
		t.Fatalf("Expected %s, got %s", contact.Name, contacts[0].Name)
	}

}
