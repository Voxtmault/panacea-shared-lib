package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rotisserie/eris"
	"github.com/voxtmault/panacea-shared-lib/config"
)

var (
	db     *sql.DB
	authDb *sql.DB
)

type MariaDatabaseStats struct {
	OpenConnections      int           `json:"open_connections"`
	ConnectionInUse      int           `json:"connection_in_use"`
	ConnectionIdle       int           `json:"connection_idle"`
	WaitingForConnection int           `json:"waiting_for_connection"`
	TotalWaitTime        time.Duration `json:"total_wait_time"`
}

func validateMariaDBConfig(config *config.DBConfig) error {
	if config.DBUser == "" {
		return eris.New("db username is empty")
	}
	if config.DBPassword == "" {
		return eris.New("db password is empty")
	}
	if config.DBHost == "" || config.DBPort == "" {
		return eris.New("invalid db address and or port")
	}
	if config.DBName == "" {
		return eris.New("invalid db name")
	}

	return nil
}

// InitMaria Establish a connection using the provided credentials with the mariadb service
func InitMariaDB(config *config.DBConfig, authConfig *config.DBConfig) error {
	log.Println("Opening Connection to Database")
	var err error

	// Validation
	if err := validateMariaDBConfig(config); err != nil {
		return eris.Wrap(err, "invalid Data MariaDB configuration")
	}

	dsn := mysql.Config{
		User:                 config.DBUser,
		Passwd:               config.DBPassword,
		AllowNativePasswords: config.AllowNativePasswords,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", config.DBHost, config.DBPort),
		DBName:               config.DBName,
		TLSConfig:            config.TSLConfig,
		MultiStatements:      config.MultiStatements,
		Params: map[string]string{
			"charset": "utf8",
		},
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

	// Init for Auth DB
	if err = validateMariaDBConfig(authConfig); err != nil {
		return eris.Wrap(err, "invalid Auth MariaDB configuration")
	}

	dsn = mysql.Config{
		User:                 authConfig.DBUser,
		Passwd:               authConfig.DBPassword,
		AllowNativePasswords: authConfig.AllowNativePasswords,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", authConfig.DBHost, authConfig.DBPort),
		DBName:               authConfig.DBName,
		TLSConfig:            authConfig.TSLConfig,
		MultiStatements:      authConfig.MultiStatements,
		Params: map[string]string{
			"charset": "utf8",
		},
	}

	authDb, err = sql.Open(authConfig.DBDriver, dsn.FormatDSN())
	if err != nil {
		return eris.Wrap(err, "Opening auth MySQL/MariaDB Connection")
	}

	authDb.SetMaxOpenConns(20)
	authDb.SetMaxIdleConns(5)
	authDb.SetConnMaxLifetime(time.Second * 5)

	err = authDb.Ping()
	if err != nil {
		return eris.Wrap(err, "Error verifying auth database connection")
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
//
// Under normal circumstances, this shouldn't be called by anyone other than main
func Close() error {
	if err := db.Close(); err != nil {
		return eris.Wrap(err, "Closing DB")
	} else {
		return nil
	}
}
