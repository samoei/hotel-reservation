package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/samoei/hotel-reservation/api/db"
	"github.com/samoei/hotel-reservation/types"
)

func insertTestUser(t *testing.T, store db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     "testuser@mail.com",
		FirstName: "Firstname",
		LastName:  "Lastname",
		Password:  "superhardtestpassword",
	})

	if err != nil {
		t.Fatal(err)
	}

	insertedUser, err := store.InsertUser(context.Background(), user)

	if err != nil {
		t.Fatal(err)
	}

	return insertedUser
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertedUser := insertTestUser(t, tdb.User)

	authHandler := NewAuthHandler(&tdb.Store)
	app := fiber.New()
	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "testuser@mail.com",
		Password: "superhardtestpassword",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	//Test if the token was returned
	if authResp.Token == "" {
		t.Fatal("Expected a token to be returned but nil was returned")
	}

	//Test we are setting the token for the right user
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, &authResp.User) {
		fmt.Println(insertedUser)
		fmt.Println(authResp.User)
		t.Fatal("Expected the user to be the same user inserted")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertTestUser(t, tdb.User)

	authHandler := NewAuthHandler(&tdb.Store)
	app := fiber.New()
	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "testuser@mail.com",
		Password: "notcorrectpassword",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400 but got %d", resp.StatusCode)
	}

	var genericResp GenericResponse

	if err := json.NewDecoder(resp.Body).Decode(&genericResp); err != nil {
		t.Fatal(err)
	}

	if genericResp.Type != "error" {
		t.Fatalf("Expected generic response of type error but got %v", genericResp.Type)
	}

	if genericResp.Msg != "invalid credentials" {
		t.Fatalf("Expected generic response message of <invalid credentials> but got %v", genericResp.Msg)
	}

}
