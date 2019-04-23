package models

import (
	"github.com/HydroProtocol/hydro-sdk-backend/test"
	"os"
	"testing"
)

func TestConnect(t *testing.T) {
	test.PreTest()

	db := Connect(os.Getenv("HSK_DATABASE_URL"))
}
