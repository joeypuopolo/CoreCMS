package main

import (
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"cms/db"
	"cms/handlers"
)

// curl -X POST -H "Content-Type: application/json" -d '{"content":"<h1>Example!!</h1>","title":"Main Headline","description": "Example." "active":1}' http://localhost:8080/code_blocks
// curl -X POST -H "Content-Type: application/json" -d '{"title":"Main Template", "parent_template_id": -1, "active": 1}' http://localhost:8080/templates
// curl -X POST -H "Content-Type: application/json" -d '{"title":"Homepage","url":"/home", "hidden": -1, "active": 1, "parent_page": -1, "template_id": 1}' http://localhost:8080/pages
// curl -X POST -H "Content-Type: application/json" -d '{"title":"Homepage","url":"/home", "hidden": -1, "active": 1, "parent_page": -1, "template_id": 1}' http://localhost:8080/pages

func main() {
	database := db.Connect()
	defer database.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Println("Path: " + path)
		// Get the file extension

		// Detect the MIME type based on file extension
		mimeType := mime.TypeByExtension(filepath.Ext(path))
		if mimeType != "" {
			w.Header().Set("Content-Type", mimeType)
		}

		http.ServeFile(w, r, "./front-end/index.html")
	})

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	path := r.URL.Path
	// 	if path == "/" {
	// 		path = "/index.html" // Default to index.html
	// 	}

	// // Get the file extension
	// ext := strings.ToLower(filepath.Ext(path))
	// switch ext {
	// case ".js":
	// 	w.Header().Set("Content-Type", "application/javascript")
	// case ".css":
	// 	w.Header().Set("Content-Type", "text/css")
	// case ".html":
	// 	w.Header().Set("Content-Type", "text/html")
	// }
	// 	log.Println("Path" + path)
	// 	http.ServeFile(w, r, "./front-end"+path)
	// })

	// Pages Routes
	r.Route("/pages", func(r chi.Router) {
		// Create
		r.Post("/", handlers.CreatePage(database))
		// Read
		r.Get("/", handlers.GetPages(database))
		r.Get("/{pageID}", handlers.RenderPage(database))
		// Update
		r.Patch("/{pageID}", handlers.UpdatePage(database))
		// Delete
		r.Delete("/{pageID}", handlers.DeletePage(database))
	})

	// Templates Routes
	r.Route("/templates", func(r chi.Router) {
		// Create
		r.Post("/", handlers.CreateTemplate(database))
		r.Post("/duplicate/{templateID}", handlers.DuplicateTemplate(database))
		// Read
		r.Get("/", handlers.GetTemplates(database))
		r.Get("/{templateID}", handlers.GetTemplate(database))
		// Update
		r.Patch("/{templateID}/name", handlers.UpdateTemplate(database))
		// Delete
		r.Delete("/{templateID}", handlers.DeleteTemplate(database))
	})

	// Code Blocks Routes
	r.Route("/code_blocks", func(r chi.Router) {
		// Create
		r.Post("/", handlers.CreateCodeBlock(database))
		// Read
		r.Get("/", handlers.GetCodeBlocks(database))
		r.Get("/{codeBlockID}", handlers.GetCodeBlock(database))
		// Update
		r.Patch("/{codeBlockID}", handlers.UpdateCodeBlock(database))
		// Delete
		r.Delete("/{codeBlockID}", handlers.DeleteCodeBlock(database))
	})

	// Settings Routes (add later)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
