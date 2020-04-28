package main

import (
	"log"
	"net/http"

	"github.com/vabshere/vernacular-auth/routes"
	"github.com/vabshere/vernacular-auth/utils"
	_ "github.com/vabshere/vernacular-auth/utils/session/providers/memory"
)

func main() {
	utils.Run()
	r := routes.Init()
	println("running server")
	log.Fatal(http.ListenAndServe(":8080", r))
}
