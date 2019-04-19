package admincli

import (
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"testing"
)

var app *cli.App

func TestCli(t *testing.T) {

	assert.Nil(t, app.Run([]string{"hydro-dex-cli", "market"}))
	assert.Nil(t, app.Run([]string{"hydro-dex-cli", "address"}))
	assert.Nil(t, app.Run([]string{"hydro-dex-cli", "engine"}))
	assert.Nil(t, app.Run([]string{"hydro-dex-cli", "status"}))
}

func TestMarket(t *testing.T) {
	preTest()

}

func preTest() {
	app = NewDexCli()
}
