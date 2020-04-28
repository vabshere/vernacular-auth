package utils

import (
	"net/http"
	"github.com/vabshere/vernacular-auth/models"
	"github.com/vabshere/vernacular-auth/utils/session"
)

// GlobalSessions is the global variable for managing sessions
var GlobalSessions *session.Manager

// Run initializes the one time configurations required for using sessions
func Run() {
	GlobalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go GlobalSessions.GC()
}

// SessionSetUser is used for setting given user's details in given session
func SessionSetUser(user *models.User, session *session.Session, r *http.Request) {
	(*session).Set("id", user.Id)
	(*session).Set("name", user.Name)
	(*session).Set("email", user.Email)
}

// SessionGetUser returns user details from given session
func SessionGetUser(session *session.Session, r *http.Request) *models.User {
	id := (*session).Get("id").(int)
	name := (*session).Get("name").(string)
	email := (*session).Get("email").(string)
	u := models.User{Id: id, Name: name, Email: email}
	return &u
}
