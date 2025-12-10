# CoreCMS — A Database-Driven Content Management System in Go

CoreCMS is a custom content management system built with Go, Chi, and a database-driven templating architecture.
Its core features include dynamic page creation, template inheritance, nestable code blocks, and an inline-editable admin panel.
This project serves both as a functional CMS and an exercise in full-stack engineering.


## Key Features

- Database-driven pages: Create pages programmatically with stored routes, content types, and JSON content structures.
- Template system: All templates are structured for inheritance, and can be stored in the database for dynamic rendering.
- Nestable CodeBlocks: Code blocks editable by a web developer to build a website project. They can render inside templates or other code blocks.


## Admin Panel Features:

- CRUD for pages, templates, and code blocks
- Allow page ordering, naming, routing, hierarchy and visibility toggles
- Login/session management


## Why I Chose Go for This Project

I intentionally chose Go over more traditional web languages (like PHP, Node, Python etc.) for three reasons:

#### 1. Performance & Simplicity

Go is faster than the alternatives, and that was going to matter for this project. For a CMS that is serving many dynamic routes, dynamic pages, and dynamic HTML code, Go would offer scalability for any website size with less effect on render speed than other popular web.

#### 2. Exercise in a New Language

I wanted to deepen my backend engineering skills by working in a different language with familiar constructs. Go is very simple and I love their alternative of being "object oriented" using structs.

#### 3. Compiled Binary File (Protecting Source Code)

Because Go compiles to a self-contained binary, I can distribute this CMS to potential clients without exposing the source code.
This makes it possible to license out a custom CMS product while still protecting the underlying IP.

## Technologies Used

- Go (Golang)
- Chi router
- Go html/template engine
- SQL (Currently using SQLite, but in deployment can use PostgreSQL)
- Vanilla JS for admin UI
- go-session for login sessions
