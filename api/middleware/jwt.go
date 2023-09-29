package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("Authenticating....")

	token, ok := c.GetReqHeaders()["X-Api-Token"]

	if !ok {
		return fmt.Errorf("Unauthorized")
	}

	fmt.Println("token", token)

	return nil
}
