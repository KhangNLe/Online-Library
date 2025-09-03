# 📚 My Library

A web-based application for discovering, tracking, and organizing books using the [Open Library API](https://openlibrary.org/developers/api).
It allows readers to search for books, build a personal library, and track reading progress with a clean, dynamic UI.

---

## 📜 Features

* **Book Search** – Look up books and authors via the Open Library API.
* **Personal Library** – Save books into categories:

  * Currently Reading
  * Planning / To Read
  * Finished
  * Favorites
* **Dynamic Updates** – Powered by [HTMX](https://htmx.org/), enabling partial reloads without a full page refresh.
* **Progress Tracking** – Organize and revisit your reading journey.
* **Authentication** – User accounts with login & signup functionality.
* **Profile Management** – Update passwords, profile details, and preferences.

---

## 🛠 Tech Stack

* **Frontend:** HTML, Bootstrap 5, HTMX
* **Backend:** Go 1.23+, Gin, SQLx, Gorilla Sessions
* **Database:** SQLite
* **Deployment:** Docker, Docker Compose
* **Data Source:** [Open Library API](https://openlibrary.org/developers/api)

---

## 📦 Dependencies

Defined in [`go.mod`](go.mod):

* `github.com/gin-gonic/gin` – HTTP web framework
* `github.com/jmoiron/sqlx` – SQL utilities
* `github.com/mattn/go-sqlite3` – SQLite driver
* `github.com/gorilla/sessions` – Session management
* `github.com/go-playground/validator/v10` – Request validation
* `github.com/bytedance/sonic` – High-performance JSON serialization
* *(and more, see `go.mod` for full list)*

---

## 🚀 Getting Started

### Prerequisites

* Go 1.23+
* Docker & Docker Compose

### Clone the Repository

```bash
git clone https://github.com/ctrl-MizteRy/Online-Library.git
cd Online-Library
```

### Run with Docker

```bash
docker-compose up --build
```

The app will be available at:
👉 **[http://localhost:8080](http://localhost:6969)**

---

## 📂 Project Structure

```
.
├── clientSide/       # Frontend static files (HTML, CSS, JS)
├── server/           # Backend Go modules
│   ├── author/       # Author-related endpoints
│   ├── homepage/     # Home & book listing
│   ├── htmxSwap/     # HTMX fragment rendering
│   ├── loginsignup/  # Authentication & password hashing
│   ├── move/         # Moving books between categories
│   ├── mybook/       # User's personal library handlers
│   ├── recomend/     # Recommendation system
│   ├── search/       # Search books, subjects, and details
│   └── user/         # User profile & account management
├── library.db        # SQLite database
├── main.go           # Application entry point
├── go.mod / go.sum   # Dependencies
├── Dockerfile        # Build instructions
└── docker-compose.yml# Deployment config
```

---

## 🧪 Testing

Some modules include Go test files (e.g., `search/book_test.go`).
Run all tests with:

```bash
go test ./...
```

---

## 📄 License

This project is licensed under the MIT License – see the [LICENSE](LICENSE) file for details.

