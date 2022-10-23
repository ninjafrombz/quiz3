// Filename

package main 

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)


func (app *application) routes() *httprouter.Router {
	// Create
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/Notes", app.listNotesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/Notes", app.createNoteHandler)
	router.HandlerFunc(http.MethodGet, "/v1/Notes/:id", app.showNoteHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/Notes/:id", app.updateNoteHandler)
    router.HandlerFunc(http.MethodDelete, "/v1/Notes/:id", app.deleteNoteHandler)

	return router
}