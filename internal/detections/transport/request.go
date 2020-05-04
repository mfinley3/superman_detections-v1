package transport

import (
	"errors"
	"net"

	"github.com/mfinley3/superman_detections-v1/internal/detections"
)

var (
	ErrInvalidIP        = errors.New("invalid IP")
	ErrMissingID        = errors.New("missing event ID")
	ErrMissingUsername  = errors.New("missing username")
	ErrMissingTimestamp = errors.New("missing timestamp")
)

type LoginReqest struct {
	Login detections.Login
}

func (r LoginReqest) Validate() error {
	if net.ParseIP(r.Login.IP) == nil {
		return ErrInvalidIP
	}
	if r.Login.EventID == "" {
		return ErrMissingID
	}
	if r.Login.Username == "" {
		return ErrMissingUsername
	}
	if r.Login.Timestamp == 0 {
		return ErrMissingTimestamp
	}
	return nil
}
