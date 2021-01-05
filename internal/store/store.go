package store

import (
	"database/sql"

	_ "github.com/lib/pq" // ...
	"github.com/sirupsen/logrus"
)

// Store ...
type Store struct {
	Config *Config
	DB     *sql.DB
}

// NewStore ...
func NewStore(config *Config) *Store {
	return &Store{
		Config: config,
	}
}

// Open ...
func (s *Store) Open() error {
	db, err := sql.Open(
		"postgres",
		"host="+s.Config.DBHost+
			" user="+s.Config.DBUser+
			" password="+s.Config.DBPassword+
			" dbname="+s.Config.DBName+
			" sslmode="+s.Config.DBMode,
	)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.DB = db
	logrus.Info("successfuly connected to database")

	return nil
}

// Close ...
func (s *Store) Close() {
	s.DB.Close()
}
