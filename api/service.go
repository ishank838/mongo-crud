package api

import (
	"context"
	"mongo-crud/models"
	"mongo-crud/store"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
		CreateOfferRequest: req,
		CreatedAt:          time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (o offer) ListOffer(req models.ListOffers) ([]models.OfferResponse, error) {
	offset := (req.Page - 1) * req.Limit
	opts := options.Find().SetLimit(int64(req.Limit)).SetSkip(int64(offset)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := o.storeDeps.GetMany(context.TODO(), models.CollectionOffersV2, bson.D{}, opts)
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
