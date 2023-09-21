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

func main() {

	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))

	if err != nil {
		log.Fatal(err)
	}

	///delete any existing tables
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(*client)
	roomStore := db.NewMongoRoomStore(*client, hotelStore)
	hotel := types.Hotel{
		Name:     "Hotel Rwanda",
		Location: "Kigali",
		Rooms:    []primitive.ObjectID{},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)

	fmt.Println("Hotel inserted successfully ", insertedHotel)

	if err != nil {
		log.Fatal(err)
	}

	rooms := []types.Room{
		{
			Type:      types.Single,
			BasePrice: 99.99,
		},
		{
			Type:      types.Double,
			BasePrice: 149.99,
		},
		{
			Type:      types.SeaView,
			BasePrice: 199.99,
		},
		{
			Type:      types.Deluxe,
			BasePrice: 249.99,
		},
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
