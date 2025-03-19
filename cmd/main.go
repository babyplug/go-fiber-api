package main

import (
	"go-fiber-api/cmd/app"
	"log"

	"github.com/go-resty/resty/v2"
	// _ "time/tzdata"
)

func main() {
	// TODO: init global & remove all log in DI.
	client := resty.New()
	application, err := app.New(client)
	if err != nil {
		log.Fatalf("initial application failed: %v", err)
	}

	if err := application.Server.Listen(":" + application.Config.Port); err != nil {
		log.Fatal(err)
	}
	log.Printf("server listening on port %s", application.Config.Port)
}
