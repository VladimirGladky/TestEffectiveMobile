package suberrors

import "errors"

var (
	ErrIdSubscriptionNotFound = errors.New("subscription id not found")
	ErrUserIdNotFound         = errors.New("user id not found")
)
