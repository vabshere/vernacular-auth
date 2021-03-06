# Vernacular-auth

A backend based on MVC pattern built using [Golang](https://golang.org) + [MySQL](https://www.mysql.com/).

# Installation

_This app_:<br />
`go get github.com/vabshere/github.com/vabshere/vernacular-auth`<br />

_External packages_:<br />
`go get http://golang.org/x/crypto/bcrypt`<br />
`go get github.com/go-sql-driver/mysql`<br />
`go get github.com/gorilla/mux`<br />

_MySQL_:<br />
You will also need to install [MySQL Server](https://www.mysql.com). This project was built using `version 8.0`. Refer the offficial docs for installation.

Update `connectDb()` function in `basic/models/user.go` with your username, password and database-name. Ensure you have a table named `user` in your database. `CREATE TABLE` for the table `user`:<br />
```
CREATE TABLE user (
  id int NOT NULL AUTO_INCREMENT PRIMARY KEY UNIQUE,
  email varchar(50) DEFAULT NULL UNIQUE,
  name varchar(30) DEFAULT NULL,
  password binary(60) DEFAULT NULL,
);
```

# Usage

Run `go run /path/to/main.go`. This should start the local server on any interface `:8080`.
