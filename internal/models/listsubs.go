package models

type ListSubscriptionsResponse struct {
	Subscriptions []*Subscription `json:"subscriptions"`
}
