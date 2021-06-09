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
type UserInterface interface{
	FindByEmail(email string) (err error)
	DeleteAllUsers() (rowsAffected int64, err error)
	Create() (err error)
	Validate() (err error)

	DB() *sql.DB
	ID() int
	Uuid() string
	Name() string
	Email() string
	Password() string

	setID(int)
	SetUuid(string)
	SetName(string)
	SetEmail(string)
	SetPassword(string)
}

type SessionInterface interface{
	FindByUuid(Uuid string) (err error)
	Delete() (err error)
	Create() (err error)
	DeleteAllSessions() (rowsAffected int64, err error)

	DB() *sql.DB
	ID() int
	Uuid() string
	Email() string
	User() UserInterface

	SetID(int)
	SetUuid(string)
	SetEmail(string)
	SetUser(UserInterface)

}

type user struct {
	db        *sql.DB
	id        int
	uuid      string
	name      string
	email     string
	password  string
	createdAt time.Time
}

type session struct {
	db   *sql.DB
	id   int
	uuid string

	// email is essentially from within user, but this field
	// is important to Scan back the row from DB
	email string

	user      UserInterface
	createdAt time.Time
}

func NewUser(db *sql.DB) UserInterface {
	return &user{db: db}
}

func NewSession(db *sql.DB, u UserInterface) SessionInterface{
	return &session{db: db, user: u}
}


func (u *user)FindByEmail(email string) (err error){
	row := u.db.QueryRow(
		"Select id, uuid, name, email, password, created_at FROM users where email=$1",
		email)

	if err = row.Scan(&u.id, &u.uuid, &u.name, &u.email, &u.password, &u.createdAt); err != nil {
		return
	}
	return
}

func (u *user)DeleteAllUsers() (rowsAffected int64, err error){
	query := "delete FROM users"
	result, err := u.db.Exec(query)
	if err != nil{
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (u *user) Create() (err error){
	if err = u.Validate(); err != nil{
		return
	}

	query := "insert INTO users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5) returning " +
		"id, uuid, name, email, password, created_at"
	stmt , err := u.db.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(CreateUUID(), u.name, u.email, Encrypt(u.password), time.Now()).Scan(
		&u.id, &u.uuid, &u.name, &u.email, &u.password, &u.createdAt)
	if err != nil {
		return
	}
	log.Info("user created: ", u)
	return
}

func (u *user) Validate() (err error){
	if u.name == ""{
		return &InvalidUser{Reason: "Empty Name"}
	}

	if u.password == "" {
		return &InvalidUser{Reason: "Password not set"}
	}

	// simple email validation
	if u.email == "" && !strings.Contains(u.email, "@") {
		return &InvalidUser{Reason: "email not valid"}
	}
	return
}

func (u *user) DB() *sql.DB {
	return u.db
}
func (u *user) ID() int {
	return u.id
}
func (u *user) Uuid() string {
	return u.uuid
}
func (u *user) Name() string {
	return u.name
}
func (u *user) Email() string {
	return u.email
}
func (u *user) Password() string {
	return u.password
}

func (u *user) setID(Id int) {
	u.id = Id
}
func (u *user) SetUuid(Uuid string) {
	u.uuid = Uuid
}
func (u *user) SetName(Name string) {
	u.name = Name
}
func (u *user) SetEmail(Email string) {
	u.email = Email
}
func (u *user) SetPassword(Password string) {
	u.password = Password
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
		s.user.Email(),
		s.user.ID(),
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
		Uuid).Scan(&s.id, &s.uuid, &s.email, s.user.ID(), &s.createdAt)
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

func (s *session) User() UserInterface{
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

func (s *session) SetUser(u UserInterface) {
	s.user = u
}
