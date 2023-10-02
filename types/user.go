package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (u UpdateUserParams) ToBson() bson.M {
	m := bson.M{}
	if len(u.FirstName) > 0 {
		m["firstName"] = u.FirstName
	}

	if len(u.LastName) > 0 {
		m["lastName"] = u.LastName
	}
	return m
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPass" json:"-"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}

	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("First name length should be at least %d characters long", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("Last name length should be at least %d characters long", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("Password length should be at least %d characters long", minPasswordLen)
	}
	if !isValidEmail(params.Email) {
		errors["email"] = "Email is not valid"
	}

	return errors
}
func isValidEmail(email string) bool {
	// Define a regular expression pattern for a valid email address
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Use the regular expression to match the email
	return regex.MatchString(email)
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}

func IsValidPassword(encPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encPassword), []byte(password)) == nil

}
