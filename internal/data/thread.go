package data

import (
	"time"
)

type Thread struct {
	Id int
	Uuid string
	Topic string
	UserId int
	CreatedAt time.Time
}

// Threads return threads from the DB
func Threads() (threads []Thread,err error) {
	rows, err := db.Query("Select id, uuid, topic, user_id, created_at from threads order by created_at desc")
	if err != nil {
		return
	}

	for rows.Next() {
		th := Thread{}
		if err = rows.Scan(&th.Id,
			&th.Uuid,
			&th.CreatedAt,
			); err != nil{
			return
		}
		threads = append (threads, th)
	}
	rows.Close()
	return
}