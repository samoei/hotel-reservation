package db

import (
	"context"

	"github.com/samoei/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomCollectionName = "rooms"

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
	UpdateRoom(context.Context, bson.M, bson.M) error
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	roomCollection := client.Database(DBNAME).Collection(roomCollectionName)
	return &MongoRoomStore{
		client:     &client,
		collection: roomCollection,
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)
	room.Booked = false
	// Update Hotel with this room id
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err := s.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	resp, err := s.collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var rooms []*types.Room

	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *MongoRoomStore) UpdateRoom(ctx context.Context, filter, update bson.M) error {
	_, err := s.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	return nil
}
