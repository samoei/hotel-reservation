package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/samoei/hotel-reservation/api"
	"github.com/samoei/hotel-reservation/api/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"

// Create a new fiber instance with custom config
var config = fiber.Config{
	// Override default error handler
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	//set up mongo DB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	//Get port from Command line else default to 3000
	listenAddr := flag.String("listenAddr", ":3000", "The listen address of the API server")
	flag.Parse()

	//Initiate fibre
	app := fiber.New(config)

	// Separate into versions
	apiv1 := app.Group("api/v1")

	//Handlers
	userHandler := api.NewUserHandler(db.NewMongoUserStore(*client))
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	//Start the server
	err = app.Listen(*listenAddr)
	if err != nil {
		log.Fatal(err)
	}
}
