package main

import (
	"backend/api/http"
	"backend/config"
	"backend/internal"
	"fmt"
)

func main() {
	cfg, err := config.GenerateAppConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("config parsed to %v\n", *cfg)

	app := internal.NewApp(cfg)
	err = app.Init()
	if err != nil {
		panic(err)
	}
	fmt.Printf("app created: %v\n", app)

	httpServer := http.NewHttpServer(cfg)
	err = httpServer.Init()
	if err != nil {
		panic(err)
	}
	err = httpServer.MapHandlers(app)
	if err != nil {
		panic(err)
	}
	err = httpServer.Run()
	if err != nil {
		panic(err)
	}
}
