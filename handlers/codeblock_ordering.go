package handlers

import (
	"database/sql"
	// "encoding/json"
	"net/http"
	// "strings"
	// "github.com/go-chi/chi/v5"
)

type CodeBlockOrdering struct {
	ID          int `json:"id"`
	PageID      int `json:"page_id"`
	TemplateID  int `json:"template_id"`
	CodeBlockID int `json:"codeblock_id"`
	Ordering    int `json:"ordering"`
	Active      int `json:"active"`
}

func AddToPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func AddToTemplate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
