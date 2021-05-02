package data

import (
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

func UserByEmail(email string) (u *User, err error){
	u = &User{}
	row, err := db.Query(
		"Select id, uuid, name, email, password, created_at FROM users where email=$1",
		email)
	if err != nil {
		return
	}
	if err= row.Scan(&u.Id, &u.Uuid, &u.Name, &u.Email, &u.Password, &u.CreatedAt); err != nil{
		return
	}
	return
}


