package data

import (
	"fmt"
	"time"
)

type Thread struct {
	//db        *sql.DB
	id        int
	uuid      string
	topic     string
	userId    int
	createdAt time.Time
}

func (a *App) Threads() ([]*Thread, error) {
	if a.DB == nil{
		return nil, &InvalidDBConn{Reason: "db nil"}
	}
	rows, err := a.DB.Query("SELECT id, uuid, topic, user_id, created_at FROM threads order by created_at desc")
	defer rows.Close()
	if err != nil {
		fmt.Println("flag1", err)
		return nil, err
	}

	// TODO : make the underlying array of a known length
	var threads []*Thread
	for rows.Next() {
		th := &Thread{}
		if err = rows.Scan(th.id,
			th.uuid,
			th.createdAt,
			); err != nil{
			return threads, err
		}
		threads = append (threads, th)
	}

	return threads, nil
}