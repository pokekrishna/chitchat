package data

import (
	"database/sql"
	"time"
)

type ThreadInterface interface{
	FetchAll() ([]ThreadInterface, error)
	GetDB() *sql.DB
}

type thread struct {
	DB *sql.DB
	Id int
	Uuid string
	Topic string
	UserId int
	CreatedAt time.Time
}

func NewThread(db *sql.DB) ThreadInterface {
	return &thread{
		DB: db,
	}
}

// FetchAll return threads from the DB
func (t *thread) FetchAll() (threads []ThreadInterface, err error) {
	if t.DB == nil{
		return nil, &InvalidDBConn{Reason: "db nil"}
	}
	rows, err := t.DB.Query("Select id, uuid, topic, user_id, created_at from threads order by created_at desc")
	if err != nil {
		return
	}

	for rows.Next() {
		th := &thread{}
		if err = rows.Scan(th.Id,
			th.Uuid,
			th.CreatedAt,
			); err != nil{
			return
		}
		threads = append (threads, th)
	}
	rows.Close()
	return
}

func (t *thread) GetDB() *sql.DB{
	return t.DB
}