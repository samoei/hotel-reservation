package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/samoei/hotel-reservation/api"
	"github.com/samoei/hotel-reservation/api/db"
	"github.com/samoei/hotel-reservation/api/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create a new fiber instance with custom config
var config = fiber.Config{
	// Override default error handler
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	//set up mongo DB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	//Get port from Command line else default to 3000
	listenAddr := flag.String("listenAddr", ":3000", "The listen address of the API server")
	flag.Parse()

	//Initiate fibre
	app := fiber.New(config)

	//Initialise handlers
	var (
		hotelStore   = db.NewMongoHotelStore(*client)
		roomStore    = db.NewMongoRoomStore(*client, hotelStore)
		userStore    = db.NewMongoUserStore(*client)
		bookingStore = db.NewMongoBookingStore(*client)
		store        = &db.Store{
			Room:    roomStore,
			Hotel:   hotelStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(store)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
	)
	//Auth
	apiBase := app.Group("api")
	apiBase.Post("/auth", authHandler.HandleAuth)

	// Separate handlers into versions
	apiv1 := app.Group("api/v1", middleware.JWTAuthentication(userStore))

	//User Handlers
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	//Hotel Handlers
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Put("/hotel/:id", hotelHandler.HandleUpdateHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	//Room handlers
	apiv1.Post("/room/:id/book", roomHandler.HandleBookingRoom)

	//Booking handlers
	apiv1.Get("/bookings", bookingHandler.HandleGetBookings)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)

	//Start the server
	err = app.Listen(*listenAddr)
	if err != nil {
		log.Fatal(err)
	}
}
