package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// connectDB connects to the database server
func connectDb() (*sql.DB, error) {
	username := "root"
	password := "mysql"
	dbName := "basic"
	db, err := sql.Open("mysql", ""+username+":"+password+"@/"+dbName+"?charset=utf8")
	if err != nil {
		return nil, err
	}

	return db, nil
}

type password []byte

// MarshalJSON is the custom method used by JSON.Marshal for the "password" type.
func (password) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

// User is the tye of all users
type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password password `json:"password"`
	Id       int      `json:"id"`
}

// SaveUser saves a user into the database
func SaveUser(u *User) error {
	db, err := connectDb()
	if err != nil {
		return err
	}

	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO user (name, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(u.Name, u.Email, u.Password)
	if err == nil {
		id, _ := res.LastInsertId()
		u.Id = int(id)
	}

	return err
}

// GetUserByEmail returns the user associated with given email
func GetUserByEmail(email string) (*User, error) {
	db, err := connectDb()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("SELECT * FROM user WHERE email=?")
	if err != nil {
		return nil, err
	}

	var user User
	err = stmt.QueryRow(email).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	return &user, err
}
