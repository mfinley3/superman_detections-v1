package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/assert"

	"github.com/mfinley3/superman_detections-v1/internal/detections"
	detection "github.com/mfinley3/superman_detections-v1/internal/detections/service"
	loginrepository "github.com/mfinley3/superman_detections-v1/internal/detections/sqlite"
	"github.com/mfinley3/superman_detections-v1/internal/detections/transport"
)

func newService(t *testing.T) detection.Service {
	db, err := loginrepository.ConnectAndMigrateDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	geobase, err := geoip2.Open("../../../../resources/GeoLite2-City.mmdb")
	if err != nil {
		t.Fatal(err)
	}

	lr := loginrepository.New(db)
	return detection.New(lr, geobase)
}

func TestLoginEndpoint(t *testing.T) {

	mux := chi.NewRouter()
	mux.Handle("/logins", Handler(newService(t)))
	ts := httptest.NewServer(mux)

	cases := []struct {
		desc               string
		body               detections.Login
		expectedResp       interface{}
		expectedStatusCode int
	}{
		{"Valid Request", detections.Login{IP: "88.27.141.35", EventID: "1", Username: "test", Timestamp: 1514763000}, detections.Detection{Current: detections.GeoLocation{Latitude: 40.4143, Longitude: -3.7016, Radius: 10}}, 201},

		{"Invalid IP", detections.Login{IP: "not valid"}, errorResponse{Message: transport.ErrInvalidIP.Error()}, 400},
		{"Missing Event ID", detections.Login{IP: "0.0.0.0"}, errorResponse{Message: transport.ErrMissingID.Error()}, 400},
		{"Missing Username", detections.Login{IP: "0.0.0.0", EventID: "1"}, errorResponse{Message: transport.ErrMissingUsername.Error()}, 400},
		{"Missing Timestamp", detections.Login{IP: "0.0.0.0", EventID: "1", Username: "test"}, errorResponse{Message: transport.ErrMissingTimestamp.Error()}, 400},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			payload, err := json.Marshal(tc.body)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPost, ts.URL+"/logins", bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}
			resp, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			expectedResponse, err := json.Marshal(tc.expectedResp)
			if err != nil {
				t.Fatal(err)
			}
			actualResponse, err := ioutil.ReadAll(resp.Body)
			assert.JSONEq(t, string(expectedResponse), string(actualResponse))
		})
	}
}
