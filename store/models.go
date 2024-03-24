package store

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type event struct {
	DocumentData documentKey `bson:"document_key"`
	Operation    string      `bson:"operationType"`
	Document     document    `bson:"_id"`
}

type document struct {
	Data string `bson:"_data"`
}

type documentKey struct {
	ID primitive.ObjectID `bson:"_id"`
}
