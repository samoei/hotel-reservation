package db

import (
	"context"

	"github.com/samoei/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collection = "bookings"

type MongoBookingStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoBookingStore(client mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     &client,
		collection: client.Database(DBNAME).Collection(collection),
	}
}

type BookingStore interface {
	InsertBookings(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	var bookings []*types.Booking

	curr, err := s.collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	if err := curr.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) InsertBookings(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	result, err := s.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = result.InsertedID.(primitive.ObjectID)
	return booking, nil
}
