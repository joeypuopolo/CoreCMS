package handlers

// TODO
// Codeblocks able to set to Hidden
// Settings { "hidden_codeblocks": [1]}

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Page struct {
	ID         int                 `json:"id"`
	Title      string              `json:"title"`
	Url        string              `json:"url"`
	Hidden     int                 `json:"hidden"`
	Active     int                 `json:"active"`
	Link       *string             `json:"link,omitempty"`
	LinkNewTab *int                `json:"link_new_tab,omitempty"`
	ParentPage int                 `json:"parent_page"`
	Settings   *string             `json:"settings,omitempty"`
	TemplateID int                 `json:"template_id"`
	CodeBlocks []CodeBlockOrdering `json:"codeblocks"`
}

func CreatePage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body into the Page struct
		var requestData struct {
			Title      string `json:"title"`
			Url        string `json:"url"`
			TemplateID int    `json:"template_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validate the required fields
		if requestData.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}
		if requestData.Url == "" {
			http.Error(w, "Url is required", http.StatusBadRequest)
			return
		}
		if requestData.TemplateID == 0 {
			// Query for a template
			http.Error(w, "Template ID is required", http.StatusBadRequest)
			return
		}

		// Insert the new page into the database
		result, err := db.Exec(`
            INSERT INTO pages (title, url, template_id) VALUES (?, ?, ?)`,
			requestData.Title, requestData.Url, requestData.TemplateID,
		)
		if err != nil {
			http.Error(w, "Failed to create page: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the ID of the newly created page
		newPageID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Failed to retrieve new page ID: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the created page ID
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Page created successfully",
			"id":      newPageID,
		})
	}
}

func GetPages(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM pages")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var pages []Page
		for rows.Next() {
			var page Page
			if err := rows.Scan(
				&page.ID,
				&page.Title,
				&page.Url,
				&page.Hidden,
				&page.Active,
				&page.Link,
				&page.LinkNewTab,
				&page.ParentPage,
				&page.Settings,
				&page.TemplateID,
			); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pages = append(pages, page)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pages)
	}
}

func GetPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id") // Get page ID from route
		var page Page

		// Fetch page details
		err := db.QueryRow("SELECT id, title, url, parent_page FROM pages WHERE id = ?", id).
			Scan(&page.ID, &page.Title, &page.Url, &page.ParentPage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch ordered codeblocks
		rows, err := db.Query("SELECT id, entity_id, codeblock_id, ordering, active FROM codeblocks_ordering WHERE entity_type = 'page' AND entity_id = ? ORDER BY ordering", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var cb CodeBlockOrdering
			if err := rows.Scan(&cb.ID, &cb.PageID, &cb.TemplateID, &cb.CodeBlockID, &cb.Ordering); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			page.CodeBlocks = append(page.CodeBlocks, cb)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
	}
}

func UpdatePage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageID := chi.URLParam(r, "pageID")

		var input struct {
			Title      *string `json:"title"`
			Url        *string `json:"url"`
			Hidden     *int    `json:"hidden"`
			Active     *int    `json:"active"`
			Link       *string `json:"link"`
			LinkNewTab *int    `json:"link_new_tab"`
			ParentPage *int    `json:"parent_page"`
			Settings   *string `json:"settings"`
			TemplateID *int    `json:"template_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Build the query to update only the fields that are provided
		query := "UPDATE pages SET"
		params := []interface{}{}

		if input.Title != nil {
			query += " title = ?,"
			params = append(params, *input.Title)
		}
		if input.Url != nil {
			query += " url = ?,"
			params = append(params, *input.Url)
		}
		if input.Hidden != nil {
			query += " hidden = ?,"
			params = append(params, *input.Hidden)
		}
		if input.Active != nil {
			query += " active = ?,"
			params = append(params, *input.Active)
		}
		if input.Link != nil {
			query += " link = ?,"
			params = append(params, *input.Link)
		}
		if input.LinkNewTab != nil {
			query += " link_new_tab = ?,"
			params = append(params, *input.LinkNewTab)
		}
		if input.ParentPage != nil {
			query += " parent_page = ?,"
			params = append(params, *input.ParentPage)
		}
		if input.Settings != nil {
			query += " settings = ?,"
			params = append(params, *input.Settings)
		}
		if input.TemplateID != nil {
			query += " template_id = ?,"
			params = append(params, *input.TemplateID)
		}

		// Remove trailing comma and add the WHERE clause
		query = query[:len(query)-1] + " WHERE id = ?"
		params = append(params, pageID)

		_, err := db.Exec(query, params...)
		if err != nil {
			http.Error(w, "Failed to update page", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Page updated successfully"))
	}
}

// Helper function to format a slice of integers as JSON
// func formatIntSliceToJSON(slice []int) string {
// 	jsonData, _ := json.Marshal(slice)
// 	return string(jsonData)
// }

func RenderPage(db *sql.DB) http.HandlerFunc {
	// Needs to:
	// if link, then 301 redirect to link
	// if no link, then return list of codeblock ids
	return func(w http.ResponseWriter, r *http.Request) {
		pageID := chi.URLParam(r, "pageID")
		if pageID == "" {
			http.Error(w, "Page ID is required", http.StatusBadRequest)
			return
		}

		// Fetch the page data
		var page Page

		row := db.QueryRow(`
			SELECT id, title, url, hidden, active, link, link_new_tab, parent_page, settings, template_id FROM pages WHERE id = ?`, pageID)
		err := row.Scan(
			&page.ID,
			&page.Title,
			&page.Url,
			&page.Hidden,
			&page.Active,
			&page.Link,
			&page.LinkNewTab,
			&page.ParentPage,
			&page.Settings,
			&page.TemplateID,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Page not found", http.StatusNotFound)
			} else {
				http.Error(w, "Error fetching page data: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Fetch ordered codeblocks
		rows, err := db.Query("SELECT id, page_id, template_id, codeblock_id, ordering, active FROM codeblocks_ordering WHERE page_id = ? OR template_id = ? ORDER BY ordering", pageID, &page.TemplateID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var cb CodeBlockOrdering
			if err := rows.Scan(&cb.ID, &cb.PageID, &cb.TemplateID, &cb.CodeBlockID, &cb.Ordering); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			page.CodeBlocks = append(page.CodeBlocks, cb)
		}

		renderedBlocks, err := renderCodeBlocks(db, page.CodeBlocks)
		if err != nil {
			http.Error(w, "Error rendering code blocks: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Serve the final rendered content
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(strings.Join(renderedBlocks, "")))
	}
}

// Helper function to fetch code blocks from a template and its parent templates
// func fetchTemplateCodeBlocks(db *sql.DB, templateID int) ([]int, error) {
// 	var template Template
// 	var codeBlocksOrder, templateCodeBlocks string

// 	err := db.QueryRow(`
// 		SELECT id, name, parent_template_id, all_codeblocks, template_codeblocks
// 		FROM templates WHERE id = ?`, templateID).
// 		Scan(&template.ID, &template.Title, &template.ParentTemplateID, &codeBlocksOrder, &templateCodeBlocks)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Parse JSON fields
// 	if err := json.Unmarshal([]byte(codeBlocksOrder), &template.AllCodeBlocks); err != nil {
// 		return nil, err
// 	}
// 	if err := json.Unmarshal([]byte(templateCodeBlocks), &template.TemplateCodeBlocks); err != nil {
// 		return nil, err
// 	}

// 	// If there's a parent template, recursively fetch its code blocks
// 	if template.ParentTemplateID != nil {
// 		parentBlocks, err := fetchTemplateCodeBlocks(db, *template.ParentTemplateID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return append(parentBlocks, *template.TemplateCodeBlocks...), nil
// 	}

// 	return *template.TemplateCodeBlocks, nil
// }

// Helper function to render code blocks based on their IDs
func renderCodeBlocks(db *sql.DB, blockIDs []CodeBlockOrdering) ([]string, error) {
	var renderedBlocks []string

	for _, blockID := range blockIDs {
		var content string
		err := db.QueryRow(`
			SELECT content FROM code_blocks WHERE id = ?`, blockID.CodeBlockID).Scan(&content)
		if err != nil {
			return nil, err
		}
		renderedBlocks = append(renderedBlocks, content)
	}

	return renderedBlocks, nil
}

// func fetchInheritedCodeBlocksForPage(db *sql.DB, pageID int) ([]int, error) {
// 	var inheritedCodeBlocks []int
// 	var templateID int

// 	// Fetch the template ID from the page
// 	err := db.QueryRow("SELECT template_id FROM pages WHERE id = ?", pageID).Scan(&templateID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Fetch the inherited code blocks from the template
// 	templateCodeBlocks, err := fetchInheritedCodeBlocks(db, templateID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	inheritedCodeBlocks = append(inheritedCodeBlocks, templateCodeBlocks...)

// 	// Fetch the code blocks from the parent page
// 	var parentPageID int
// 	err = db.QueryRow("SELECT parent_page FROM pages WHERE id = ?", pageID).Scan(&parentPageID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if parentPageID != -1 {
// 		parentPageCodeBlocks, err := fetchInheritedCodeBlocksForPage(db, parentPageID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		inheritedCodeBlocks = append(inheritedCodeBlocks, parentPageCodeBlocks...)
// 	}

// 	return inheritedCodeBlocks, nil
// }

// func updateAllCodeBlocksForPage(db *sql.DB, pageID int) error {
// 	inheritedCodeBlocks, err := fetchInheritedCodeBlocksForPage(db, pageID)
// 	if err != nil {
// 		return err
// 	}

// 	// Fetch the existing code block IDs from the page
// 	var existingCodeBlocks []int
// 	err = db.QueryRow("SELECT all_codeblocks FROM pages WHERE id = ?", pageID).Scan(&existingCodeBlocks)
// 	if err != nil {
// 		return err
// 	}

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
// 	_, err = db.Exec("UPDATE pages SET all_codeblocks = ? WHERE id = ?", existingCodeBlocks, pageID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func DeletePage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the page ID from the URL
		pageID := chi.URLParam(r, "pageID")
		if pageID == "" {
			http.Error(w, "Page ID is required", http.StatusBadRequest)
			return
		}

		// Execute the DELETE query
		result, err := db.Exec("DELETE FROM pages WHERE id = ?", pageID)
		if err != nil {
			http.Error(w, "Failed to delete page: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if any rows were affected
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "Failed to confirm deletion: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}

		// Return success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Page deleted successfully"))
	}
}
