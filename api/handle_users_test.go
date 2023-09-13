package api

import (
	"context"
	"log"
	"testing"

	"github.com/samoei/hotel-reservation/api/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testdburi = "mongodb://localhost:27017"

type testdb struct {
	db.UserStore
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(*client),
	}
}

func (tdb *testdb) teardown(t *testing.T) {

}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	t.Fail()

}
