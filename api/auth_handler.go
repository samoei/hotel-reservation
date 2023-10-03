package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samoei/hotel-reservation/api/db"
	"github.com/samoei/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	store *db.Store
}

func NewAuthHandler(store *db.Store) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  types.User `json:"user"`
	Token string     `json:"token"`
}

type GenericResponse struct {
	Type string `json:"type"`
	Msg  string `json:"message"`
}

func invalideCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(GenericResponse{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

// Handlers should only
// - Deserialise incoming request
// - Fetch data
// - Call some business logic
// - Prepare response (Serialise)
// - Send back the response to the caller
func (h *AuthHandler) HandleAuth(c *fiber.Ctx) error {
	var params AuthParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.store.User.GetUserByEmail(c.Context(), params.Email)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalideCredentials(c)
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalideCredentials(c)
	}

	resp := AuthResponse{
		User:  *user,
		Token: createTokenFromUser(user),
	}

	fmt.Println("Authenticated -> ", user.FirstName)

	return c.JSON(resp)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4)
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	tokenStr, _ := token.SignedString([]byte(secret))

	return tokenStr

}
