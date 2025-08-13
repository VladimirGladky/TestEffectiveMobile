package models

type GoodResponse struct {
	Message string `json:"message"`
}

type BadResponse struct {
	Error string `json:"error"`
}

type SumSubscriptionsResponse struct {
	Sum int `json:"sum"`
}
