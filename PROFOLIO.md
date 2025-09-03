# Portfolio Overview: My Library

---

## 🚀 Project Summary

A web-based personal reading tracker that integrates with the [Open Library API](https://openlibrary.org/developers/api) to provide:

* **Search & Explore** – Discover books and authors with detailed metadata.
* **Personal Library** – Organize books into Reading, To Read, Finished, and Favorites.
* **User Accounts** – Register, log in, and manage personal libraries.
* **Dynamic UI** – Built with HTMX for seamless, partial content updates.
* **Lightweight Backend** – Powered by Go + SQLite for efficiency and simplicity.

---

## 🧠 Technologies Used

* **Backend:** Go 1.23+, Gin, SQLx, Gorilla Sessions
* **Frontend:** HTML, Bootstrap 5, HTMX
* **Database:** SQLite
* **Deployment:** Docker, Docker Compose
* **Data Source:** Open Library API

---

## 🧹 Key Components Implemented

**Authentication & User Management**

* Secure password hashing (`bcrypt`) and session handling.
* Profile updates & password changes with validation.

**Book Management**

* Handlers for adding, moving, and removing books across categories.
* Recommendation engine (`recomend/`) for suggested books.

**Search System**

* Search by title, author, or subject.
* HTMX-powered results with dynamic book detail pages.

**HTMX Swap Integration**

* Modular handlers for partial HTML fragments (`htmxSwap/`).
* Enables SPA-like interactivity without a heavy JS framework.

---

## 📊 Architecture & Design Documentation

* **Layered Design**

  * **Presentation Layer**: Static HTML, Bootstrap, HTMX fragments.
  * **Business Logic**: Modular Go packages (`author`, `search`, `move`, `user`).
  * **Persistence Layer**: SQLite via `sqlx`.

* **Patterns & Practices**

  * **Modular Package Organization** – Each domain has a dedicated folder.
  * **SRP (Single Responsibility Principle)** – Each handler focuses on one action (e.g., `addToLibrary.go`, `favoriteAuthor.go`).
  * **Composition over Inheritance** – Go idioms applied for clean handler composition.
  * **RESTful Endpoints** – Clear separation of resources (user, search, books).

---

## 📊 Sample Workflow

1. User signs up and logs in.
2. Searches for a book (e.g., "The Hobbit").
3. Views book details, then adds it to "Want to Read".
4. Later moves it to "Reading" or "Finished".
5. Profile updates dynamically via HTMX swaps.

---

## 🧪 Testing Strategy

* Unit tests for core handlers (`search/book_test.go`).
* Manual testing of HTMX fragment rendering.
* SQLite database reset for clean test runs.

---

## 💬 What I Learned

* How to integrate **HTMX** for dynamic UI updates without full-page reloads.
* Best practices for structuring a **Go backend with modular packages**.
* Working with **Open Library API** for real-world data.
* Combining **SQLite + SQLx** for lightweight, production-ready persistence.
* Containerizing a Go application with **Docker & Docker Compose**.

---

## 🌟 Resume Summary

Developed a full-stack **Online Library** web app using Go, HTMX, and SQLite:

* Built modular backend with RESTful endpoints for books, search, and users.
* Integrated with Open Library API for live book/author data.
* Implemented authentication, password hashing, and profile management.
* Designed dynamic frontend using **HTMX** for seamless updates.
* Containerized with Docker for reproducible deployment.
* Stack: **Go, HTMX, Bootstrap, SQLite, Docker**.

