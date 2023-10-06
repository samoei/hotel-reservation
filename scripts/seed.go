package main

import (
	"context"
	"fmt"
	"log"

	"github.com/samoei/hotel-reservation/api/db"
	"github.com/samoei/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore *db.MongoHotelStore
	userStore  *db.MongoUserStore
	roomStore  *db.MongoRoomStore
	ctx        context.Context
)

func seedHotel(name, location string, stars int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Stars:    stars,
	}

	rooms := []types.Room{
		{
			RoomNumber:    "101A",
			RoomType:      "Single",
			PricePerNight: 99.99,
		},
		{
			RoomNumber:    "101B",
			RoomType:      "Double",
			PricePerNight: 120.99,
		},
		{
			RoomNumber:    "201A",
			RoomType:      "Single",
			PricePerNight: 99.99,
		},
		{
			RoomNumber:    "201B",
			RoomType:      "Double",
			PricePerNight: 120.99,
		},
		{
			RoomNumber:    "301A",
			RoomType:      "Suite",
			PricePerNight: 169.99,
		},
		{
			RoomNumber:    "301B",
			RoomType:      "Suite",
			PricePerNight: 169.99,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)

	fmt.Println("Hotel inserted successfully ", insertedHotel)

	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal("Could not insert room!", err)
		}
		fmt.Println("Room inserted successfully ", insertedRoom)
	}
}

func seedUser(firstName, lastName, email, password string) {

	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	})

	if err != nil {
		log.Fatal("Could not create user from params", err)
	}

	userStore.InsertUser(ctx, user)

}

func main() {
	seedHotel("Hotel Lamada", "Nairobi", 3)
	seedHotel("Hotel One", "Berlin", 4)
	seedHotel("Temple Hotel", "Dublin", 2)

	seedUser("Samson", "Tallam", "samT@mail.com", "after123489")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	///delete any existing tables
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(*client)
	roomStore = db.NewMongoRoomStore(*client, hotelStore)
	userStore = db.NewMongoUserStore(*client)

}
