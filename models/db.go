package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OfferDbModel struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Title      string             `json:"title" bson:"title"`
	Status     string             `json:"status" bson:"status"`
	Targetting Targetting         `json:"targetting" bson:"targetting"`
	CreatedAt  time.Time          `bson:"created_at"`
}
