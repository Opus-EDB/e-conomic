package economic

// Invoice represents a booked or draft invoice schema.
type Invoice struct {
	BookedInvoiceNumber            int     `json:"bookedInvoiceNumber"`            // A reference number for the booked invoice document.
	DraftInvoiceNumber             int     `json:"draftInvoiceNumber"`             // A reference number for the draft invoice document.
	Date                           string  `json:"date"`                           // Invoice issue date. Format according to ISO-8601 (YYYY-MM-DD).
	Currency                       string  `json:"currency"`                       // The ISO 4217 currency code of the invoice.
	ExchangeRate                   float64 `json:"exchangeRate"`                   // The exchange rate between the invoice currency and the base currency of the agreement. The exchange rate expresses how much it will cost in base currency to buy 100 units of the invoice currency.
	NetAmount                      float64 `json:"netAmount"`                      // The total invoice amount in the invoice currency before all taxes and discounts have been applied. For a credit note this amount will be negative.
	NetAmountInBaseCurrency        float64 `json:"netAmountInBaseCurrency"`        // The total invoice amount in the base currency of the agreement before all taxes and discounts have been applied. For a credit note this amount will be negative.
	GrossAmount                    float64 `json:"grossAmount"`                    // The total invoice amount in the invoice currency after all taxes and discounts have been applied. For a credit note this amount will be negative.
	GrossAmountInBaseCurrency      float64 `json:"grossAmountInBaseCurrency"`      // The total invoice amount in the base currency of the agreement after all taxes and discounts have been applied. For a credit note this amount will be negative.
	VatAmount                      float64 `json:"vatAmount"`                      // The total amount of VAT on the invoice in the invoice currency. This will have the same sign as net amount.
	RoundingAmount                 float64 `json:"roundingAmount"`                 // The total rounding error, if any, on the invoice in base currency.
	Remainder                      float64 `json:"remainder"`                      // Remaining amount to be paid.
	RemainderInBaseCurrency        float64 `json:"remainderInBaseCurrency"`        // Remaining amount to be paid in base currency.
	DueDate                        string  `json:"dueDate"`                        // The date the invoice is due for payment. Only used if the terms of payment is of type 'duedate', in which case it is mandatory. Format according to ISO-8601 (YYYY-MM-DD).
	PaymentTermsNumber             int     `json:"paymentTermsNumber"`             // A unique identifier of the payment term.
	DaysOfCredit                   int     `json:"daysOfCredit"`                   // The number of days of credit on the invoice. This field is only valid if terms of payment is not of type 'duedate'.
	PaymentTermsName               string  `json:"paymentTermsName"`               // The name of the payment terms.
	PaymentTermsType               string  `json:"paymentTermsType"`               // The type of payment term.
	PaymentTermsSelf               string  `json:"paymentTermsSelf"`               // A unique reference to the payment term resource.
	CustomerNumber                 int     `json:"customerNumber"`                 // The customer number is a positive unique numerical identifier with a maximum of 9 digits.
	CustomerSelf                   string  `json:"customerSelf"`                   // A unique reference to the customer resource.
	RecipientName                  string  `json:"recipientName"`                  // The name of the actual recipient.
	RecipientAddress               string  `json:"recipientAddress"`               // The street address of the actual recipient.
	RecipientZip                   string  `json:"recipientZip"`                   // The zip code of the actual recipient.
	RecipientCity                  string  `json:"recipientCity"`                  // The city of the actual recipient.
	RecipientCountry               string  `json:"recipientCountry"`               // The country of the actual recipient.
	RecipientEAN                   string  `json:"recipientEAN"`                   // The 'European Article Number' of the actual recipient.
	RecipientPublicEntryNumber     string  `json:"recipientPublicEntryNumber"`     // The public entry number of the actual recipient.
	AttentionCustomerContactNumber int     `json:"attentionCustomerContactNumber"` // Unique identifier of the customer employee.
	AttentionSelf                  string  `json:"attentionSelf"`                  // A unique reference to the customer employee.
	VatZoneNumber                  int     `json:"vatZoneNumber"`                  // Unique identifier of the vat zone.
	VatZoneSelf                    string  `json:"vatZoneSelf"`                    // A unique reference to the vat zone.
	RecipientCVR                   string  `json:"recipientCVR"`                   // The Corporate Identification Number of the recipient for example CVR in Denmark.
	DeliveryLocationNumber         int     `json:"deliveryLocationNumber"`         // A unique identifier for the delivery location.
	DeliveryLocationSelf           string  `json:"deliveryLocationSelf"`           // A unique reference to the delivery location resource.
	DeliveryAddress                string  `json:"deliveryAddress"`                // Street address where the goods must be delivered to the customer.
	DeliveryZip                    string  `json:"deliveryZip"`                    // The zip code of the place of delivery.
	DeliveryCity                   string  `json:"deliveryCity"`                   // The city of the place of delivery.
	DeliveryCountry                string  `json:"deliveryCountry"`                // The country of the place of delivery.
	DeliveryTerms                  string  `json:"deliveryTerms"`                  // Details about the terms of delivery.
	DeliveryDate                   string  `json:"deliveryDate"`                   // The date of delivery.
	NotesHeading                   string  `json:"notesHeading"`                   // The invoice heading. Usually displayed at the top of the invoice.
	NotesTextLine1                 string  `json:"notesTextLine1"`                 // The first line of supplementary text on the invoice. This is usually displayed right under the heading in a smaller font.
	NotesTextLine2                 string  `json:"notesTextLine2"`                 // The second line of supplementary text in the notes on the invoice. This is usually displayed as a footer on the invoice.
	Customer Customer `json:"customer"`
	CustomerContactNumber          int     `json:"customerContactNumber"`          // Unique identifier of the customer contact.
	CustomerContactSelf            string  `json:"customerContactSelf"`            // A unique reference to the customer contact resource.
	SalesPersonEmployeeNumber      int     `json:"salesPersonEmployeeNumber"`      // Unique identifier of the employee.
	SalesPersonSelf                string  `json:"salesPersonSelf"`                // A unique reference to the employee resource.
	VendorReferenceEmployeeNumber  int     `json:"vendorReferenceEmployeeNumber"`  // Unique identifier of the employee.
	VendorReferenceSelf            string  `json:"vendorReferenceSelf"`            // A unique reference to the employee resource.
	OtherReference                 string  `json:"otherReference"`                 // A text field that can be used to save any custom reference on the invoice.
	PdfDownload                    string  `json:"pdfDownload"`                    // The unique reference of the pdf representation for this booked invoice.
	LayoutNumber                   int     `json:"layoutNumber"`                   // The unique identifier of the layout.
	LayoutSelf                     string  `json:"layoutSelf"`                     // A unique link reference to the layout item.
	ProjectNumber                  int     `json:"projectNumber"`                  // A unique identifier of the project.
	ProjectSelf                    string  `json:"projectSelf"`                    // A unique reference to the project resource.
	LineNumber                     int     `json:"lineNumber"`                     // The line number is a unique number within the invoice.
	SortKey                        int     `json:"sortKey"`                        // A sort key used to sort the lines in ascending order within the invoice.
	LineDescription                string  `json:"lineDescription"`                // A description of the product or service sold.
	LineDeliveryDate               string  `json:"lineDeliveryDate"`               // Invoice delivery date. The date is formatted according to ISO-8601.
	Quantity                       float64 `json:"quantity"`                       // The number of units of goods on the invoice line.
	UnitNetPrice                   float64 `json:"unitNetPrice"`                   // The price of 1 unit of the goods or services on the invoice line in the invoice currency.
	DiscountPercentage             float64 `json:"discountPercentage"`             // A line discount expressed as a percentage.
	UnitCostPrice                  float64 `json:"unitCostPrice"`                  // The cost price of 1 unit of the goods or services in the invoice currency.
	VatRate                        float64 `json:"vatRate"`                        // The VAT rate in % used to calculate the vat amount on this line.
	LineVatAmount                  float64 `json:"lineVatAmount"`                  // The total amount of VAT on the invoice line in the invoice currency. This will have the same sign as total net amount.
	TotalNetAmount                 float64 `json:"totalNetAmount"`                 // The total invoice line amount in the invoice currency before all taxes and discounts have been applied. For a credit note this amount will be negative.
	UnitNumber                     int     `json:"unitNumber"`                     // The unique identifier of the unit.
	UnitName                       string  `json:"unitName"`                       // The name of the unit (e.g. 'kg' for weight or 'l' for volume).
	UnitSelf                       string  `json:"unitSelf"`                       // A unique reference to the unit resource.
	ProductNumber                  string  `json:"productNumber"`                  // The unique product number. This can be a stock keeping unit identifier (SKU).
	ProductSelf                    string  `json:"productSelf"`                    // A unique reference to the product resource.
	DepartmentalDistributionNumber int     `json:"departmentalDistributionNumber"` // A unique identifier of the departmental distribution.
	DepartmentalDistributionSelf   string  `json:"departmentalDistributionSelf"`   // A unique reference to the departmental distribution resource.
	Sent                           string  `json:"sent"`                           // A convenience link to see if the invoice has been sent or not.
	Self                           string  `json:"self"`                           // The unique self reference of the booked invoice.
}
