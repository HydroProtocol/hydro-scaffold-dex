package cli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCli(t *testing.T) {
	app := NewDexCli()
	assert.Nil(t, app.Run([]string{"hydro-dex-cli", "market"}))
	assert.Nil(t, app.Run([]string{"hydro-dex-cli", "address"}))
	assert.Nil(t, app.Run([]string{"hydro-dex-cli", "engine"}))
	assert.Nil(t, app.Run([]string{"hydro-dex-cli", "status"}))
}
