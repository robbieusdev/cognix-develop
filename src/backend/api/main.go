package main

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
	_ "go.uber.org/fx"
)

// @title Cognix API
// @version 1.0
// @description This is Cognix Golang API Documentation

// @contact.name API Support
// @contact.url
// @contact.email

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /api
// @query.collection.format multi

func main() {
	utils.InitLogger(true)
	app := fx.New(Module) //	fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
	//	return &fxevent.ZapLogger{Logger: log}
	//})

	app.Run()
}
