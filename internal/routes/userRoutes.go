package routes

import "net/http"

var UserRoutes = map[string]http.HandlerFunc{
	"/register": handlers.registerUser,
}
