// TODO:
// - DuplicateTemplate(){}
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Template struct {
	ID               int                 `json:"id"`
	Title            string              `json:"title"`
	ParentTemplateID *int                `json:"parent_template_id"`
	CodeBlocks       []CodeBlockOrdering `json:"codeblocks"`
}

func CreateTemplate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Title string `json:"title"`
		}

		// Decode the JSON request body into the input struct
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate that the title is provided
		if input.Title == "" {
			http.Error(w, "Template title is required", http.StatusBadRequest)
			return
		}

		// Insert the new template into the database
		query := "INSERT INTO templates (title) VALUES (?)"
		result, err := db.Exec(query, input.Title)
		if err != nil {
			http.Error(w, "Failed to create template", http.StatusInternalServerError)
			return
		}

		// Get the ID of the newly inserted template
		templateID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Failed to retrieve new template ID", http.StatusInternalServerError)
			return
		}

		// Prepare the response with the created template
		template := Template{
			ID:    int(templateID),
			Title: input.Title,
		}

		// Return the created template as a JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(template); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func DuplicateTemplate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the template ID to duplicate from the URL
		templateID := chi.URLParam(r, "templateID")
		if templateID == "" {
			http.Error(w, "Template ID is required", http.StatusBadRequest)
			return
		}

		// Fetch the template details
		var template Template
		err := db.QueryRow(`
			SELECT 
			id, 
			title, 
			parent_template_id
			--json_extract(template_codeblocks, '$') AS template_codeblocks 
			FROM templates 
			WHERE id = ?`,
			templateID).Scan(
			&template.ID,
			&template.Title,
			&template.ParentTemplateID,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Template not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to fetch template: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Create a new template record with a duplicated title
		newTemplateName := template.Title + " (Copy)"
		res, err := db.Exec(`
			INSERT INTO templates (title, parent_template_id) 
			VALUES (?, ?)`,
			newTemplateName, template.ParentTemplateID,
		)
		if err != nil {
			http.Error(w, "Failed to duplicate template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the new template ID
		newTemplateID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Failed to get new template ID: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Duplicate the template's code blocks
		// pulls codeblocks from CodeblocksOrdering

		// for _, codeBlockID := range *template.TemplateCodeBlocks {
		// 	_, err = db.Exec(`
		// 		INSERT INTO code_blocks (template_id, content)
		// 		SELECT ?, content
		// 		FROM code_blocks
		// 		WHERE id = ?`, newTemplateID, codeBlockID)
		// 	if err != nil {
		// 		http.Error(w, "Failed to duplicate code block: "+err.Error(), http.StatusInternalServerError)
		// 		return
		// 	}
		// }

		// Respond with the new template ID
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"new_template_id": %d}`, newTemplateID)))
	}
}

func GetTemplates(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, parent_template_id FROM templates")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var templates []Template
		for rows.Next() {
			var tmpl Template
			if err := rows.Scan(&tmpl.ID, &tmpl.Title, &tmpl.ParentTemplateID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			templates = append(templates, tmpl)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(templates)
	}
}

func GetTemplate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the templateID from the URL parameters
		templateID := chi.URLParam(r, "templateID")

		// Convert templateID to an integer
		id, err := strconv.Atoi(templateID)
		if err != nil {
			http.Error(w, "Invalid template ID", http.StatusBadRequest)
			return
		}

		// Query the database for the template
		row := db.QueryRow("SELECT id, title, parent_template_id FROM templates WHERE id = ?", id)

		// Scan the result into a Template struct
		var tmpl Template
		err = row.Scan(&tmpl.ID, &tmpl.Title, &tmpl.ParentTemplateID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Template not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Respond with the template as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tmpl)
	}
}

func UpdateTemplate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the template ID from the URL
		templateID := chi.URLParam(r, "id")
		if templateID == "" {
			http.Error(w, "Template ID is required", http.StatusBadRequest)
			return
		}

		// Parse the template ID into an integer
		id, err := strconv.Atoi(templateID)
		if err != nil {
			http.Error(w, "Invalid template ID", http.StatusBadRequest)
			return
		}

		// Decode the JSON payload into a partial Template struct
		var updatedFields map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updatedFields); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if len(updatedFields) == 0 {
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		// Build the dynamic SQL query
		query := "UPDATE templates SET "
		args := []interface{}{}
		for field, value := range updatedFields {
			switch field {
			case "title":
				query += "title = ?, "
				args = append(args, value)
			case "parent_template_id":
				query += "parent_template_id = ?, "
				args = append(args, value)
			default:
				http.Error(w, fmt.Sprintf("Field '%s' cannot be updated", field), http.StatusBadRequest)
				return
			}
		}

		// Remove the trailing comma and space, and add the WHERE clause
		query = strings.TrimSuffix(query, ", ") + " WHERE id = ?"
		args = append(args, id)

		// Execute the query
		_, err = db.Exec(query, args...)
		if err != nil {
			http.Error(w, "Failed to update template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Template updated successfully"))
	}
}

// func fetchInheritedCodeBlocks(db *sql.DB, templateID int) ([]int, error) {
// 	var inheritedCodeBlocks []int
// 	var parentTemplateID *int

// 	// Fetch the parent template ID
// 	err := db.QueryRow("SELECT parent_template_id FROM templates WHERE id = ?", templateID).Scan(&parentTemplateID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// If there's a parent template, recursively fetch its code blocks
// 	if parentTemplateID != nil {
// 		parentCodeBlocks, err := fetchInheritedCodeBlocks(db, *parentTemplateID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		inheritedCodeBlocks = append(inheritedCodeBlocks, parentCodeBlocks...)
// 	}

// 	// Fetch the code blocks from the current template
// 	var templateCodeBlocks []int
// 	err = db.QueryRow("SELECT template_codeblocks FROM templates WHERE id = ?", templateID).Scan(&templateCodeBlocks)
// 	if err != nil {
// 		return nil, err
// 	}
// 	inheritedCodeBlocks = append(inheritedCodeBlocks, templateCodeBlocks...)

// 	return inheritedCodeBlocks, nil
// }

// func updateAllCodeBlocks(db *sql.DB, templateID int) error {
// 	inheritedCodeBlocks, err := fetchInheritedCodeBlocks(db, templateID)
// 	if err != nil {
// 		return err
// 	}

// 	// Fetch the existing code block IDs from the template
// 	// rows, err := db.Query("SELECT id, title, parent_template_id FROM templates")
// 	// 	if err != nil {
// 	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// 		return
// 	// 	}
// 	// 	defer rows.Close()

// 	// 	var templates []Template
// 	// 	for rows.Next() {
// 	// 		var tmpl Template
// 	// 		if err := rows.Scan(&tmpl.ID, &tmpl.Title, &tmpl.ParentTemplateID); err != nil {
// 	// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// 			return
// 	// 		}
// 	// 		templates = append(templates, tmpl)
// 	// 	}

// 	rows, err := db.Query("SELECT codeblock_id FROM codeblocks_ordering WHERE entity_type = 'template' AND codeblock_id = ?", templateID).Scan(&existingCodeBlocks)
// 	// var existingCodeBlocks []int
// 	// err = db.QueryRow("SELECT all_codeblocks FROM templates WHERE id = ?", templateID).Scan(&existingCodeBlocks)
// 	// if err != nil {
// 		// return err
// 	// }

// 	// Create a map to keep track of the existing code block IDs
// 	existingCodeBlockMap := make(map[int]bool)
// 	for _, codeBlockID := range existingCodeBlocks {
// 		existingCodeBlockMap[codeBlockID] = true
// 	}

// 	// Append the inherited code block IDs to the existing code block IDs
// 	for _, codeBlockID := range inheritedCodeBlocks {
// 		if !existingCodeBlockMap[codeBlockID] {
// 			existingCodeBlocks = append(existingCodeBlocks, codeBlockID)
// 		}
// 	}

// 	// Update the AllCodeBlocks field with the updated code block IDs
// 	_, err = db.Exec("UPDATE templates SET all_codeblocks = ? WHERE id = ?", existingCodeBlocks, templateID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func DeleteTemplate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateID := chi.URLParam(r, "templateID")
		id, err := strconv.Atoi(templateID)
		if err != nil {
			http.Error(w, "Invalid template ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM templates WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Template deleted successfully"))
	}
}
