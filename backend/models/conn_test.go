package models

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	setEnvs()
	InitTestDBPG()
	var testConnect int
	_ = DB.Raw("select 1").Row().Scan(&testConnect)
	assert.EqualValues(t, 1, testConnect)
}

func TestNow(t *testing.T) {
	spew.Dump(time.Now())
	spew.Dump(time.Now().UTC())
}
