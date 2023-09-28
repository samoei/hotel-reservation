package db

import (
	"context"

	"github.com/samoei/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, bson.M, bson.M) error
	GetHotels(context.Context, bson.M) ([]*types.Hotel, error)
	GetHotelByID(context.Context, primitive.ObjectID) (*types.Hotel, error)
	UpdateHotel(ctx context.Context, filter, update bson.M) error
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client mongo.Client) *MongoHotelStore {
	hotelCollection := client.Database(DBNAME).Collection(hotelColl)
	return &MongoHotelStore{
		client:     &client,
		collection: hotelCollection,
	}
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.collection.InsertOne(ctx, hotel)

	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}
func (s *MongoHotelStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}
	return nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	var hotels []*types.Hotel

	resp, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (s MongoHotelStore) GetHotelByID(ctx context.Context, id primitive.ObjectID) (*types.Hotel, error) {
	var hotel *types.Hotel
	if err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&hotel); err != nil {
		return nil, err
	}

	return hotel, nil

}

func (s MongoHotelStore) UpdateHotel(ctx context.Context, filter, values bson.M) error {
	update := bson.D{
		{
			Key: "$set", Value: values,
		},
	}
	_, err := s.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}
	return nil
}
