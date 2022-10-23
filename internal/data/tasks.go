// Filename: /internals/data/notes.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"quiz3.desireamagwula.net/internal/validator"
)

type Note struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Task_Name   string    `json:"task_name"`
	Description string    `json:"desription"`
	Category    string    `json:"category"`
	Priority    string    `json:"priority"`
	Status      []string  `json:"status"`
	Version     int32     `json:"version"`
	// ID        int64     `json:"id"`
	// CreatedAt time.Time `json:"-"`
	// Name      string    `json:"name"`
	// Level     string    `json:"level"`
	// Contact   string    `json:"contact"`
	// Phone     string    `json:"phone"`
	// Email     string    `json:"email,omitempty"`
	// Website   string    `json:"website,omitempty"`
	// Address   string    `json:"address"`
	// Mode      []string  `json:"mode"`
	// Version   int32     `json:"version"`
}

func ValidateNote(v *validator.Validator, note *Note) {
	// Use the Check() Method to execute our validation checks
	v.Check(note.Task_Name != "", "name", "must be provided")
	v.Check(len(note.Task_Name) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(note.Description != "", "level", "must be provided")
	v.Check(len(note.Description) <= 200, "level", "must not be more than 200 bytes long")

	v.Check(note.Category != "", "contact", "must be provided")
	v.Check(len(note.Category) <= 200, "contact", "must not be more than 200 bytes long")

	v.Check(note.Priority != "", "address", "must be provided")
	v.Check(len(note.Priority) <= 500, "address", "must not be more than 200 bytes long")

	v.Check(note.Status != nil, "mode", "must be provided!")
	v.Check(len(note.Status) >= 1, "mode", "must contain at least one entry")
	v.Check(len(note.Status) <= 5, "mode", "must contain at most five entries")
	v.Check(validator.Unique(note.Status), "mode", "must not contain duplicate entries")

}

type NoteModel struct {
	DB *sql.DB
}

// Insert() allows us to create a new note

func (m NoteModel) Insert(note *Note) error {
	query := `
		INSERT INTO notes (task_name, description, category, priority, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version
	`
	// Collect the data fields into a slice
	args := []interface{}{
		note.Task_Name, note.Description,
		note.Category, note.Priority,
		pq.Array(note.Status),
	}
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&note.ID, &note.CreatedAt, &note.Version)
	//return m.DB.QueryRow(query, args...).Scan(&note.ID, &note.CreatedAt, &note.Version)
}

// Get() allows us to retrieve

func (m NoteModel) Get(id int64) (*Note, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Create the query
	query := `
		SELECT id, created_at, task_name, description, category, priority, status, version
		FROM notes
		WHERE id = $1
	`
	// Declare a note variable to hold the returned data
	var note Note
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	// Execute the query using QueryRow()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&note.ID,
		&note.CreatedAt,
		&note.Task_Name,
		&note.Description,
		&note.Category,
		&note.Priority,
		pq.Array(&note.Status),
		&note.Version,
	)
	// Handle any errors
	if err != nil {
		// Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Success
	return &note, nil
}

// Update() allows us to edit/alter a specific note

func (m NoteModel) Update(note *Note) error {
	// Create a query
	query := `
		UPDATE notes
		SET task_name = $1, description = $2, category = $3,
		    priority = $4, status = $5, version = version + 1
		WHERE id = $6
		AND version = $7
		RETURNING version
	`
	args := []interface{}{
		&note.Task_Name,
		&note.Description,
		&note.Category,
		&note.Priority,
		pq.Array(&note.Status),
		note.ID,
		note.Version,
	}

	//Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleanup to prevent memory leaks
	defer cancel()
	// Check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&note.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil

}

// Delete removes a specific note
func (m NoteModel) Delete(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}
	// Create the delete query
	query := `
		DELETE FROM notes
		WHERE id = $1
	`
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	// Execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Check how many rows were affected by the delete operation. We
	// call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}
func (m NoteModel) GetAll(task_name string, description string, status []string, filters Filters) ([]*Note, Metadata, error) {
	// Construct the query
	query := fmt.Sprintf(`
		SELECT COUNT (*) OVER(), id, created_at, task_name, description, category, priority, status, version
		FROM notes
		WHERE (to_tsvector('simple', task_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', description) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (status @> $3 OR $3 = '{}' )
		ORDER by %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortOrder())
	// Create
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{task_name, description, pq.Array(status), filters.limit(), filters.offSet()}
	// Execute the query
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	// Initialize an empty slice
	notes := []*Note{}
	// iterate over the rows in the resultset
	for rows.Next() {
		var note Note
		// SCan the valuies from the row into the note
		err := rows.Scan(
			&totalRecords,
			&note.ID,
			&note.CreatedAt,
			&note.Task_Name,
			&note.Description,
			&note.Category,
			&note.Priority,
			pq.Array(&note.Status),
			&note.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		notes = append(notes, &note)

	}
	// Check if any errors occured after looping through the resultset
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// safely return the resultset
	return notes, metadata, nil
}
