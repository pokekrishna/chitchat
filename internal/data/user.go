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

	GetDB() *sql.DB
	GetID() int
	GetUuid() string
	GetName() string
	GetEmail() string
	GetPassword() string

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

	GetDB() *sql.DB
	GetID() int
	GetUuid() string
	GetEmail() string
	GetUser() UserInterface

	SetID(int)
	SetUuid(string)
	SetEmail(string)
	SetUser(UserInterface)

}

type user struct {
	DB *sql.DB
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type session struct {
	DB        *sql.DB
	Id        int
	Uuid      string

	// Email is essentially from within User, but this field
	// is important to Scan back the row from DB
	Email     string

	User  UserInterface
	CreatedAt time.Time
}

func NewUser(db *sql.DB) UserInterface {
	return &user{DB: db}
}

func NewSession(db *sql.DB, u UserInterface) *session{
	return &session{DB: db, User: u}
}


func (u *user)FindByEmail(email string) (err error){
	row := u.DB.QueryRow(
		"Select id, uuid, name, email, password, created_at FROM users where email=$1",
		email)

	if err = row.Scan(&u.Id, &u.Uuid, &u.Name, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		return
	}
	return
}

func (u *user)DeleteAllUsers() (rowsAffected int64, err error){
	query := "delete FROM users"
	result, err := u.DB.Exec(query)
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
	stmt , err := u.DB.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(CreateUUID(), u.Name, u.Email, Encrypt(u.Password), time.Now()).Scan(
		&u.Id, &u.Uuid, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return
	}
	log.Info("User created: ", u)
	return
}

func (u *user) Validate() (err error){
	if u.Name == ""{
		return &InvalidUser{Reason: "Empty Name"}
	}

	if u.Password == "" {
		return &InvalidUser{Reason: "Password not set"}
	}

	// simple email validation
	if u.Email == "" && !strings.Contains(u.Email, "@") {
		return &InvalidUser{Reason: "Email not valid"}
	}
	return
}

func (u *user) GetDB() *sql.DB {
	return u.DB
}
func (u *user) GetID() int {
	return u.Id
}
func (u *user) GetUuid() string {
	return u.Uuid
}
func (u *user) GetName() string {
	return u.Name
}
func (u *user) GetEmail() string {
	return u.Email
}
func (u *user) GetPassword() string {
	return u.Password
}

func (u *user) setID(Id int) {
	u.Id = Id
}
func (u *user) SetUuid(Uuid string) {
	u.Uuid = Uuid
}
func (u *user) SetName(Name string) {
	u.Name = Name
}
func (u *user) SetEmail(Email string) {
	u.Email = Email
}
func (u *user) SetPassword(Password string) {
	u.Password = Password
}

func (s *session) Create()  (err error){
	var stmt *sql.Stmt
	query := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, created_at"
	stmt, err = s.DB.Prepare(query)
	if err != nil {
		log.Error("Cannot prepare stmt", err)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(CreateUUID(),
		s.User.GetEmail(),
		s.User.GetID(),
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

func (s *session)FindByUuid(Uuid string) (err error) {
	err = s.DB.QueryRow("select id, uuid, email, user_id, created_at from sessions where uuid=$1",
		Uuid).Scan(&s.Id, &s.Uuid, &s.Email, s.User.GetID(), &s.CreatedAt)
	if err != nil {
		return
	}
	return
}

func (s *session) Delete() (err error) {
	query := "Delete from sessions where uuid=$1"
	stmt, err := s.DB.Prepare(query)
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

func (s *session)DeleteAllSessions() (rowsAffected int64, err error){
	query := "delete FROM sessions"
	result, err := s.DB.Exec(query)
	if err != nil{
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return
	}
	return
}

func (s *session) GetDB() *sql.DB{
	return s.DB
}

func (s *session) GetID() int{
	return s.Id
}

func (s *session) GetUuid() string{
	return s.Uuid
}

func (s *session) GetEmail() string{
	return s.Email
}

func (s *session) GetUser() UserInterface{
	return s.User
}

func (s *session) SetUuid(Uuid string) {
	s.Uuid = Uuid
}

func (s *session) SetID(Id int) {
	s.Id = Id
}

func (s *session) SetEmail(Email string) {
	s.Email = Email
}

func (s *session) SetUser(u UserInterface) {
	s.User = u
}
