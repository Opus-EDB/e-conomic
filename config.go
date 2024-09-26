package economic

import (
	"encoding/json"
	"log"
	"os"
)

var c *Config

func init() {
	conf := os.Getenv("ECONOMIC_CONFIG_FILE")
	if conf == "" {
		agt := os.Getenv("ECONOMIC_AGREEMENT_GRANT_TOKEN")
		ast := os.Getenv("ECONOMIC_APP_SECRET_TOKEN")
		if len(agt) == 0 || len(ast) == 0 {
			fmt.Printf("WARNING: ECONOMIC_CONFIG_FILE or ECONOMIC_AGREEMENT_GRANT_TOKEN and ECONOMIC_APP_SECRET_TOKEN must be set")
		}
		log.Printf("Read config from env: Grant %sXXXXXX, App %sXXXXXX", agt[:4], ast[:4])
		c = &Config{
			AgreementGrant: agt,
			AppSecretToken: ast,
		}
		return
	}
	config, err := getConfigFromFile(conf)
	if err != nil {
		panic(err)
	}
	log.Printf("Read config from file: Grant %sXXXXXX, App %sXXXXXX", config.AgreementGrant[:4], config.AppSecretToken[:4])
	c = &config
}

func getConfigFromFile(path string) (Config, error) {
	// Read the file and return the config
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	// Unmarshal the content into a Config struct
	var c Config
	err = json.Unmarshal(content, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

type Config struct {
	AgreementGrant string `json:"agreement_grant"`
	AppSecretToken string `json:"app_secret"`
}
