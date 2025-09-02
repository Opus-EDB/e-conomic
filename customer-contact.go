package economic

import (
	"fmt"
	"log"
	"net/http"
)

func (client *Client) getCustomerContacts(customerNumber int) ([]CustomerContact, error) {
	cc := CollectionReponse[CustomerContact]{}
	err := client.callRestAPI(fmt.Sprintf("customers/%d/contacts", customerNumber), http.MethodGet, nil, &cc)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	return cc.Collection, err
}

func (c *Customer) ID() int {
	return c.CustomerNumber
}

func (c *Customer) SetID(id int) {
	c.CustomerNumber = id
}

func (client *Client) UpdateOrCreateContact(customer Customer, contact *CustomerContact) error {
	var customerInEconomic *Customer
	fmt.Println("Update or create contact")
	customerInEconomic, err := client.GetCustomer(customer)
	if err != nil {
		return err
	}
	customer = *customerInEconomic

	if contact == nil {
		return nil
	}
	contacts, err := client.getCustomerContacts(customer.CustomerNumber)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}
	// Check if contact already exists - search by email, phone or name
	for _, c := range contacts {
		if c.Email == contact.Email || c.Phone == contact.Phone || c.Name == contact.Name {
			contact.CustomerContactNumber = c.CustomerContactNumber
			customerId := customer.CustomerNumber
			path := fmt.Sprintf("customers/%d/contacts/%d", customerId, contact.CustomerContactNumber)
			err := client.callRestAPI(path, http.MethodPut, contact, &contact)
			if err != nil {
				log.Printf("Error: %s", err)
			}
			return nil
		}
	}
	_, err = client.createCustomerContact(customer.CustomerNumber, *contact)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	return err
}

func (client *Client) createCustomerContact(customerNumber int, contact CustomerContact) (CustomerContact, error) {
	var createdContact CustomerContact
	err := client.callRestAPI(fmt.Sprintf("customers/%d/contacts", customerNumber), http.MethodPost, contact, &createdContact)
	if err != nil {
		return createdContact, err
	}
	return createdContact, err
}

func (client *Client) GetCustomerContactNumber(customerNumber int) (int, error) {
	contacts, err := client.getCustomerContacts(customerNumber)
	if err != nil {
		return 0, err
	}
	numberOfContacts := len(contacts)
	if numberOfContacts < 1 {
		return 0, fmt.Errorf("no customer contact found with customer number %d", customerNumber)
	}
	return contacts[numberOfContacts-1].CustomerContactNumber, nil // return the last added contact (number)
}

type CustomerContactID struct {
	CustomerContactNumber int    `json:"customerContactNumber"` //Unique identifier of the customer contact."`
	Self                  string `json:"self,omitempty"`        //A unique reference to the customer contact resource."`
}

// CustomerContact represents a customer contact.
type CustomerContact struct {
	CustomerContactNumber int      `json:"customerContactNumber,omitempty"` // Unique numerical identifier of the customer contact.
	Email                 string   `json:"email"`                           // Customer contact e-mail address. This is where copies of sales documents are sent.
	Name                  string   `json:"name"`                            // Customer contact name.
	Phone                 string   `json:"phone"`                           // Customer contact phone number.
	EInvoiceId            string   `json:"eInvoiceId,omitempty"`            // Electronic invoicing Id. This will appear on EAN invoices in the field <cbc:ID>. Note this is not available on UK agreements.
	Notes                 string   `json:"notes,omitempty"`                 // Any notes you need to keep on a contact person.
	EmailNotifications    []string `json:"emailNotifications,omitempty"`    // This array specifies what events the contact person should receive email notifications on. Note that limited plans only have access to invoice notifications.
	Deleted               bool     `json:"deleted,omitempty"`               // Flag indicating if the contact person is deleted.
	CustomerNumber        int      `json:"customerNumber,omitempty"`        // The customer number is a positive unique numerical identifier with a maximum of 9 digits.
	CustomerSelf          string   `json:"customerSelf,omitempty"`          // A unique reference to the customer resource.
	SortKey               int      `json:"sortKey,omitempty"`               // The customer contact number displayed in the e-conomic web interface.
	Self                  string   `json:"self,omitempty"`                  // The unique self reference of the customer contact resource.
}
