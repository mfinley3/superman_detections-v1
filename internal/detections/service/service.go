package service

import (
	"context"
	"net"

	"github.com/mfinley3/superman_detections-v1/internal/detections"
	"github.com/oschwald/geoip2-golang"
	"github.com/umahmood/haversine"
)

type Service interface {
	Detect(context.Context, detections.Login) (detections.Detection, error)
}

//How we can assert that a struct completely implements an interface
var _ Service = (*detectionService)(nil)

type detectionService struct {
	logins detections.LoginRepository
	geoIP2 *geoip2.Reader
}

func New(lr detections.LoginRepository, geobase *geoip2.Reader) Service {
	return &detectionService{
		logins: lr,
		geoIP2: geobase,
	}
}

func (ds *detectionService) Detect(ctx context.Context, login detections.Login) (detections.Detection, error) {
	city, err := ds.geoIP2.City(net.ParseIP(login.IP))
	if err != nil {
		return detections.Detection{}, err
	}

	login.GeoLocation = detections.GeoLocation{
		Latitude:  city.Location.Latitude,
		Longitude: city.Location.Longitude,
		Radius:    city.Location.AccuracyRadius,
	}

	login, err = ds.logins.Save(login)
	if err != nil {
		return detections.Detection{}, err
	}

	logins, err := ds.logins.FindPreceding(login)
	if err != nil {
		return detections.Detection{}, err
	}

	r := detections.Detection{
		Current: login.GeoLocation,
	}

	l, dist := findNearestLogin(login.GeoLocation, logins)
	speed, suspicious := calculateTravelSpeed(l.Timestamp, login.Timestamp, dist)
	r.Preceding = detections.NewAccess(l, speed)
	r.IsTravelFromSuspicious = suspicious

	logins, err = ds.logins.FindSubsequent(login)
	if err != nil {
		return detections.Detection{}, err
	}

	l, dist = findNearestLogin(login.GeoLocation, logins)
	speed, suspicious = calculateTravelSpeed(login.Timestamp, l.Timestamp, dist)
	r.Subsequent = detections.NewAccess(l, speed)
	r.IsTravelToSuspicious = suspicious

	return r, nil
}

func findNearestLogin(location detections.GeoLocation, logins []detections.Login) (detections.Login, float64) {
	origin := haversine.Coord{
		Lat: location.Latitude,
		Lon: location.Longitude,
	}

	closestDistance := float64(-1)
	var closestLogin detections.Login
	for _, login := range logins {
		d, _ := haversine.Distance(origin, haversine.Coord{
			Lat: login.GeoLocation.Latitude,
			Lon: login.GeoLocation.Longitude,
		})
		if d < closestDistance || closestDistance == -1 {
			closestDistance = d
			closestLogin = login
		}
	}
	return closestLogin, closestDistance
}

//Calulates speed and returns a boolean that flags if it is suspicious (over 500)
func calculateTravelSpeed(startTime, endTime int, distance float64) (int, bool) {
	timediff := float64(endTime-startTime) / 3600
	speed := int(distance / timediff)
	return speed, speed >= 500
}
