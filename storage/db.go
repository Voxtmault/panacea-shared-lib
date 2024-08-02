package storage

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rotisserie/eris"
)

var db *sql.DB

type MariaDatabaseStats struct {
	OpenConnections      int           `json:"open_connections"`
	ConnectionInUse      int           `json:"connection_in_use"`
	ConnectionIdle       int           `json:"connection_idle"`
	WaitingForConnection int           `json:"waiting_for_connection"`
	TotalWaitTime        time.Duration `json:"total_wait_time"`
}

// InitMaria Establish a connection using the provided credentials with the mariadb service
func InitMariaDB(dsn mysql.Config) error {
	log.Println("Opening Connection to Database")
	var err error

	// Validation
	if dsn.User == "" {
		return eris.New("db username is empty")
	}
	if dsn.Passwd == "" {
		return eris.New("db password is empty")
	}
	if dsn.Addr == "" {
		return eris.New("invalid db address and or port")
	}
	if dsn.DBName == "" {
		return eris.New("invalid db name")
	}

	db, err = sql.Open(mariaDriver, dsn.FormatDSN())
	if err != nil {
		return eris.Wrap(err, "Opening MySQL/MariaDB Connection")
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Second * 5)

	err = db.Ping()
	if err != nil {
		return eris.Wrap(err, "Error verifying database connection")
	}

	log.Println("Successfully opened database connection !")
	return nil
}

func GetDBConnection() *sql.DB {
	return db
}

// GetMariaStats
func GetDBStats() MariaDatabaseStats {
	return MariaDatabaseStats{
		OpenConnections:      db.Stats().OpenConnections,
		ConnectionInUse:      db.Stats().InUse,
		ConnectionIdle:       db.Stats().Idle,
		WaitingForConnection: int(db.Stats().WaitCount),
		TotalWaitTime:        db.Stats().WaitDuration,
	}
}

// CloseMaria will close the current database connection, only do this when exiting the program
// Under normal circumstances, this shouldn't be called by anyone other than main
func Close() error {
	err := db.Close()
	if err != nil {
		return eris.Wrap(err, "Closing DB")
	}

	return nil
}
