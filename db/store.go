package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Connect(dbDriver, dbSource string) (*sql.DB, error) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return conn, err
	}

	err = conn.Ping()
	return conn, err
}

type Store struct {
	UserRepo      UserRepo
	AuthAuditRepo AuthAuditRepo
	SessionRepo   SessionRepo
}

func NewStore(db *sql.DB) Store {
	return Store{
		UserRepo{
			db: db,
		},
		AuthAuditRepo{
			db: db,
		},
		SessionRepo{
			db: db,
		},
	}
}
