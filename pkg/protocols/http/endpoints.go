package http

import (
	"context"
	"darkstalker/pkg/services"
	"github.com/go-kit/kit/endpoint"
)

func MakeStalkUserEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(stalkUserRequest)
		dossier, err := s.StalkUser(ctx, req.UserId)
		return stalkUserResponse{dossier, err}, err
	}
}

type stalkUserRequest struct {
	UserId string `json:"user_id"`
}

type stalkUserResponse struct {
	Dossier map[string]interface{} `json:"dossier"`
	Err     error                  `json:"err,omitempty"`
}
