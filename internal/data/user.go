package data

import (
	"database/sql"
	"github.com/pokekrishna/chitchat/pkg/log"
	"time"
)

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Session struct {
	Id int
	Uuid string
	Email string
	UserId int
	CreatedAt time.Time
}

func UserByEmail(email string) (u *User, err error){
	u = &User{}
	row := db.QueryRow(
		"Select id, uuid, name, email, password, created_at FROM users where email=$1",
		email)

	if err = row.Scan(&u.Id, &u.Uuid, &u.Name, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		return
	}
	return
}

func (u *User) CreateSession() (s *Session, err error){
	// insert into db
	s = &Session{}
	var stmt *sql.Stmt
	query := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at"
	stmt, err = db.Prepare(query)
	if err != nil {
		log.Error("Cannot prepare stmt", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(createUUID(),
		u.Email,
		u.Id,
		time.Now(),
	).Scan(&s.Id, &s.Uuid, &s.Email, &s.UserId, &s.CreatedAt)
	if err != nil {
		log.Error("Cannot scan back created session", err)
		return
	} else {
		log.Info("Session created for user email", u.Email)
	}
	return
}

func SessionByUuid(Uuid string) (s *Session, err error) {
	s = &Session{}
	err = db.QueryRow("select id, uuid, email, user_id, created_at from sessions where uuid=$1",
		Uuid).Scan(&s.Id, &s.Uuid, &s.Email, &s.UserId, &s.CreatedAt)
	if err != nil {
		return
	}
	return
}

func (s *Session)DeleteSessionByUuid() (err error) {
	query := "Delete from sessions where uuid=$1"
	stmt, err := db.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.Uuid)
	if err != nil {
		return
	}
	return
}

