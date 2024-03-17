package main

import (
	"context"
	"fmt"
	"log"
	"mongo-crud/api"
	"mongo-crud/models"
	"mongo-crud/store"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type config struct {
	DbConnection string `required:"true"`
	Database     string `default:"offers"`
}

func main() {
	var config config
	err := envconfig.Process("ENV", &config)
	if err != nil {
		panic(err)
	}

	//Create Store
	storeDeps, err := store.NewMongoStore(config.DbConnection, config.Database)
	if err != nil {
		panic(err)
	}

	//Init Collection
	storeDeps.InitCollection(models.CollectionOffersV2)

	err = storeDeps.ExecTxn(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		result, err := storeDeps.Insert(ctx, models.CollectionOffersV2, models.OfferDbModel{
			ID:     primitive.NewObjectID(),
			Title:  "transaction testing failed",
			Status: "ACTIVE",
			Targetting: models.Targetting{
				CountyCode: []string{"IN"},
			},
			CreatedAt: time.Now(),
		})
		if err != nil {
			return nil, err
		}

		//log.Println(result.InsertedID)

		//return nil, errors.New("err")

		offer, err := storeDeps.GetMany(ctx, models.CollectionOffersV2, bson.D{{
			Key:   "_id",
			Value: result.InsertedID,
		}})
		if err != nil {
			return nil, err
		}

		var res models.OfferDbModel
		offer.Next(ctx)
		offer.Decode(&res)
		defer offer.Close(ctx)

		log.Println(fmt.Sprintf("%+v", res))

		return nil, nil
	})
	if err != nil {
		log.Println(err)
	}

	//return

	svc := api.NewOfferSvc(storeDeps)

	// err = svc.CreateOffer(models.CreateOfferRequest{
	// 	Title:  "Test Offer",
	// 	Status: "INACTIVE",
	// 	Targetting: models.Targetting{
	// 		CountyCode: []string{"US", "IN"},
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	resp, err := svc.ListOffer(models.ListOffers{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		log.Println(err)
	}

	for _, offer := range resp {
		log.Println(offer)
	}
}
