package main

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
)

func main() {
	utils.InitLogger(true)
	app := fx.New(Module)

	app.Run()
}
