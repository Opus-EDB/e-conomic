# Usage

The package will try to read environment variable from `ECONOMIC_AGREEMENT_GRANT_TOKEN` and `ECONOMIC_APP_SECRET_TOKEN` on initialization. Alternatively, call `InititializeConfig` with the desired config. 

## Example
```go
import (
	economic "github.com/Opus-EDB/e-conomic"
)

func init() {
	config := economic.Config{
		AgreementGrant: "MY_AGREEMENT_GRANT_TOKEN",
		AppSecretToken: "MY_SECRET_APP_TOKEN",
	}
	economic.InitializeConfig(&config)
}
```
