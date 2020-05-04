package endpoints

import (
	"net/http"

	"github.com/mfinley3/superman_detections-v1/internal/detections"
)

var (
	_ detections.Response = (*LoginResponse)(nil)
)

// LINK RESPONSE

type LoginResponse struct {
	detections.Detection
	headers map[string]string
	created bool
}

func (r LoginResponse) StatusCode() int {
	if r.created {
		return http.StatusCreated
	}
	return http.StatusOK
}

// If the endpoints need them, they can return custom headers
func (r LoginResponse) Headers() map[string]string {
	return r.headers
}

func (r LoginResponse) Body() interface{} {
	return r
}

func (r LoginResponse) Empty() bool {
	return false
}
