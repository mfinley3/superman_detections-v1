package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/mfinley3/superman_detections-v1/internal/detections"
	"github.com/mfinley3/superman_detections-v1/internal/detections/endpoints"
	detection "github.com/mfinley3/superman_detections-v1/internal/detections/service"
	"github.com/mfinley3/superman_detections-v1/internal/detections/transport"
)

func Handler(ds detection.Service) http.Handler {
	mux := chi.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
	}

	mux.Post("/logins", kithttp.NewServer(
		endpoints.Login(ds),
		decodeLoginRequest,
		encodeResponse,
		opts...,
	).ServeHTTP)

	return mux
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	var login detections.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		return nil, err
	}

	req := transport.LoginReqest{
		Login: login,
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("content-type", "application/json")

	r, ok := resp.(detections.Response)
	if !ok {
		return errors.New("error") //make proper error handling with encode error
	}

	for k, v := range r.Headers() {
		w.Header().Set(k, v)
	}
	w.WriteHeader(r.StatusCode())

	if r.Empty() {
		return nil
	}
	return json.NewEncoder(w).Encode(r.Body())
}

type errorResponse struct {
	Message string
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	resp := errorResponse{
		Message: err.Error(),
	}
	switch err {
	case transport.ErrInvalidIP, transport.ErrMissingID, transport.ErrMissingUsername, transport.ErrMissingTimestamp:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(resp)
}
