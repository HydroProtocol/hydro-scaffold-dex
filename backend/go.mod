module github.com/HydroProtocol/hydro-box-dex/backend

go 1.12

require (
	github.com/HydroProtocol/hydro-sdk-backend v0.0.13
	github.com/davecgh/go-spew v1.1.1
	github.com/go-redis/redis v6.15.1+incompatible
	github.com/jinzhu/gorm v1.9.4
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.2.8
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/stretchr/testify v1.3.0
	github.com/urfave/cli v1.20.0
	gopkg.in/go-playground/validator.v9 v9.28.0
)

// for local test only
replace github.com/HydroProtocol/hydro-sdk-backend => ../../hydro-sdk-backend

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190419153524-e8e3143a4f4a
