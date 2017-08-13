package db

import (
	"database/sql"
	"fmt"

	"github.com/alde/fusion/config"

	// This package takes care of all database-related activities
	_ "github.com/go-sql-driver/mysql"
)

// FusionDAO sturct is a DAO for database related actions
type FusionDAO struct {
	conf config.DatabaseConfig
}

// New Creates a new FusionDAO
func New(cfg config.DatabaseConfig) *FusionDAO {
	return &FusionDAO{conf: cfg}
}

func (f *FusionDAO) connect() (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		f.conf.User, f.conf.Password, f.conf.Host, f.conf.Port, f.conf.Name)
	return sql.Open("mysql", connectionString)
}
