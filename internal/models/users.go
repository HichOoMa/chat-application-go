package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserSignUpForm struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email"    bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UserLoginForm struct {
	Email    string `json:"email"    bson:"email"`
	Password string `json:"password" bson:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Friend struct {
	ID   primitive.ObjectID `json:"_id"  bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type User struct {
	ID             primitive.ObjectID `json:"_id"            bson:"_id"`
	Name           string             `json:"name"           bson:"name"`
	Email          string             `json:"email"          bson:"email"`
	HashedPassword string             `json:"hashedPassword" bson:"hashedPassword"`
	PdpUrl         string             `json:"pdp_url"        bson:"pdp_url"`
	Friends        []Friend           `json:"friends"        bson:"friends"`
}

func (user *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return err
	}
	user.HashedPassword = string(hashedPassword)
	return nil
}

func (user *User) PasswordMatch(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
