package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/samoei/hotel-reservation/api/db"
	"github.com/samoei/hotel-reservation/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(db db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: db,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	user, err := h.userStore.GetUserByID(context.Background(), c.Params("id"))
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Phil",
		LastName:  "Samoei",
	}

	return c.JSON(u)
}
