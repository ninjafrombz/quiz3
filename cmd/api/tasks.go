// Filename:
package main

import (
	"errors"
	"fmt"
	"net/http"

	"quiz3.desireamagwula.net/internal/data"
	"quiz3.desireamagwula.net/internal/validator"
)

// CreateNoteHandler for the POST /v1/Notes" endpoint

func (app *application) createNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Our target decode destination fmt.Fprintln(w, "create a new Note..")
	var input struct {
		Name    string   `json:"name"`
		Level   string   `json:"level"`
		Contact string   `json:"contact"`
		Phone   string   `json:"phone"`
		Email   string   `json:"email"`
		Website string   `json:"website"`
		Address string   `json:"address"`
		Mode    []string `json:"mode"`
	}

	// Initialize a new json.Decoder instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new Note struct
	Note := &data.Note{
		Name:    input.Name,
		Level:   input.Level,
		Contact: input.Contact,
		Phone:   input.Phone,
		Email:   input.Email,
		Website: input.Website,
		Address: input.Address,
		Mode:    input.Mode,
	}

	//Initialize a new validator instance
	v := validator.New()

	// Check the map to determine if there were any validation errors
	if data.ValidateNote(v, Note); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// CReate a Note
	err = app.models.Notes.Insert(Note)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// CReate a location header for the newly created
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/Notes/%d", Note.ID))
	//Write the JSON response with 201 - Created status code with the body
	// being the Note data and the header being the headers map

	err = app.writeJSON(w, http.StatusCreated, envelope{"Note": Note}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}

}

func (app *application) showNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the specific Note
	Note, err := app.models.Notes.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	// Write the sdata returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"Note": Note}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	// This method does a partial replacement
	// Get the id for the Note that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the orginal record from the database
	Note, err := app.models.Notes.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Create an input struct to hold data read in from the client
	// We update input struct to use pointers because pointers have a
	// default value of nil
	// If a field remains nil then we know that the client did not update it
	var input struct {
		Name    *string  `json:"name"`
		Level   *string  `json:"level"`
		Contact *string  `json:"contact"`
		Phone   *string  `json:"phone"`
		Email   *string  `json:"email"`
		Website *string  `json:"website"`
		Address *string  `json:"address"`
		Mode    []string `json:"mode"`
	}

	// Initialize a new json.Decoder instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Check for updates
	if input.Name != nil {
		Note.Name = *input.Name
	}
	if input.Level != nil {
		Note.Level = *input.Level
	}
	if input.Contact != nil {
		Note.Contact = *input.Contact
	}
	if input.Phone != nil {
		Note.Phone = *input.Phone
	}
	if input.Email != nil {
		Note.Email = *input.Email
	}
	if input.Website != nil {
		Note.Website = *input.Website
	}
	if input.Address != nil {
		Note.Address = *input.Address
	}
	if input.Mode != nil {
		Note.Mode = input.Mode
	}

	// Perform validation on the updated Note. If validation fails, then
	// we send a 422 - Unprocessable Entity respose to the client
	// Initialize a new Validator instance
	v := validator.New()

	// Check the map to determine if there were any validation errors
	if data.ValidateNote(v, Note); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Let's pass the updated Note record to the Update() method
	err = app.models.Notes.Update(Note)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"Note": Note}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteNoteHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the Note from the Database. Send a 404 not found status cide to the client
	// if not found

	err = app.models.Notes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}
	// Return 200 Status OK to the client with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Note successfuly deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// Allows the client to see a listing of Notes based on a set of criterias

func (app *application) listNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Create an input struct to hold our query paraneters
	var input struct {
		Name  string
		Level string
		Mode  []string
		data.Filters
	}
	v := validator.New()
	// Get the url values map
	qs := r.URL.Query()
	// Use the helper methods to extfract the values
	input.Name = app.readString(qs, "name", "")
	input.Level = app.readString(qs, "level", "")
	input.Mode = app.readCSV(qs, "mode", []string{})
	//Get the page information
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Get the sort info
	input.Filters.Sort = app.readString(qs, "sort", "id")
	// Specify the allowed sort values
	input.Filters.SortList = []string{"id", "name", "level", "-id", "-name", "-level"}
	// CHeck for validation error
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Get a listing of all Notes
	Notes, metadata, err := app.models.Notes.GetAll(input.Name, input.Level, input.Mode, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containing all the Notes
	err = app.writeJSON(w, http.StatusOK, envelope{"Notes": Notes, "metadata ": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
