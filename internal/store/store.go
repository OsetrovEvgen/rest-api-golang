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

// InitTestData ...
func (s *Store) InitTestData() error {
	if _, err := s.DB.Exec(
		`INSERT INTO projects (id, name, description)
			VALUES ('def_id', 'def_name', 'def_description')`,
	); err != nil {
		logrus.Fatal(err)
	}

	if _, err := s.DB.Exec(
		`INSERT INTO columns (id, project_id, name, position)
			VALUES ('def_id', 'def_id', 'def_name', 0)`,
	); err != nil {
		logrus.Fatal(err)
	}

	if _, err := s.DB.Exec(
		`INSERT INTO tasks (id, column_id, name, description, position)
			VALUES ('def_id', 'def_id', 'def_name', 'def_description', 0)`,
	); err != nil {
		logrus.Fatal(err)
	}

	if _, err := s.DB.Exec(
		`INSERT INTO comments (id, task_id, text)
			VALUES ('def_id', 'def_id', 'def_text')`,
	); err != nil {
		logrus.Fatal(err)
	}

	return nil
}

// Close ...
func (s *Store) Close() {
	s.DB.Close()
}
