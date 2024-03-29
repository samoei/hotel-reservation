package db

import (
	"context"
	"fmt"

	"github.com/samoei/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper
	GetUserByID(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(ctx context.Context, filter, update bson.M) error
}

type MongoUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserStore(client mongo.Client) *MongoUserStore {
	usersCollection := client.Database(DBNAME).Collection(userColl)
	return &MongoUserStore{
		client:     &client,
		collection: usersCollection,
	}
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("... Dropping users collection")
	return s.collection.Drop(ctx)
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {

	// oid, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return nil, err
	// }

	var user types.User
	if err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {

	var users []*types.User

	cur, err := s.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
func (s *MongoUserStore) UpdateUser(ctx context.Context, filter, values bson.M) error {

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
func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return err
	}

	//TODO: Might need to log if no record was deleted

	return nil
}
