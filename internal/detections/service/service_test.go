package service

import (
	"context"
	"errors"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/assert"

	"github.com/mfinley3/superman_detections-v1/internal/detections"
	loginrepository "github.com/mfinley3/superman_detections-v1/internal/detections/sqlite"
)

func newService(t *testing.T) Service {
	db, err := loginrepository.ConnectAndMigrateDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	geobase, err := geoip2.Open("../../../resources/GeoLite2-City.mmdb")
	if err != nil {
		t.Fatal(err)
	}

	lr := loginrepository.New(db)
	return New(lr, geobase)
}

func TestDetect(t *testing.T) {
	service := newService(t)

	cases := []struct {
		desc         string
		login        detections.Login
		expectedResp detections.Detection
		expectedErr  error
	}{
		{
			desc: "First call - no Subsequent or Preceding",
			login: detections.Login{
				Username:  "mike",
				Timestamp: 1514763000,
				EventID:   "a49adbf7-1874-4a8e-b61f-7d972d1d2fe4",
				IP:        "146.140.67.239",
			},
			expectedResp: detections.Detection{
				Current: detections.GeoLocation{
					Latitude: 51.2993, Longitude: 9.491, Radius: 0xc8,
				},
			},
			expectedErr: nil,
		},

		{
			desc: "Second Call - No Subsequent",
			login: detections.Login{
				Username:  "mike",
				Timestamp: 1514765000,
				EventID:   "df2dd38a-8818-46fc-a966-d826cbed6734",
				IP:        "68.195.164.188",
			},
			expectedResp: detections.Detection{
				Current: detections.GeoLocation{
					Latitude: 40.8777, Longitude: -73.908, Radius: 0x1,
				},
				IsTravelFromSuspicious: true,
				Preceding: detections.Access{
					IP:        "146.140.67.239",
					Speed:     6912,
					Timestamp: 1514763000,
					GeoLocation: detections.GeoLocation{
						Latitude: 51.2993, Longitude: 9.491, Radius: 0xc8,
					},
				},
			},
			expectedErr: nil,
		},

		{
			desc: "Third call - Complete",
			login: detections.Login{
				Username:  "mike",
				Timestamp: 1514764000,
				EventID:   "7116b633-50d5-45c0-8376-491264600d7b",
				IP:        "145.139.67.30",
			},
			expectedResp: detections.Detection{
				Current: detections.GeoLocation{
					Latitude: 52.3824, Longitude: 4.8995, Radius: 0x64,
				},
				IsTravelToSuspicious:   true,
				IsTravelFromSuspicious: true,
				Preceding: detections.Access{
					IP:        "146.140.67.239",
					Speed:     754,
					Timestamp: 1514763000,
					GeoLocation: detections.GeoLocation{
						Latitude: 51.2993, Longitude: 9.491, Radius: 0xc8,
					},
				},
				Subsequent: detections.Access{
					IP:        "68.195.164.188",
					Speed:     13070,
					Timestamp: 1514765000,
					GeoLocation: detections.GeoLocation{
						Latitude: 40.8777, Longitude: -73.908, Radius: 0x1,
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "Bad IP address",
			login: detections.Login{
				Username:  "mike",
				Timestamp: 1514763000,
				EventID:   "c8b9f233-dd0c-48c9-81d4-14e712503403",
				IP:        "Not An IP",
			},
			expectedErr: errors.New("IP passed to Lookup cannot be nil"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			resp, err := service.Detect(context.Background(), tc.login)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedResp, resp)
		})
	}
}

func TestFindNearestLogin(t *testing.T) {

	testLoginOne := detections.Login{
		GeoLocation: detections.GeoLocation{
			Latitude:  30.3764,
			Longitude: -97.7078,
		},
	}

	testLoginTwo := detections.Login{
		GeoLocation: detections.GeoLocation{
			Latitude:  34.0494,
			Longitude: -118.2641,
		},
	}

	cases := []struct {
		desc             string
		origin           detections.GeoLocation
		logins           []detections.Login
		expectedLogin    detections.Login
		expectedDistance float64
	}{
		{"Find Nearest Location", detections.GeoLocation{Latitude: 32, Longitude: -108.5}, []detections.Login{testLoginOne, testLoginTwo}, testLoginTwo, 582.716595358608},
		{"Find Nearest Lcation", detections.GeoLocation{Latitude: 30, Longitude: -100}, []detections.Login{testLoginOne, testLoginTwo}, testLoginOne, 139.3155860975571},

		{"Same location - testLoginOne", testLoginOne.GeoLocation, []detections.Login{testLoginOne, testLoginTwo}, testLoginOne, 0},
		{"Same location - testLoginTwo", testLoginTwo.GeoLocation, []detections.Login{testLoginOne, testLoginTwo}, testLoginTwo, 0},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			login, distance := findNearestLogin(tc.origin, tc.logins)
			assert.Equal(t, tc.expectedLogin, login)
			assert.Equal(t, tc.expectedDistance, distance)
		})
	}
}

func TestCalculateTravelSpeed(t *testing.T) {
	cases := []struct {
		desc             string
		startTime        int
		endTime          int
		distance         float64
		expectedSpeed    int
		expectSuspicious bool
	}{
		{"480 miles in 24 hours", 1514764800, 1514851200, 480, 20, false},
		{"14400 miles in 24 hours", 1514764800, 1514851200, 14400, 600, true},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			speed, suspicious := calculateTravelSpeed(tc.startTime, tc.endTime, tc.distance)
			assert.Equal(t, tc.expectedSpeed, speed)
			assert.Equal(t, tc.expectSuspicious, suspicious)
		})
	}
}
