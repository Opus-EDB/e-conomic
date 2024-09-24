package economic

import (
	"fmt"
	"log"
	"net/http"
)

func CreateInvoice(order *Order) (err error) {
	resp := map[string]any{}
	err = callRestAPI("invoices/drafts", http.MethodPost, order, &resp)
	if err != nil {
		log.Printf("ERROR: %#v", err)
		return err
	}
	log.Printf("Order created: %#v", resp)
	if resp["draftInvoiceNumber"] == nil {
		return fmt.Errorf("draftInvoiceNumber not found in response")
	}
	invoiceNo := int(resp["draftInvoiceNumber"].(float64))
	body := DraftInvoiceBody{
		DraftInvoice: DraftInvoice{
			DraftInvoiceNumber: invoiceNo,
		},
	}
	resp = map[string]any{}
	err = callRestAPI("invoices/booked", http.MethodPost, body, &resp)
	if err != nil {
		log.Printf("ERROR: %#v", err)
		return err
	}
	log.Printf("Order booked: %#v", resp)
	return
}

type DraftInvoiceBody struct {
	DraftInvoice DraftInvoice `json:"draftInvoice"`
}

type DraftInvoice struct {
	DraftInvoiceNumber int `json:"draftInvoiceNumber"`
}

func (o *Order) SetID(id int) {
	o.Soap.OrderHandle.ID = id
}

func (o Order) ID() int {
	return o.Soap.OrderHandle.ID
}

// required: date, currency, layout, paymentTerms, customer, recipient, recipient.name, recipient.vatZone
type Order struct {
	Date                 string            `json:"date"`                           //Order issue date. Format according to ISO-8601 (YYYY-MM-DD)."`
	Currency             string            `json:"currency"`                       //The ISO 4217 3-letter currency code of the order."`
	ExchangeRate         *float64          `json:"exchangeRate,omitempty"`         //The desired exchange rate between the order currency and the base currency of the agreement. The exchange rate expresses how much it will cost in base currency to buy 100 units of the order currency. If no exchange rate is supplied, the system will get the current daily rate, unless the order currency is the same as the base currency, in which case it will be set to 100."`
	DueDate              *string           `json:"dueDate,omitempty"`              //The date the order is due for payment. This property is only used if the terms of payment is of type 'duedate', in which case it is a mandatory property. Format according to ISO-8601 (YYYY-MM-DD)."`
	GrossAmount          *float64          `json:"grossAmount,omitempty"`          //The total order amount in the order currency after all taxes and discounts have been applied."`
	MarginInBaseCurrency *float64          `json:"marginInBaseCurrency,omitempty"` //The difference between the cost price of the items on the order and the sales net order amount in base currency."`
	MarginPercentage     *float64          `json:"marginPercentage,omitempty"`     //The margin expressed as a percentage. If the net order amount is less than the cost price this number will be negative."`
	NetAmount            float64           `json:"netAmount,omitempty"`            //The total order amount in the order currency before all taxes and discounts have been applied."`
	RoundingAmount       float64           `json:"roundingAmount,omitempty"`       //The total rounding error, if any, on the order in base currency."`
	VatAmount            float64           `json:"vatAmount,omitempty"`            //The total amount of VAT on the order in the order currency. This will have the same sign as net amount"`
	Layout               Layout            `json:"layout"`                         //The layout used by the order."`
	Project              *Project          `json:"project,omitempty"`              //The project the order is connected to."`
	PaymentTerms         PaymentTerms      `json:"paymentTerms"`                   //The terms of payment for the order."`
	Customer             CustomerID        `json:"customer"`                       //The customer of the order."`
	Recipient            Recipient         `json:"recipient"`                      //The actual recipient of the order. This may be the same info found on the customer (and will probably be so in most cases) but it may also be a different recipient. For instance, the customer placing the order may be ACME Headquarters, but the recipient of the order may be ACME IT."`
	DeliveryLocation     *DeliveryLocation `json:"deliveryLocation,omitempty"`     //A reference to the place of delivery for the goods on the order"`
	Delivery             *Delivery         `json:"delivery,omitempty"`             //The actual place of delivery for the goods on the order. This is usually the same place as the one referenced in the deliveryLocation property, but may be edited as required."`
	Notes                *Notes            `json:"notes,omitempty"`                //Notes on the order."`
	References           *References       `json:"references,omitempty"`           //Customer and company references related to this order."`
	Pdf                  *Pdf              `json:"pdf,omitempty"`                  //References a pdf representation of this order."`
	Lines                []OrderLine       `json:"lines"`                          //An array containing the specific order lines."`
	Soap                 *struct {
		OrderHandle struct {
			ID int `json:"id"`
		} `json:"orderHandle"`
	} `json:"soap,omitempty"`
}

type Project struct {
	ProjectNumber int    `json:"projectNumber"`  //A unique identifier of the project."`
	Self          string `json:"self,omitempty"` //A unique reference to the project resource."`
}

type CustomerID struct {
	CustomerNumber int    `json:"customerNumber"` //The customer id number. The customer id number can be either positive or negative, but it can't be zero."`
	Self           string `json:"self,omitempty"` //A unique reference to the customer resource."`
}

