package economic

import (
	"os"
)

// put shared test code here

func getTestClient() *Client {
	return &Client{
		AgreementGrant: os.Getenv("ECONOMIC_AGREEMENT_GRANT_TOKEN"),
		AppSecretToken: os.Getenv("ECONOMIC_APP_SECRET_TOKEN"),
	}
}
