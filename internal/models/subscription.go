package models

type Subscription struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	Id          string `json:"id"`
	UserId      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type CreateSubscription struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserId      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}
