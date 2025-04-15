# Gitbook - Backend

Gitbook is a self-hosted code management platform designed to organize, trach an analyze your software repositories. This repository contains the backend service responsible for managing
project metadata, exposing APIs, and interfacing with the local git server.

This service is part of a larger system:

- [`gitbook-client`](https://github.com/Arihantawasthi/gitbook-client) - Frontend user interface.
- [`gitbook-crons`](https://github.com/Arihantawasthi/gitbook-crons) - Background jobs for calculating statistics across your repositories.

---

## ⚠️ Prerequisite: Git Repositories

This service **does not host Git repositories**. It reads from a directory on disk where your Git bare repositories already exist.

You are expected to:
- Clone/push repositories into a common local directory (e.g., `/home/user/repos`)
- Pass that directory path via the `REPO_BASE_PATH` environment variable

> If you're new to Git server setup (SSH, HTTP), refer to the [official Git documentation](https://git-scm.com/book/en/v2/Git-on-the-Server-Setting-Up-the-Server).

🧪 **Planned Feature**: Native Git hosting (with push/pull over SSH/HTTP) and repository creation via UI is on the roadmap. This will allow GitBook to act as a full Git server — stay tuned!

## 📌 Overview

This backend service provides:
- RESTful APIs for repositories, commits and stats.
- Integration with Git repositories stored locally.
- PostgreSQL-backed data persistence.
- Hooks for future features like post commit processing for statistics.

---

## 🧰 Tech Stack
- **Language**: Go (Golang)
- **Database**: PostgreSQL
- **Routing**: Go's standard `net/http`
- **Git Integration**: Local bare repositories + Git hook support

## 📡 API Endpoints
### 🗂 Repositories
| Method      | Endpoint                                         | Description
| :---        | :---                                             | :---
| `GET`       | `/api/v1/repos`                                  | List all tracked repositories
| `GET`       | `/api/v1/repo/{name}/{type}/metadata/{branch}/`  | Get file/folder structure of a repo at a branch
| `POST`      | `/api/v1/update-last-commit`                     | Manually trigger update of last commit data

### 📊 Statistics
| Method      | Endpoint          | Description
| :---        | :---              | :---
| `GET`       | `/api/v1/stats`   | Fetch overall repositories statistics (repos, commits, lines, etc.)

### 📜 Commits
| Method      | Endpoint                              | Description
| :---        | :---                                  | :---
| `GET`       | `/api/v1/repo/logs/{name}/{branch}`   | Get commit history for a branch
| `GET`       | `/api/v1/repo/commit/{name}/{hash}`   | Get file/folder structure of a repo at a branch
| `GET`       | `/api/v1/repo/{name}/{file}`          | Get commit history for a specific file

The `{type}` param in `/metadata/` is either `tree` (for directories) or `blob` (for files).

## 📁 Project Structure
```bash
gitbook
├── app
│   ├── handler
│   │   ├── comm_handler.go
│   │   └── repo_handler.go
│   ├── router.go
│   ├── services
│   │   ├── comm_service.go
│   │   └── repo_service.go
│   ├── storage
│   │   ├── connect.go
│   │   └── queries.go
│   └── types
│       └── types.go
├── cmd
│   └── main.go
├── go.mod
├── go.sum
├── runserver.sh
├── runserver.sh.example
└── utils
    ├── helpers.go
    └── logger.go
```

## ⚙️ Local Development Setup

### 1. Clone the Repository

```
git clone https://github.com/Arihantawasthi/gitbook.git && cd gitbook
```

### 2. Configure `runserver.sh`
Copy the sample runserver file and fill in the necessary values:
```
cp runserver.sh.example runserver.sh
```

Update values like:
```bash
#!/bin/bash

export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=gitbook
export SERVER_PORT=8080

go run cmd/main.go
```
Make it executable:
```
chmod +x runserver.sh
```

### 3. Start PostgreSQL
Ensure Postgres is running and create the `gitbook` database:
```
createdb gitbook
```

Alternatively, run Postgres using Docker:
```
docker run --name gitbook-postgres -e POSTGRES_PASSWORD=yourpassword -p 5432:5432 -d postgres
```

### 4. Run the server
```
./runserver.sh
```

The API will be available at http://localhost:8000

### 🧠 Notes
- Git repositories are stored locally as bare repositories.
- This service doesn't manage Git push/pull - it assumes respositories are updated externally via your git server

## 🚀 Future Enhancements
The current version of GitBook focuses on parsing local Git repositories and exposing relevant metadata via APIs.
But there's more features I'm planning to implement. Here’s a sneak peek into the roadmap:
- **🔐 Role-Based Access Control**
- **🏗 Repository Management via UI**
- **📤 Git Server Integration**
- **🤝 Collaboration Features**
