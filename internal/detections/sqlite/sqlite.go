package sqlite

import (
	"os"

	"github.com/jinzhu/gorm"

	"github.com/mfinley3/superman_detections-v1/internal/detections"
)

var _ detections.LoginRepository = (*loginRepository)(nil)

type loginRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) detections.LoginRepository {
	return &loginRepository{
		db: db,
	}
}

func (lr *loginRepository) Save(login detections.Login) (detections.Login, error) {
	insertStatement := `INSERT INTO login (event_uuid, username, unix_timestamp, ip_address, latitude, longitude, radius)
						VALUES (? ,? ,?, ?, ?, ?, ?);`
	db := lr.db.Exec(insertStatement, login.EventID, login.Username, login.Timestamp, login.IP, login.GeoLocation.Latitude, login.GeoLocation.Longitude, login.GeoLocation.Radius)
	return login, db.Error
}

func (lr *loginRepository) FindPreceding(login detections.Login) ([]detections.Login, error) {
	var logins []detections.Login
	selectStatement := `SELECT * FROM login WHERE username = ? AND unix_timestamp < ? ORDER BY unix_timestamp DESC;`
	db := lr.db.Raw(selectStatement, login.Username, login.Timestamp).Scan(&logins)
	return logins, db.Error
}

func (lr *loginRepository) FindSubsequent(login detections.Login) ([]detections.Login, error) {
	var logins []detections.Login
	selectStatement := `SELECT * FROM login WHERE username = ? AND unix_timestamp > ? ORDER BY unix_timestamp DESC;`
	db := lr.db.Raw(selectStatement, login.Username, login.Timestamp).Scan(&logins)
	return logins, db.Error
}

// This should live elsewhere/not exist with proper migrations and a full RDS
func ConnectAndMigrateDB(location string) (*gorm.DB, error) {

	//Create SQLite db if it doesn't exist
	if location != ":memory:" {
		_, err := os.Stat(location)
		if err == os.ErrNotExist {
			_, err := os.Create(location)
			if err != nil {
				return nil, err
			}
		}
	}

	//Open connection
	db, err := gorm.Open("sqlite3", location)
	if err != nil {
		return db, err
	}
	db.SingularTable(true)
	db.LogMode(false)

	// "migrate"
	db = db.Exec(
		`CREATE TABLE IF NOT EXISTS login (
			event_uuid TEXT NOT NULL PRIMARY KEY, 
			username TEXT NOT NULL, 
			unix_timestamp INTEGER NOT NULL, 
			ip_address TEXT NOT NULL,
			latitude REAL NOT NULL, 
			longitude REAL NOT NULL, 
			radius NOT NULL
		)`)

	return db, db.Error
}
