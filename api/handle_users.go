package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samoei/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Phil",
		LastName:  "Samoei",
	}

	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	u := types.User{
		ID:        c.Params("id"),
		FirstName: "Phil",
		LastName:  "Samoei",
	}

	return c.JSON(u)
}
