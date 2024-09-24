package learnx

import (
	"encoding/json"
	"fmt"
	"opusEDB/economics/economic"
	"time"
)

func (order *Order) HandleOrder() error {
	if order.Paid {
		// Order is already paid, no need to create it in economic
		// however, we do need to create a journal entry for the payment
		amount := 0.0
		for _, item := range order.OrderItems {
			totalPrice, err := item.TotalPrice.Float64()
			if err != nil {
				return fmt.Errorf("error converting total price to float64: %w", err)
			}
			amount += totalPrice
		}
		entr := economic.JournalEntry{
			Amount:              json.Number(fmt.Sprint(amount)),
			Currency:            "DKK",
			VoucherNumber:       order.TikkoOrderID,
			Date:                time.Now().Format(time.DateOnly),
			AccountNumber:       6724,
			ContraAccountNumber: 6730,
			EntryTypeNumber:     2,
			JournalNumber:       2,
			Text:                fmt.Sprintf("Event id: %d", order.EventID),
		}
		err := entr.CreateEntry()
		if err != nil {
			return fmt.Errorf("error creating draft entry: %w", err)
		}
		return err
	}
	customer := order.customer()
	contact := order.contact()
	c, err := economic.GetOrCreateCustomer(customer, contact)
	if err != nil {
		return fmt.Errorf("customer error: %w", err)
	}
	orderLines := order.orderLines()
	o := economic.Order{
		Date:     time.Now().Format(time.DateOnly),
		Currency: order.OrderCurrency,
		Layout: economic.Layout{
			LayoutNumber: 20,
		},
		PaymentTerms: c.PaymentTerms,
		Customer:     economic.CustomerID{CustomerNumber: c.CustomerNumber},
		Lines:        orderLines,
		Recipient: economic.Recipient{
			Name:    c.Name,
			Address: order.InvoiceAddress1 + " " + order.InvoiceAddress2,
			City:    order.InvoiceCity,
			Zip:     order.InvoiceZip,
			VatZone: economic.VatZone{VatZoneNumber: 1},
		},
	}
	return economic.CreateInvoice(&o)
}

func (order *Order) contact() economic.CustomerContact {
	return economic.CustomerContact{
		Name:  order.InvoicePerson,
		Email: order.InvoiceEmail,
		Phone: order.InvoiceTelephone,
	}
}

func (order *Order) orderLines() []economic.OrderLine {
	orderLines := []economic.OrderLine{}
	for _, item := range order.OrderItems {
		unitPrice, err := item.UnitPrice.Float64()
		if err != nil {
			panic(err)
		}
		orderLine := economic.OrderLine{
			LineNumber:  item.ProductID,
			SortKey:     item.SortKey,
			Description: item.Description,
			Product: &economic.Product{
				ProductNumber: learnXIdMap(item.ProductID),
			},
			Quantity:     float64(item.Quantity),
			UnitNetPrice: unitPrice,
			DepartmentalDistribution: &economic.DepartmentalDistribution{
				DepartmentalDistributionNumber: order.EventID,
			},
		}
		orderLines = append(orderLines, orderLine)
	}
	return orderLines
}

func (order *Order) customer() economic.Customer {
	return economic.Customer{
		Name:                          order.InvoicePerson,
		Address:                       order.InvoiceAddress1,
		City:                          order.InvoiceCity,
		Zip:                           order.InvoiceZip,
		Email:                         order.InvoiceEmail,
		TelephoneAndFaxNumber:         order.InvoiceTelephone,
		Country:                       order.InvoiceCountryCode,
		CorporateIdentificationNumber: order.InvoiceCVR,
		Currency:                      "DKK",
		CustomerGroup: economic.CustomerGroup{
			CustomerGroupNumber: 1,
		},
		PaymentTerms: economic.PaymentTerms{
			PaymentTermsNumber: 4, // default to 14 days
		},
		VatZone: economic.VatZone{
			VatZoneNumber: 1,
		},
	}
}

func learnXIdMap(id int) string {
	switch id {
	case 63:
		//InvoiceGiftcardProductId
		return "15"
	case 47:
		//InvoiceTicketProductId
		return "10"
	case 10:
		//InvoiceTicketFeeProductId
		return "11"
	case 8:
		//InvoiceAddonProductId
		return "14"
	case 58:
		//InvoiceInvoiceFeeProductId
		return "12"
	case 68:
		//InvoiceSMSFeeProductId
		return "13"
	default:
		return ""
	}
}

// The data we receive from LearX when we create an order.
type Order struct {
	DueDate              *string     `json:"due_date"`               // Due date of the order
	EventDate            string      `json:"event_date"`             // Event date
	EventID              int         `json:"event_id"`               // Event ID
	IncludeVAT           bool        `json:"include_vat"`            // Include VAT
	InvoiceAddress1      string      `json:"invoice_address_1"`      // Invoice address line 1
	InvoiceAddress2      string      `json:"invoice_address_2"`      // Invoice address line 2
	InvoiceCity          string      `json:"invoice_city"`           // Invoice city
	InvoiceCompany       string      `json:"invoice_company"`        // Invoice company name
	InvoiceCountryCode   string      `json:"invoice_country_code"`   // Invoice country code
	InvoiceCVR           string      `json:"invoice_cvr"`            // Invoice CVR number
	InvoiceEAN           string      `json:"invoice_ean_ref"`        // Invoice EAN number
	InvoiceEmail         string      `json:"invoice_email"`          // Invoice email
	InvoicePerson        string      `json:"invoice_person"`         // Invoice person name
	InvoiceTelephone     string      `json:"invoice_telephone"`      // Invoice telephone number
	InvoiceZip           string      `json:"invoice_zip"`            // Invoice ZIP code
	OrderCreatedDatetime string      `json:"order_created_datetime"` // Order creation datetime
	OrderCreator         string      `json:"order_creator"`          // Order creator
	OrderCurrency        string      `json:"order_currency"`         // Order currency
	OrderDescription     string      `json:"order_description"`      // Order description
	OrderItems           []OrderItem `json:"order_items"`            // List of order items
	Paid                 bool        `json:"paid"`                   // Payment status
	SalesPersonEmail     *string     `json:"sales_person_email"`     // Sales person email
	TikkoOrderID         int         `json:"tikko_order_id"`         // Tikko order ID
}

type OrderItem struct {
	Description string       `json:"description"` // Item description
	ProductID   int          `json:"product_id"`  // Product ID
	Quantity    int          `json:"quantity"`    // Quantity of the item
	SortKey     int          `json:"sort_key"`    // Sort key
	TotalPrice  json.Number  `json:"total_price"` // Total price of the item
	UnitPrice   json.Number  `json:"unit_price"`  // Unit price of the item
	VATAmount   *json.Number `json:"vat_amount"`  // VAT amount
	VATPercent  *json.Number `json:"vat_percent"` // VAT percentage
}
