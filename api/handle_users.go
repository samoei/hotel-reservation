package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/samoei/hotel-reservation/api/db"
	"github.com/samoei/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(db db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: db,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); len(err) > 0 {
		return c.JSON(err)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	user, err := h.userStore.GetUserByID(c.Context(), c.Params("id"))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "user not found"})
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		values bson.M
		userId = c.Params("id")
	)

	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	if err := c.BodyParser(&values); err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	if err := h.userStore.UpdateUser(c.Context(), filter, values); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": userId})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	err := h.userStore.DeleteUser(c.Context(), userId)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{"deleted": userId})

}
