package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

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
