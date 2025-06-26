package economic

import (
	"fmt"
	"log"
	"net/http"
)

func (client *Client) GetCustomerByNumber(number int) (*Customer, error) {
	var customer Customer
	err := client.callRestAPI(fmt.Sprintf("customers/%d", number), http.MethodGet, nil, &customer)
	return &customer, err
}

/*
func (client *Client) GetCustomer(customer *Customer) error {
	return client.callRestAPI(fmt.Sprintf("customers/%d", customer.CustomerNumber), http.MethodGet, nil, customer)
}
*/

func (client *Client) CreateCustomer(customer *Customer, contact *CustomerContact) (*Customer, error) {
	r := Customer{}
	err := client.callRestAPI("customers", http.MethodPost, customer, &r)
	if err != nil || contact == nil {
		return &r, err
	}
	err = client.UpdateOrCreateContact(r, *contact)
	return &r, err
}

func (client *Client) UpdateCustomer(customer *Customer, contact *CustomerContact) error {
	err := client.callRestAPI(fmt.Sprintf("customers/%d", customer.CustomerNumber), http.MethodPut, customer, nil)
	if err != nil || contact == nil {
		return err
	}
	err = client.UpdateOrCreateContact(*customer, *contact)
	return err
}

func (client *Client) DeleteCustomer(customer *Customer) error {
	err := client.callRestAPI(fmt.Sprintf("customers/%d", customer.CustomerNumber), http.MethodDelete, nil, nil)
	return err
}

// GetCustomer gets a customer from economic by customer number. If the
// customer does not exist, it creates a new customer in economic using the
// provided.  `customer` is read and modified in-place.
func (client *Client) GetOrCreateCustomer(customer *Customer, contact CustomerContact) error {
	if customer.CorporateIdentificationNumber == "" && customer.VatNumber == "" {
		return fmt.Errorf("no corporate identification number or vat number provided")
	}
	if customer.CorporateIdentificationNumber != "" {
		customer.VatNumber = customer.CorporateIdentificationNumber
	}
	customers := client.FindCustomerByOrgNumber(customer.CorporateIdentificationNumber)
	if len(customers) == 0 {
		log.Printf("No customer found with org number %s - creating", customer.CorporateIdentificationNumber)
		c, err := client.CreateCustomer(customer, &contact)
		if err != nil {
			log.Printf("Error: %s", err)
			return err
		}
		customers = append(customers, *c)
	}
	if len(customers) > 1 {
		return fmt.Errorf("multiple customers found with org number %s", customer.CorporateIdentificationNumber)
	}
	*customer = customers[0]
	return client.UpdateOrCreateContact(*customer, contact)
}

// Updates or creates a company based on the corporate identification number.
func (client *Client) UpdateOrCreateCustomer(customer Customer, contact CustomerContact) error {
	customers := client.FindCustomerByOrgNumber(customer.CorporateIdentificationNumber)
	if len(customers) == 0 {
		log.Printf("No customer found with org number %s - creating", customer.CorporateIdentificationNumber)
		c, err := client.CreateCustomer(&customer, nil) // don't include contact here
		if err != nil {
			log.Printf("Error: %s", err)
			return err
		}
		customers = append(customers, *c)
	}
	if len(customers) > 1 {
		return fmt.Errorf("multiple customers found with org number %s", customer.CorporateIdentificationNumber)
	}
	customer.CustomerNumber = customers[0].CustomerNumber
	return client.UpdateCustomer(&customer, &contact)
}

func (client *Client) FindCustomerByOrgNumber(org string) []Customer {
	filter := &Filter{}
	filter.AndCondition("corporateIdentificationNumber", FilterOperatorEquals, org)
	resp := CollectionReponse[Customer]{}
	err := client.callRestAPI("customers?filter="+filter.String(), http.MethodGet, nil, &resp)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	return resp.Collection
}

