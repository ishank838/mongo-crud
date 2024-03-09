package main

import (
	"log"
	"mongo-crud/api"
	"mongo-crud/models"
	"mongo-crud/store"

	"github.com/kelseyhightower/envconfig"
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
