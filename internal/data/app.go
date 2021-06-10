package data

import "database/sql"

// TODO: Better name than App?
type App struct {
	DB *sql.DB
}
