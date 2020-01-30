package http

import (
	"context"
	"darkstalker/pkg/services"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

var customErrorEncoder = httptransport.ServerErrorEncoder(encodeError)

func MakeHTTPHandler(apiVersion string, s services.Service) http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/" + apiVersion).Subrouter()
	{
		stalk := api.Path("/stalk/{id}")
		{
			stalk.Methods("GET").Handler(MakeStalkUserHTTPHandler(s))
		}
	}

	return r
}

func MakeStalkUserHTTPHandler(s services.Service) http.Handler {
	return httptransport.NewServer(
		MakeStalkUserEndpoint(s),
		decodeStalkUserRequest,
		encodeStalkUserResponse,
		customErrorEncoder)
}

func decodeStalkUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	return stalkUserRequest{vars["id"]}, nil
}

func encodeStalkUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(stalkUserResponse)

	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(res.Dossier)
}

// Encode errors from business logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
