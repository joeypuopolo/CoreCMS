package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CodeBlock struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Active      int     `json:"active"`
	Description *string `json:"description,omitempty"`
	Content     string  `json:"content"`
}

func CreateCodeBlock(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cb CodeBlock
		if err := json.NewDecoder(r.Body).Decode(&cb); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("INSERT INTO code_blocks (title, description, content) VALUES (?, ?, ?)", cb.Title, cb.Description, cb.Content)
		if err != nil {
			http.Error(w, "Failed to create code block", http.StatusInternalServerError)
			return
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Failed to retrieve last inserted ID", http.StatusInternalServerError)
			return
		}

		cb.ID = int(lastID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(cb)
	}
}

func GetCodeBlocks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, active, description, content FROM code_blocks ORDER BY id DESC")
		if err != nil {
			http.Error(w, "Failed to retrieve code blocks", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var codeBlocks []CodeBlock
		for rows.Next() {
			var cb CodeBlock
			if err := rows.Scan(
				&cb.ID,
				&cb.Title,
				&cb.Active,
				&cb.Description,
				&cb.Content,
			); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			codeBlocks = append(codeBlocks, cb)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(codeBlocks)
	}
}

func GetCodeBlock(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		codeBlockID := chi.URLParam(r, "codeBlockID")
		var codeBlock CodeBlock

		row := db.QueryRow("SELECT id, title, active, description, content FROM code_blocks WHERE id = ?", codeBlockID)
		if err := row.Scan(
			&codeBlock.ID,
			&codeBlock.Title,
			&codeBlock.Active,
			&codeBlock.Description,
			&codeBlock.Content,
		); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Code block not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to retrieve code block", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(codeBlock); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func UpdateCodeBlock(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		codeBlockID := chi.URLParam(r, "codeBlockID")
		var input struct {
			Title       *string `json:"title"`
			Active      *int    `json:"active"`
			Description *string `json:"description"`
			Content     *string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Build the query to update only the fields that are provided
		query := "UPDATE code_blocks SET"
		params := []interface{}{}
		if input.Title != nil {
			query += " title = ?,"
			params = append(params, *input.Title)
		}
		if input.Active != nil {
			query += " active = ?,"
			params = append(params, *input.Active)
		}
		if input.Description != nil {
			query += " description = ?,"
			params = append(params, *input.Description)
		}
		if input.Content != nil {
			query += " content = ?,"
			params = append(params, *input.Content)
		}

		// Removes trailing comma
		query = query[:len(query)-1] + " WHERE id = ?"
		params = append(params, codeBlockID)

		_, err := db.Exec(query, params...)
		if err != nil {
			http.Error(w, "Failed to update code block", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Code block updated successfully"))
	}
}

func DeleteCodeBlock(db *sql.DB) http.HandlerFunc {
	// also delete from codeblock_ordering
	return func(w http.ResponseWriter, r *http.Request) {
		codeBlockID := chi.URLParam(r, "codeBlockID")

		_, err := db.Exec("DELETE FROM code_blocks WHERE id = ?", codeBlockID)
		if err != nil {
			http.Error(w, "Failed to delete code block", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Code block deleted successfully"))
	}
}
