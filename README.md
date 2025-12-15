# CoreCMS — A Database-Driven Content Management System in Go

CoreCMS is a custom content management system built with Go and Chi. Its architecture is inspired by all the CMS systems I've worked in, and so it reflects the way I think about websites. This code is not yet completed, but this README file is an attempt to explain all its current and future features. 

Additionally, this project was an exercise in practicing/learning a new language. I think if this were to go into actual use/production, this would best be written in PHP, that way we could easily run this on a low-cost shared hosting environment with default settings, which reflects the needs of typical website clients. However with this Go version, a great use case for this CMS could be hosting an app within a client's website as this CMS has that flexibility.


## Key Features

- Database-driven: All website code is stored in the database.
- Everything is template based: All templates, pages and code components are structured for inheritance, and they are stored in the database for dynamic rendering.
- Pages: Every page loads an associated template.
- Templates: Every template loads associated code components.
- Nestable CodeBlocks: Code components, called CodeBlocks, are editable by a web developer to build a website project. They render inside templates, pages or other code blocks. Website admin users have a UI to sort/add/remove CodeBlocks from the template and page views.
- Inherent unpublished system: The website is all served from a separate database table that stores all the compiled data and content for each page.


## Admin Panel Features:

The main admin panel has an easy UI that allows website admins to rename pages, sort them, set up page hierarchies, set pages to inactive to the public or hidden, customize the routing and swap page templates. For the customer's web developer, they can easily add and assemble code components to load on the website.



## Technologies Used

- Go (Golang)
- Chi router
- Go html/template engine
- SQL (Currently using SQLite, but in deployment can use PostgreSQL or MariaDB)
- Vanilla JS for admin UI
- go-session for login sessions
