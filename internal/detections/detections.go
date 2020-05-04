package detections

type Login struct {
	Username    string `json:"username"`
	Timestamp   int    `json:"unix_timestamp" gorm:"column:unix_timestamp"`
	EventID     string `json:"event_uuid" gorm:"column:event_uuid"`
	IP          string `json:"ip_address" gorm:"column:ip_address"`
	GeoLocation `json:"-"`
}

type GeoLocation struct {
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"lon,omitempty"`
	Radius    uint16  `json:"radius,omitempty"`
}

type Access struct {
	IP        string `json:"ip_address,omitempty"`
	Speed     int    `json:"speed,omitempty"`
	Timestamp int    `json:"timestamp,omitempty"`
	GeoLocation
}

type Detection struct {
	Current                GeoLocation `json:"currentGeo,omitempty"`
	IsTravelToSuspicious   bool        `json:"traveltoCurrentGeoSuspicious"`
	IsTravelFromSuspicious bool        `json:"travelfromCurrentGeoSuspicious"`
	Preceding              Access      `json:"precedingIPAcces"`
	Subsequent             Access      `json:"subsequentIpAccess"`
}

func NewAccess(login Login, speed int) Access {
	return Access{
		IP:          login.IP,
		Speed:       speed,
		Timestamp:   login.Timestamp,
		GeoLocation: login.GeoLocation,
	}
}

type LoginRepository interface {
	Save(Login) (Login, error)
	FindPreceding(Login) ([]Login, error)
	FindSubsequent(Login) ([]Login, error)
}

type Request interface {
	Validate() bool
}

type Response interface {
	StatusCode() int
	Headers() map[string]string
	Body() interface{}
	Empty() bool
}
