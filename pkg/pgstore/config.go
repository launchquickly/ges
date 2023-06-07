package pgstore

import (
	"database/sql"
	"fmt"
)

const (
	driverName = "postgres"
)

// PasswordConfig database connection details that authenticates via password.
type PasswordConfig struct {
	Host     string `json:"host"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Username string `json:"username"`
}

// NewPasswordDB returns a database handle representing a pool of zero or more underlying connections that authenticate
// via username and password.
//
// * c      -  used to configure database connection
// * local  -  true if application running locally, false if running on AWS
func NewPasswordDB(c PasswordConfig, local bool) (*sql.DB, error) {
	dsn := createConnectionString(c.Host, c.Port, c.Username, c.Password, c.Name, local)
	DB, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return DB, nil
}

func createConnectionString(host string, port int, username, password, name string, local bool) string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, username, password, name)

	if local {
		dsn = fmt.Sprintf("%s sslmode=disable", dsn)
	} else {
		dsn = fmt.Sprintf("%s sslmode=require", dsn)
	}

	return dsn
}
