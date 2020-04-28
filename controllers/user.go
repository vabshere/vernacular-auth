package controllers

import (
	"net/http"
	"os"
	"regexp"
	"text/template"

	"github.com/vabshere/vernacular-auth/models"
	"github.com/vabshere/vernacular-auth/utils"

	"golang.org/x/crypto/bcrypt"
)

// newUser returns an instance of user with the given name, email and password
func newUser(name, email, password string) models.User {
	var u models.User
	if len(name) != 0 {
		u.Name = template.HTMLEscapeString(name)
	}
	u.Email = template.HTMLEscapeString(email)
	u.Password = []byte(template.HTMLEscapeString(password))
	return u
}

var saveUser = models.SaveUser

// SignUp creates a new user in the database and creates its session. Returns a pointer to user instance on success, nil otherwise.
func SignUp(w http.ResponseWriter, r *http.Request) *models.User {
	r.ParseForm()
	u := newUser(r.FormValue("name"), r.FormValue("email"), r.FormValue("password"))

	if len(u.Name) == 0 || len(u.Email) == 0 || len(u.Password) == 0 {
		utils.Respond(1, "Invalid submission", http.StatusBadRequest, w, r)
		return nil
	}

	if ok, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", u.Email); !ok {
		utils.Respond(1, "Invalid email", http.StatusBadRequest, w, r)
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		utils.Respond(1, "Error", http.StatusInternalServerError, w, r)
		return nil
	}

	u.Password = hash
	err = saveUser(&u)
	if err != nil {
		s := string(err.Error())
		if s[len("Error "):len("Error 1062")] == "1062" {
			utils.Respond(1, "Email already taken", http.StatusOK, w, r)
			return nil
		}

		utils.Respond(1, "Error", http.StatusInternalServerError, w, r)
		return nil
	}

	return &u
}

var getUserByEmail = models.GetUserByEmail

// SignIn checks if the user exists in the database and creates a session on successful attempt. Returns a pointer to user instance on success, nil otherwise
func SignIn(w http.ResponseWriter, r *http.Request) *models.User {
	r.ParseForm()
	u := newUser("", r.FormValue("email"), r.FormValue("password"))
	if len(u.Email) == 0 || len(u.Password) == 0 {
		utils.Respond(1, "Invalid submission", http.StatusBadRequest, w, r)
		return nil
	}

	if ok, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, u.Email); !ok {
		utils.Respond(1, "Invalid email", http.StatusBadRequest, w, r)
		return nil
	}

	user, err := getUserByEmail(u.Email)
	if err != nil {
		utils.Respond(1, "Error", http.StatusInternalServerError, w, r)
		return nil
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, u.Password); err != nil {
		utils.Respond(1, "Authentication failed", http.StatusOK, w, r)
		return nil
	}

	return user
}

// GetUser returns user from the session
func GetUser(w http.ResponseWriter, r *http.Request) {
	session, b := utils.GlobalSessions.SessionCheck(r)
	if b {
		u := utils.SessionGetUser(&session, r)
		utils.RespondJson(0, u, http.StatusOK, w, r)
		return
	}

	utils.RespondJson(1, nil, http.StatusOK, w, r)
	return
}

// SignOut deletes the user session
func SignOut(w http.ResponseWriter, r *http.Request) {
	utils.GlobalSessions.SessionDestroy(w, r)
	utils.Respond(0, "Success", http.StatusOK, w, r)
	return
}

// exists returns whether the given path (file or directory) exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}
