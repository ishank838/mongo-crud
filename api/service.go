package api

import (
	"context"
	"mongo-crud/models"
	"mongo-crud/store"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Offer interface {
	CreateOffer(models.CreateOfferRequest) error
	ListOffer(models.ListOffers) ([]models.OfferResponse, error)
}

type offer struct {
	storeDeps store.MongoStore
}

func NewOfferSvc(storeDeps store.MongoStore) Offer {
	return offer{
		storeDeps: storeDeps,
	}
}

func (o offer) CreateOffer(req models.CreateOfferRequest) error {
	_, err := o.storeDeps.Insert(context.TODO(), models.CollectionOffersV2, models.OfferDbModel{
		ID:          primitive.NewObjectID(),
		Title:       req.Title,
		Status:      req.Status,
		Targetting:  req.Targetting,
		CreatedAt:   time.Now(),
		CreatedDate: primitive.NewDateTimeFromTime(time.Now()),
	})
	if err != nil {
		return err
	}
	return nil
}

func (o offer) ListOffer(req models.ListOffers) ([]models.OfferResponse, error) {
	offset := int64((req.Page - 1) * req.Limit)
	limit := int64(req.Limit)
	opts := options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	}

	filter := bson.M{}

	//figure out how to have dynamic filters
	// if val, ok := req.Filters["status"]; ok {
	// 	filter["status"] = val
	// }

	cursor, err := o.storeDeps.GetMany(context.TODO(), models.CollectionOffersV2, filter, &opts)
	if err != nil {
		return nil, err
	}

	offers := []models.OfferResponse{}
	for cursor.Next(context.TODO()) {
		var offer models.OfferDbModel
		err = cursor.Decode(&offer)
		if err != nil {
			return nil, err
		}
		offers = append(offers, models.OfferResponse{OfferDbModel: offer})
	}

	return offers, nil
}
