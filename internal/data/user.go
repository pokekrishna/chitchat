package data

import (
	"database/sql"
	"github.com/pokekrishna/chitchat/pkg/log"
	"strings"
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
	Id   int
	Uuid string

	// email is essentially from within user, but this field
	// is important to Scan back the row from DB
	Email string

	User      *User
	CreatedAt time.Time
}

//func NewSession(db *sql.DB, u *User) SessionInterface {
//	return &Session{db: db, User: u}
//}

func (a *App) FindUserByEmail(u *User) (err error) {
	row := a.DB.QueryRow(
		"SELECT id, uuid, name, email, password, created_at FROM users WHERE email=$1",
		u.Email)

	if err = row.Scan(&u.Id, &u.Uuid, &u.Name, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		return
	}
	return
}

func (a *App) DeleteAllUsers() (rowsAffected int64, err error) {
	query := "delete FROM users"
	result, err := a.DB.Exec(query)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (a *App) CreateUser(u *User) (err error) {
	if err = u.Validate(); err != nil {
		return
	}

	query := "INSERT INTO users (uuid, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING " +
		"id, uuid, name, email, password, created_at"
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(CreateUUID(), u.Name, u.Email, Encrypt(u.Password), time.Now()).Scan(
		&u.Id, &u.Uuid, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return
	}
	log.Info("user created: ", u)
	return
}

func (u *User) Validate() (err error) {
	if u.Name == "" {
		return &InvalidUser{Reason: "Empty Name"}
	}

	if u.Password == "" {
		return &InvalidUser{Reason: "Password not set"}
	}

	// simple email validation
	if u.Email == "" && !strings.Contains(u.Email, "@") {
		return &InvalidUser{Reason: "email not valid"}
	}
	return
}

func (a *App) CreateSession(s *Session) (err error) {
	if err = s.User.Validate(); err != nil {
		return
	}
	var stmt *sql.Stmt

	query := "INSERT INTO sessions (uuid, email, user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id, uuid, email, created_at"
	stmt, err = a.DB.Prepare(query)
	if err != nil {
		log.Error("Cannot prepare stmt", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(CreateUUID(),
		s.User.Email,
		s.User.Id,
		time.Now(),
	).Scan(&s.Id, &s.Uuid, &s.Email, &s.CreatedAt)
	if err != nil {
		log.Error("Cannot scan back created session", err)
		return
	} else {
		log.Info("Session created for user email", s.Email)
	}
	return
}

func (a *App) FindSessionByUuid(s *Session) (err error) {
	err = a.DB.QueryRow("select id, uuid, email, user_id, created_at from sessions where uuid=$1",
		s.Uuid).Scan(&s.Id, &s.Uuid, &s.Email, s.User.Id, &s.CreatedAt)
	if err != nil {
		return
	}
	return
}

func (a *App) DeleteSession(s *Session) (err error) {
	query := "Delete from sessions where uuid=$1"
	stmt, err := a.DB.Prepare(query)
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

func (a *App) DeleteAllSessions() (rowsAffected int64, err error) {
	query := "delete FROM sessions"
	result, err := a.DB.Exec(query)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}
