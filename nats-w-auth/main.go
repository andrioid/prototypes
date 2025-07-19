package main

import (
	"html/template"
	"nats-w-auth/internal/users"

	"github.com/alexedwards/scs/v2"
	"github.com/nats-io/nats-server/v2/server"
)

var natsd *server.Server
var templates *template.Template
var sessionManager *scs.SessionManager
var userManager *users.UserManager

func main() {
	setupNats()
	setupSession()
	userManager = users.New(natsd.ClientURL())
	httpServe()

}
