package economic

import (
	"net/http"
)

func (client *Client) GetPaymentTerms() ([]PaymentTerm, error) {
	var paymentTerms CollectionReponse[PaymentTerm]
	err := client.callRestAPI("payment-terms", http.MethodGet, nil, &paymentTerms)
	if err != nil {
		return paymentTerms.Collection, err
	}
	return paymentTerms.Collection, nil
}

// PaymentTerm represents a specific payment term on the agreement.
type PaymentTerm struct {
	PaymentTermsNumber int    `json:"paymentTermsNumber"` // A unique identifier of the payment term.
	DaysOfCredit       int    `json:"daysOfCredit"`       // The number of days before payment must be made.
	Description        string `json:"description"`        // A description of the payment term.
	Name               string `json:"name"`               // The name of the payment term.
	PaymentTermsType   string `json:"paymentTermsType"`   // The type of payment term.
	Self               string `json:"self"`               // A unique link reference to the payment term item.
}
