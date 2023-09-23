package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Stars    int                  `bson:"stars" json:"stars"`
}

type Room struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomType      string             `bson:"room_type" json:"room_type"`
	RoomNumber    string             `bson:"room_number" json:"room_number"`
	PricePerNight float64            `bson:"price_per_night" json:"price_per_night"`
	HotelID       primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
