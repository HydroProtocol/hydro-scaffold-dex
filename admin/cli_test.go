package main

import (
	"testing"
)

func TestCli(t *testing.T) {
	process([]string{"admin", "nmt", "hello"})
}
