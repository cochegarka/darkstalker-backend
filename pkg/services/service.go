package services

import "context"

// Service is a service intended for retrieving data about user via VK API.
type Service interface {
	stalker
}

type stalker interface {
	// Returns data about user with given id.
	StalkUser(ctx context.Context, userId string) (interface{}, error)
}
