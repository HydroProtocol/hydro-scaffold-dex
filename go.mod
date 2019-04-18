module github.com/HydroProtocol/hydro-box-dex

go 1.12

require (
	github.com/HydroProtocol/hydro-sdk-backend v0.0.0-20190416090304-3de52adf566c
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-redis/redis v6.15.1+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.2.8
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/stretchr/testify v1.3.0
	github.com/urfave/cli v1.20.0
	github.com/valyala/fasttemplate v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20190411191339-88737f569e3a // indirect
	google.golang.org/appengine v1.5.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.28.0
)

// for local test only
replace github.com/HydroProtocol/hydro-sdk-backend => ../hydro-sdk-backend
