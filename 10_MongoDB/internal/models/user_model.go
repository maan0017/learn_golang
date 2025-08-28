package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user document in MongoDB
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Gender    string             `bson:"gender" json:"gender" validate:"required,oneof=male female other"`
	Age       uint16             `bson:"age" json:"age" validate:"gte=1,lte=120"`
	Password  string             `bson:"password" json:"password" validate:"required,min=6"` // json:"-" so password won't be exposed
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