// Customer represents a customer, aka. Debtor.
type Customer struct {
	// Mandatory fields
	VatZone       VatZone       `json:"vatZone"`       // Indicates in which VAT-zone the customer is located (e.g.: domestically, in Europe or elsewhere abroad).
	Name          string        `json:"name"`          // The customer name.
	Currency      string        `json:"currency"`      // Default payment currency.
	CustomerGroup CustomerGroup `json:"customerGroup"` // In order to set up a new customer, it is necessary to specify a customer group. It is useful to group a company’s customers (e.g., ‘domestic’ and ‘foreign’ customers) and to link the group members to the same account when generating reports.
	PaymentTerms  PaymentTerms  `json:"paymentTerms"`  // The default payment terms for the customer.
	// Optional fields
	Barred                        bool         `json:"barred,omitempty"`                        // Boolean indication of whether the customer is barred from invoicing.
	Address                       string       `json:"address,omitempty"`                       // Address for the customer including street and number.
	Balance                       float64      `json:"balance,omitempty"`                       // The outstanding amount for this customer.
	CorporateIdentificationNumber string       `json:"corporateIdentificationNumber,omitempty"` // Corporate Identification Number. For example CVR in Denmark.
	PNumber                       string       `json:"pNumber,omitempty"`                       // Extension of corporate identification number (CVR). Identifying separate production unit (p-nummer).
	City                          string       `json:"city,omitempty"`                          // The customer's city.
	Country                       string       `json:"country,omitempty"`                       // The customer's country.
	CreditLimit                   float64      `json:"creditLimit,omitempty"`                   // A maximum credit for this customer. Once the maximum is reached or passed in connection with an order/quotation/invoice for this customer you see a warning in e-conomic.
	CustomerNumber                int          `json:"customerNumber,omitempty"`                // The customer number is a positive unique numerical identifier with a maximum of 9 digits. If no customer number is specified a number will be supplied by the system.
	EAN                           string       `json:"ean,omitempty"`                           // European Article Number. EAN is used for invoicing the Danish public sector.
	Email                         string       `json:"email,omitempty"`                         // Customer e-mail address where e-conomic invoices should be emailed. Note: you can specify multiple email addresses in this field, separated by a space. If you need to send a copy of the invoice or write to other e-mail addresses, you can also create one or more customer contacts.
	Layout                        *Layout      `json:"layout,omitempty"`                        // Layout to be applied for invoices and other documents for this customer.
	Zip                           string       `json:"zip,omitempty"`                           // The customer's postcode.
	PublicEntryNumber             string       `json:"publicEntryNumber,omitempty"`             // The public entry number is used for electronic invoicing, to define the account invoices will be registered on at the customer.
	TelephoneAndFaxNumber         string       `json:"telephoneAndFaxNumber,omitempty"`         // The customer's telephone and/or fax number.
	MobilePhone                   string       `json:"mobilePhone,omitempty"`                   // The customer's mobile phone number.
	EInvoicingDisabledByDefault   bool         `json:"eInvoicingDisabledByDefault,omitempty"`   // Boolean indication of whether the default sending method should be email instead of e-invoice. This property is updatable only by using PATCH to /customers/:customerNumber
	VatNumber                     string       `json:"vatNumber,omitempty"`                     // The customer's value added tax identification number. This field is only available to agreements in Sweden, UK, Germany, Poland and Finland. Not to be mistaken for the danish CVR number, which is defined on the corporateIdentificationNumber property.
	Website                       string       `json:"website,omitempty"`                       // Customer website, if applicable.
	SalesPerson                   *SalesPerson `json:"salesPerson,omitempty"`                   // Reference to the employee responsible for contact with this customer.
	PriceGroup                    *PriceGroup  `json:"priceGroup,omitempty"`                    // A unique link reference to the price-group resource.
}

// CustomerGroup represents a customer group.
type CustomerGroup struct {
	CustomerGroupNumber int    `json:"customerGroupNumber"` // The unique identifier of the customer group.
	Self                string `json:"self,omitempty"`      // A unique link reference to the customer group item.
}

// PriceGroup represents a price group.
type PriceGroup struct {
	Self string `json:"self"` // A unique link reference to the price-group resource.
}
