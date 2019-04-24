package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnect(t *testing.T) {
	setEnvs()
	InitTestDBPG()
	var testConnect int
	_ = DB.Raw("select 1").Row().Scan(&testConnect)
	assert.EqualValues(t, 1, testConnect)
}
