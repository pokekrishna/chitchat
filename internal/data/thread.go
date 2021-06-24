package data

import (
	"fmt"
	"time"
)

type Thread struct {
	//db        *sql.DB
	Id        int
	Uuid      string
	Topic     string
	UserId    int
	CreatedAt time.Time
}

func (a *App) Threads() ([]*Thread, error) {
	if a.DB == nil {
		return nil, &InvalidDBConn{Reason: "db nil"}
	}
	rows, err := a.DB.Query("SELECT id, uuid, topic, user_id, created_at FROM threads order by created_at desc")
	defer rows.Close()
	if err != nil {
		fmt.Println("flag1", err)
		return nil, err
	}

	var threads []*Thread
	for rows.Next() {
		th := Thread{}
		if err = rows.Scan(&th.Id,
			&th.Uuid,
			&th.Topic,
			&th.UserId,
			&th.CreatedAt,
		); err != nil {
			return threads, err
		}
		threads = append(threads, &th)
	}

	return threads, nil
}
