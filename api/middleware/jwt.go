package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("Authenticating....")

	token, ok := c.GetReqHeaders()["X-Api-Token"]

	if !ok {
		return fmt.Errorf("Unauthorized")
	}

	claims, err := validateToken(token)

	if err != nil {
		return err
	}

	expiresFloat := claims["expires"].(float64)
	expires := int64(expiresFloat)

	//Check if the token had expired
	if time.Now().Unix() > expires {
		return fmt.Errorf("token has expired")
	}

	fmt.Println(claims)

	return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")
		fmt.Println("NEVER PRINT A SECRET:", secret)

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Unauthorized")
}
