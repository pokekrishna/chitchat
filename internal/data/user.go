package data

import (
	"database/sql"
	"github.com/pokekrishna/chitchat/pkg/log"
	"strings"
	"time"
)

// TODO: Something about these following interfaces ...
// TODO: ...does not feel right. Only to mock out the db interactions...
// TODO: ...these interfaces are so specific that they serve purpose of
// TODO: ...only mocking and nothing else.
// TODO: ...Figure out a way to solve this.
type SessionInterface interface{
	FindByUuid(Uuid string) (err error)
	Delete() (err error)
	Create() (err error)
	DeleteAllSessions() (rowsAffected int64, err error)

	DB() *sql.DB
	ID() int
	Uuid() string
	Email() string
	User() *User

	SetID(int)
	SetUuid(string)
	SetEmail(string)
	SetUser(*User)

}

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type session struct {
	db   *sql.DB
	id   int
	uuid string

	// email is essentially from within user, but this field
	// is important to Scan back the row from DB
	email string

	user      *User
	createdAt time.Time
}

func NewSession(db *sql.DB, u *User) SessionInterface {
	return &session{db: db, user: u}
}

func (a *App) FindUserByEmail(u *User) (err error) {
	row := a.DB.QueryRow(
		"Select id, uuid, name, email, password, created_at FROM users where email=$1",
		u.Email)

	if err = row.Scan(&u.Id, &u.Uuid, &u.Name, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		return
	}
	return
}

func (a *App)DeleteAllUsers() (rowsAffected int64, err error){
	query := "delete FROM users"
	result, err := a.DB.Exec(query)
	if err != nil{
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

func (u *User) Validate() (err error){
	if u.Name == ""{
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

func (s *session) Create()  (err error){
	var stmt *sql.Stmt
	query := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, created_at"
	stmt, err = s.db.Prepare(query)
	if err != nil {
		log.Error("Cannot prepare stmt", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(CreateUUID(),
		s.user.Email,
		s.user.Id,
		time.Now(),
	).Scan(&s.id, &s.uuid, &s.email, &s.createdAt)
	if err != nil {
		log.Error("Cannot scan back created session", err)
		return
	} else {
		log.Info("Session created for user email", s.email)
	}
	return
}

func (s *session)FindByUuid(Uuid string) (err error) {
	err = s.db.QueryRow("select id, uuid, email, user_id, created_at from sessions where uuid=$1",
		Uuid).Scan(&s.id, &s.uuid, &s.email, s.user.Id, &s.createdAt)
	if err != nil {
		return
	}
	return
}

func (s *session) Delete() (err error) {
	query := "Delete from sessions where uuid=$1"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.uuid)
	if err != nil {
		return
	}
	return
}

func (s *session)DeleteAllSessions() (rowsAffected int64, err error){
	query := "delete FROM sessions"
	result, err := s.db.Exec(query)
	if err != nil{
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (s *session) DB() *sql.DB{
	return s.db
}

func (s *session) ID() int{
	return s.id
}

func (s *session) Uuid() string{
	return s.uuid
}

func (s *session) Email() string{
	return s.email
}

func (s *session) User() *User{
	return s.user
}

func (s *session) SetUuid(Uuid string) {
	s.uuid = Uuid
}

func (s *session) SetID(Id int) {
	s.id = Id
}

func (s *session) SetEmail(Email string) {
	s.email = Email
}

func (s *session) SetUser(u *User) {
	s.user = u
}
