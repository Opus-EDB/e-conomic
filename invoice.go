package economic

// Invoice represents a booked or draft invoice schema.
type Invoice struct {
	BookedInvoiceNumber            int           `json:"bookedInvoiceNumber"`       // A reference number for the booked invoice document.
	DraftInvoiceNumber             int           `json:"draftInvoiceNumber"`        // A reference number for the draft invoice document.
	Date                           string        `json:"date"`                      // Invoice issue date. Format according to ISO-8601 (YYYY-MM-DD).
	Currency                       string        `json:"currency"`                  // The ISO 4217 currency code of the invoice.
	ExchangeRate                   float64       `json:"exchangeRate"`              // The exchange rate between the invoice currency and the base currency of the agreement. The exchange rate expresses how much it will cost in base currency to buy 100 units of the invoice currency.
	NetAmount                      float64       `json:"netAmount"`                 // The total invoice amount in the invoice currency before all taxes and discounts have been applied. For a credit note this amount will be negative.
	NetAmountInBaseCurrency        float64       `json:"netAmountInBaseCurrency"`   // The total invoice amount in the base currency of the agreement before all taxes and discounts have been applied. For a credit note this amount will be negative.
	GrossAmount                    float64       `json:"grossAmount"`               // The total invoice amount in the invoice currency after all taxes and discounts have been applied. For a credit note this amount will be negative.
	GrossAmountInBaseCurrency      float64       `json:"grossAmountInBaseCurrency"` // The total invoice amount in the base currency of the agreement after all taxes and discounts have been applied. For a credit note this amount will be negative.
	VatAmount                      float64       `json:"vatAmount"`                 // The total amount of VAT on the invoice in the invoice currency. This will have the same sign as net amount.
	RoundingAmount                 float64       `json:"roundingAmount"`            // The total rounding error, if any, on the invoice in base currency.
	Remainder                      float64       `json:"remainder"`                 // Remaining amount to be paid.
	RemainderInBaseCurrency        float64       `json:"remainderInBaseCurrency"`   // Remaining amount to be paid in base currency.
	DueDate                        string        `json:"dueDate"`                   // The date the invoice is due for payment. Only used if the terms of payment is of type 'duedate', in which case it is mandatory. Format according to ISO-8601 (YYYY-MM-DD).
	PaymentTermsNumber             int           `json:"paymentTermsNumber"`        // A unique identifier of the payment term.
	DaysOfCredit                   int           `json:"daysOfCredit"`              // The number of days of credit on the invoice. This field is only valid if terms of payment is not of type 'duedate'.
	PaymentTerms                   *PaymentTerms `json:"paymentTerms"`
	CustomerNumber                 int           `json:"customerNumber"` // The customer number is a positive unique numerical identifier with a maximum of 9 digits.
	CustomerSelf                   string        `json:"customerSelf"`   // A unique reference to the customer resource.
	Recipient                      *Recipient    `json:"recipient"`
	Delivery                       *Delivery     `json:"delivery,omitempty"`
	Notes                          *Notes        `json:"notes"`
	Customer                       *Customer     `json:"customer,omitempty"`
	CustomerContactNumber          int           `json:"customerContactNumber"`         // Unique identifier of the customer contact.
	CustomerContactSelf            string        `json:"customerContactSelf"`           // A unique reference to the customer contact resource.
	SalesPersonEmployeeNumber      int           `json:"salesPersonEmployeeNumber"`     // Unique identifier of the employee.
	SalesPersonSelf                string        `json:"salesPersonSelf"`               // A unique reference to the employee resource.
	VendorReferenceEmployeeNumber  int           `json:"vendorReferenceEmployeeNumber"` // Unique identifier of the employee.
	VendorReferenceSelf            string        `json:"vendorReferenceSelf"`           // A unique reference to the employee resource.
	References                     *References   `json:"references"`
	Pdf                            *Pdf          `json:"pdf"`
	Layout                         *Layout       `json:"layout"`
	ProjectNumber                  int           `json:"projectNumber"`                  // A unique identifier of the project.
	ProjectSelf                    string        `json:"projectSelf"`                    // A unique reference to the project resource.
	Lines                          []OrderLine   `json:"lines"`                          // The line number is a unique number within the invoice.
	UnitNumber                     int           `json:"unitNumber"`                     // The unique identifier of the unit.
	UnitName                       string        `json:"unitName"`                       // The name of the unit (e.g. 'kg' for weight or 'l' for volume).
	UnitSelf                       string        `json:"unitSelf"`                       // A unique reference to the unit resource.
	DepartmentalDistributionNumber int           `json:"departmentalDistributionNumber"` // A unique identifier of the departmental distribution.
	DepartmentalDistributionSelf   string        `json:"departmentalDistributionSelf"`   // A unique reference to the departmental distribution resource.
	Sent                           string        `json:"sent"`                           // A convenience link to see if the invoice has been sent or not.
	Self                           string        `json:"self"`                           // The unique self reference of the booked invoice.
	LastUpdated                    string        `json:"lastUpdated,omitempty"`
}
