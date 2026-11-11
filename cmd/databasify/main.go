package main

import (
	"os"

	sdk "github.com/databasifylabs/engine-sdk"
	"github.com/databasifylabs/engine-template/engine"
)

func main() {
	os.Exit(sdk.Run(engine.New()))
}
