package economic

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
