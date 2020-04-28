package routes

import (
	"net/http"

	"github.com/vabshere/vernacular-auth/controllers"
	"github.com/vabshere/vernacular-auth/middleware"

	"github.com/gorilla/mux"
)

//Init initializes routes for the app
func Init() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/reg", middleware.SessionReset(controllers.SignUp)).Methods(http.MethodPost)
	r.Handle("/oauth", middleware.SessionReset(controllers.SignIn)).Methods(http.MethodPost)
	r.HandleFunc("/home", controllers.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/signOut", controllers.SignOut).Methods(http.MethodDelete)
	return r
}
