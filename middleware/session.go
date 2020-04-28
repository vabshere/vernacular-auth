package middleware

import (
	"net/http"

	"github.com/vabshere/vernacular-auth/models"
	"github.com/vabshere/vernacular-auth/utils"
)

type SessionReset func(http.ResponseWriter, *http.Request) *models.User

func (handler SessionReset) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	utils.GlobalSessions.SessionDestroy(w, r)
	if user := handler(w, r); user != nil {
		session := utils.GlobalSessions.SessionStart(w, r)
		utils.SessionSetUser(user, &session, r)
		utils.RespondJson(0, user, http.StatusOK, w, r)
	}
}
