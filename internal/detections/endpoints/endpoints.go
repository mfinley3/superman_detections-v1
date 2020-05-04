package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	detection "github.com/mfinley3/superman_detections-v1/internal/detections/service"
	"github.com/mfinley3/superman_detections-v1/internal/detections/transport"
)

func Login(ds detection.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.LoginReqest)
		if err := req.Validate(); err != nil {
			return nil, err
		}

		d, err := ds.Detect(ctx, req.Login)
		if err != nil {
			return nil, err
		}

		return LoginResponse{
			Detection: d,
			created:   true,
		}, nil
	}
}
