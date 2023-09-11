package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/samoei/hotel-reservation/api"
)

func main() {
	//Get port from Command line else default to 3000
	listenAddr := flag.String("listenAddr", ":3000", "The listen address of the API server")
	flag.Parse()

	//Initiate fibre
	app := fiber.New()

	// Separate into versions
	apiv1 := app.Group("api/v1")

	//Handlers
	apiv1.Get("/users", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)

	//Start the server
	err := app.Listen(*listenAddr)
	if err != nil {
		log.Fatal(err)
	}
}
