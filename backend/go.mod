module github.com/HydroProtocol/hydro-box-dex/backend

go 1.12

require (
	github.com/HydroProtocol/hydro-sdk-backend v0.0.13
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-playground/validator v9.28.0+incompatible

	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.2.8
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/stretchr/testify v1.3.0
	github.com/urfave/cli v1.20.0
	google.golang.org/appengine v1.5.0 // indirect
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.28.0
)

//replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190419153524-e8e3143a4f4a

// replace gopkg.in/go-playground/validator.v9 => github.com/go-playground/validator v9.28.0+incompatible

// replace gopkg.in/fsnotify.v1 => github.com/fsnotify/fsnotify v1.4.7

// for local test only
replace github.com/HydroProtocol/hydro-sdk-backend => ../../hydro-sdk-backend
