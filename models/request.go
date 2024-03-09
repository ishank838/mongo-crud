package models

type CreateOfferRequest struct {
	Title      string     `json:"title" bson:"title"`
	Status     string     `json:"status" bson:"status"`
	Targetting Targetting `json:"targetting" bson:"targetting"`
}

type Targetting struct {
	CountyCode []string `json:"country" bson:"country_code"`
}

type ListOffers struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
