package data

import (
	"database/sql"
	"time"
)

type ThreadInterface interface{
	FetchAll() ([]ThreadInterface, error)
	DB() *sql.DB
}

type thread struct {
	db        *sql.DB
	id        int
	uuid      string
	topic     string
	userId    int
	createdAt time.Time
}

func NewThread(db *sql.DB) ThreadInterface {
	return &thread{
		db: db,
	}
}

// FetchAll return threads from the DB
func (t *thread) FetchAll() (threads []ThreadInterface, err error) {
	if t.db == nil{
		return nil, &InvalidDBConn{Reason: "db nil"}
	}
	rows, err := t.db.Query("Select id, uuid, topic, user_id, created_at from threads order by created_at desc")
	if err != nil {
		return
	}

	for rows.Next() {
		th := &thread{}
		if err = rows.Scan(th.id,
			th.uuid,
			th.createdAt,
			); err != nil{
			return
		}
		threads = append (threads, th)
	}
	rows.Close()
	return
}

func (t *thread) DB() *sql.DB{
	return t.db
}