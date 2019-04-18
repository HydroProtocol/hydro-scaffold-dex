package config

import (
	"os"
)

var User1 string = "0xe36ea790bc9d7ab70c55260c66d52b1eca985f84"
var User2 string = "0xe834ec434daba538cd1b9fe1582052b880bd7e63"

func Getenv(name string) string {
	return os.Getenv(name)
}
