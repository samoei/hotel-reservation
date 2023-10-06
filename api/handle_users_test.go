package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/samoei/hotel-reservation/api/db"
	"github.com/samoei/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testdburi = "mongodb://localhost:27017"

type testdb struct {
	db.Store
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(*client)
	roomStore := db.NewMongoRoomStore(*client, hotelStore)
	UserStore := db.NewMongoUserStore(*client)
	return &testdb{
		db.Store{
			Room:  roomStore,
			Hotel: hotelStore,
			User:  UserStore,
		},
	}
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.User.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	userHandler := NewUserHandler(&tdb.Store)
	app := fiber.New()
	app.Post("/", userHandler.HandleCreateUser)
	// password := "9876hunter12345"
	params := types.CreateUserParams{
		Email:     "phil@samoei.com",
		FirstName: "Phil",
		LastName:  "Samoei",
		Password:  "9876hunter12345",
	}
	b, _ := json.Marshal(params)
	body := bytes.NewReader(b)
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	//TODO: Test if the password is blank
	//      Test if the ID is set
	// if len(user.ID)  0 {
	// 	t.Errorf("Expected firstname %s but got %s", params.FirstName, user.FirstName)
	// }

	if user.FirstName != params.FirstName {
		t.Errorf("Expected firstname %s but got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("Expected lastname %s but got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("Expected lastname %s but got %s", params.Email, user.Email)
	}

	// encpw, _ := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	// if string(encpw) != user.EncryptedPassword {
	// 	t.Errorf("Expected password %s but got %s", string(encpw), user.EncryptedPassword)
	// }

}