type Recipient struct {
	Name              string     `json:"name"`                        //The name of the actual recipient."`
	Address           string     `json:"address,omitempty"`           //The street address of the actual recipient."`
	Zip               string     `json:"zip,omitempty"`               //The zip code of the actual recipient."`
	City              string     `json:"city,omitempty"`              //The city of the actual recipient."`
	Country           string     `json:"country,omitempty"`           //The country of the actual recipient."`
	Ean               string     `json:"ean,omitempty"`               //The 'European Article Number' of the actual recipient."`
	PublicEntryNumber string     `json:"publicEntryNumber,omitempty"` //The public entry number of the actual recipient."`
	Attention         *Attention `json:"attention,omitempty"`         //The person to whom this order is addressed."`
	VatZone           VatZone    `json:"vatZone"`                     //Recipient vat zone."`
	MobilePhone       string     `json:"mobilePhone,omitempty"`       //The phone number the order was sent to (if applicable)."`
	NemHandelType     string     `json:"nemHandelType,omitempty"`     //Chosen NemHandel type used for e-invoicing."`
}

type Attention struct {
	CustomerContactNumber int    `json:"customerContactNumber,omitempty"` //Unique identifier of the customer employee."`
	Self                  string `json:"self,omitempty"`                  //A unique reference to the customer employee."`
}

type DeliveryLocation struct {
	DeliveryLocationNumber int    `json:"deliveryLocationNumber"` //A unique identifier for the delivery location."`
	Self                   string `json:"self,omitempty"`         //A unique reference to the delivery location resource."`
}

type Delivery struct {
	Address       string `json:"address,omitempty"`       //Street address where the goods must be delivered to the customer."`
	Zip           string `json:"zip,omitempty"`           //The zip code of the place of delivery."`
	City          string `json:"city,omitempty"`          //The city of the place of delivery"`
	Country       string `json:"country,omitempty"`       //The country of the place of delivery"`
	DeliveryTerms string `json:"deliveryTerms,omitempty"` //Details about the terms of delivery."`
	DeliveryDate  string `json:"deliveryDate,omitempty"`  //The date of delivery."`
}

type Notes struct {
	Heading   string `json:"heading,omitempty"`   //The order heading. Usually displayed at the top of the order."`
	TextLine1 string `json:"textLine1,omitempty"` //The first line of supplementary text on the order. This is usually displayed right under the heading in a smaller font."`
	TextLine2 string `json:"textLine2,omitempty"` //The snd line of supplementary text in the notes on the order. This is usually displayed as a footer on the order."`
}

type References struct {
	CustomerContact *CustomerContactID `json:"customerContact,omitempty"` //The customer contact is a reference to the employee at the customer to contact regarding the order."`
	SalesPerson     *SalesPerson       `json:"salesPerson,omitempty"`     //The primary sales person is a reference to the employee who sold the goods on the order."`
	VendorReference *VendorReference   `json:"vendorReference,omitempty"` //A reference to any snd employee involved in the sale."`
	Other           string             `json:"other,omitempty"`           //A text field that can be used to save any custom reference on the order."`
}

type VendorReference struct {
	EmployeeNumber int    `json:"employeeNumber"` //Unique identifier of the employee."`
	Self           string `json:"self,omitempty"` //A unique reference to the employee resource."`
}

type Pdf struct {
	Download string `json:"download,omitempty"` //The unique reference of the pdf representation for this draft order."`
}

type OrderLine struct {
	LineNumber               int                       `json:"lineNumber"`                         //The line number is a unique number within the order."`
	SortKey                  int                       `json:"sortKey,omitempty"`                  //A sort key used to sort the lines in ascending order within the order."`
	Description              string                    `json:"description,omitempty"`              //A description of the product or service sold. Please note, that when setting existing products, description field is required. While setting non-existing product, description field can remain empty."`
	Accrual                  *Accrual                  `json:"accrual,omitempty"`                  //The accrual for the order."`
	Unit                     *Unit                     `json:"unit,omitempty"`                     //The unit of measure applied to the order line."`
	Product                  *Product                  `json:"product,omitempty"`                  //The product or service offered on the order line."`
	Quantity                 float64                   `json:"quantity,omitempty"`                 //The number of units of goods on the order line."`
	UnitNetPrice             float64                   `json:"unitNetPrice,omitempty"`             //The price of 1 unit of the goods or services on the order line in the order currency."`
	DiscountPercentage       float64                   `json:"discountPercentage,omitempty"`       //A line discount expressed as a percentage."`
	UnitCostPrice            float64                   `json:"unitCostPrice,omitempty"`            //The cost price of 1 unit of the goods or services in the order currency."`
	MarginInBaseCurrency     float64                   `json:"marginInBaseCurrency,omitempty"`     //The difference between the net price and the cost price on the order line in base currency."`
	MarginPercentage         float64                   `json:"marginPercentage,omitempty"`         //The margin on the order line expressed as a percentage."`
	DepartmentalDistribution *DepartmentalDistribution `json:"departmentalDistribution,omitempty"` //A departmental distribution defines which departments this entry is distributed between. This requires the departments module to be enabled."`
}

type Accrual struct {
	StartDate string `json:"startDate,omitempty"` //The start date for the accrual. Format: YYYY-MM-DD."`
	EndDate   string `json:"endDate,omitempty"`   //The end date for the accrual. Format: YYYY-MM-DD."`
}

type Unit struct {
	UnitNumber int    `json:"unitNumber"`     //The unique identifier of the unit."`
	Self       string `json:"self,omitempty"` //A unique reference to the unit resource."`
}

type Product struct {
	ProductNumber string `json:"productNumber,omitempty"` //The unique product number. This can be a stock keeping unit identifier (SKU)."`
	Self          string `json:"self,omitempty"`          //A unique reference to the product resource."`
}

type DepartmentalDistribution struct {
	DepartmentalDistributionNumber int    `json:"departmentalDistributionNumber"` //A unique identifier of the departmental distribution."`
	DistributionType               string `json:"distributionType,omitempty"`     //Type of the distribution"`
	Self                           string `json:"self,omitempty"`                 //A unique reference to the departmental distribution resource."`
}
