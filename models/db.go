package models

import "time"

type OfferDbModel struct {
	CreateOfferRequest
	CreatedAt time.Time `bson:"created_at"`
}
