package economic

type Layout struct {
	LayoutNumber int    `json:"layoutNumber"`   //A unique identifier of the layout."`
	Self         string `json:"self,omitempty"` //A unique reference to the layout resource."`
}

// VatZone represents a VAT zone.
type VatZone struct {
	VatZoneNumber int    `json:"vatZoneNumber"`  // The unique identifier of the VAT-zone.
	Self          string `json:"self,omitempty"` // A unique link reference to the VAT-zone item.
}

// PaymentTerms represents the default payment terms for the customer.
type PaymentTerms struct {
	PaymentTermsNumber int    `json:"paymentTermsNumber"` // The unique identifier of the payment terms.
	Self               string `json:"self,omitempty"`     // A unique link reference to the payment terms item.
}

// SalesPerson represents the employee responsible for contact with the customer.
type SalesPerson struct {
	EmployeeNumber int    `json:"employeeNumber"` // The unique identifier of the employee.
	Self           string `json:"self,omitempty"` // A unique link reference to the employee resource.
}

// This is used by for collections in the Rest API
type CollectionReponse[T any] struct {
	Collection []T        `json:"collection"`
	MetaData   any        `json:"metaData"`
	Pagination Pagination `json:"pagination"`
	Self       string     `json:"self"`
}

// This is used by for collections in the Regular API
type ItemsReponse[T any] struct {
	Items      []T        `json:"items"`
	MetaData   any        `json:"metaData"`
	Pagination Pagination `json:"pagination"`
	Self       string     `json:"self"`
}

type Pagination struct {
	FirstPage            string `json:"firstPage"`
	LastPage             string `json:"lastPage"`
	MaxPageSizeAllowed   int    `json:"maxPageSizeAllowed"`
	PageSize             int    `json:"pageSize"`
	Results              int    `json:"results"`
	ResultsWithoutFilter int    `json:"resultsWithoutFilter"`
	SkipPages            int    `json:"skipPages"`
}
