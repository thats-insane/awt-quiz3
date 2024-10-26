package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *appDependencies) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.notAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthCheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/signup/:id", a.displayUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/signup/:id", a.updateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/signup", a.createUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/signup/:id", a.deleteUserHandler)
	return router
}
