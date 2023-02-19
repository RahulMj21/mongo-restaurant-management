package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	FirstName    *string            `json:"first_name" validate:"required,min=3,max=15"`
	LastName     *string            `json:"last_name" validate:"required,min=3,max=15"`
	Password     *string            `json:"password" validate:"required,min=8,max=40"`
	Email        *string            `json:"email" validate:"email,required,unique"`
	Avatar       *string            `json:"avatar"`
	AccessToken  *string            `json:"access_token"`
	RefreshToken *string            `json:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	UserId       string             `json:"user_id"`
}
