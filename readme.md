**go.mod file** tells Go which packages the project will use. its like a package.json file. it lists dependencies and their version
**go.sum file** is automatically generated whenever you add go packages using go get. It provides security guarantees that the dependencies haven’t been modified
**go mod tidy** If you ever want to clean up unused deps

**LIBRARY MANAGEMENT SYSTEM GUIDE**

# Entities

User -> id, name, email, password, role (admin or normal user)
Books -> id, title, author, stock
Borrows -> id, user_id (who borrowed -> fk to users), book_id (which book was borrowed -> fk to books), borrow_date, return_date (can be null if not yet returned)

# Admin

- manages books, monitors users, and updates stock

# Normal user

- Registers, logs in, borrows and returns books.

# High level flow

- Users sign up and log in -> uses JWT token for auth
- Admin can add, update, delete books.
- User can borrow available books, if stock is > 0.
- when borrowing, stock decreases by 1 vice versa when you return a book.
- System keeps history of who borrowed what and when.

**Relationships Btn Entities**

1. User to Borrow

# 1 to Many

- one user can borrow many books over time.
- Each borrow record belongs to only one user.

2. Book to Borrow

# 1 to many

- One book can appear in many borrow records (different users borrowing at different times)
- Each borrow record belongs to one book.

3. User to Book (many to many) (indirectly via borrow)

- user can borrow many books
- A book can be borrowed by many users.
- This r/shp will be modeled by borrows table.

**ERD REPRESENTATION**
![alt text](image.png)

**Difference Between log.default().Println() and fmt.Println()**

# fmt

- fmt belongs to the **fmt** package.
- it prints text to the console.
- has no timestamps, no logs, just raw text.
- example --> Connected to database

# log.default()

- Belongs to the log package
- log.Default returns the default logger instance in GO.
- it includes timestamps and supports logging with more structured info.
- example --> 2025/08/21 10:30:00 Connected to database

# e.Logger.Fatal

- this is Echo's logger instance.
- More advanced than Go's standard log, because it supports structured logging
- .Fatal logs the message at Fatal level, and after logging it terminates the program.
- outcome is sth like this::-> {"time":"2025-08-21T10:55:02Z","level":"FATAL","msg":"Error loading .env file","error":"open .env: no such file or directory"}

**fmt.Sprintf**

- formats a string and returns it

# Running a Go project

- go run main.go

# Difference btn **appAddress := fmt.Sprintf("localhost:%s", port)** and **appAddress :=fmt.Sprintf(":%s", port)**

- The first one explicitly binds the server to localhost only.
- The server is only accessible at htt://localhost:8080, but not via other n/w interfaces like LAN IP 127.0.0.1 or 0.0.0.0

- The second one is a shorthand for all available interfaces (like 0.0.0.0:8080) for IPV4 and [::]:8080 for IPV6
- The server is reachable not only at localhost:8080, but also from other devices on your n/w if firewall allows.

- In go fields must start with an uppercase letter to be exported (visible to GORM and JSON)

# INDEXING

- without and index, DB would do a full table scan.
- starts at row 1 and checks every email until it finds rwanjie.
- Thats O(n) time to check upto 1m records.

**With an Index**

- The DB creates a B-Tree (balanced search tree)
- In this tree - each node contains a range of sorted email values. Branch 1 contains emails from A TO D, Branch 2 contains emails from E TO H like that...
- The branches let the DB jump directly to the right section.
  **Analogy**
- Think of it like a dictionary - if you're looking for Rwanjie, you don't start from "Aaron"
  **A B-Tree works like looking for a word in a dictionary**

# Maps

- A map in go is just a key-value store
  **map[KeyType]ValueType**
- map[string]string -> keys are strings, values are strings
  **Example**
  return c.Json(http.statusBadRequest, map[string]string{
  "error":"Invalid request"
  })
- key is string, value is string

# Installing Redis

**go get github.com/redis/go-redis/v9**

- Data in redis is tored in key value pairs

# Installing Swagger for documentation

**go install github.com/swaggo/swag/cmd/swag@latest**

- the above is installed only once per machine

**go get github.com/swaggo/echo-swagger@latest**
**go get github.com/swaggo/swag@latest**

- From your module root where go.mod lives, run:
  **swag init -g cmd/api/main.go**
- -g points to the entry file containing the top-level annotations (main.go)
- This generates a docs/ package
- Everytime you change annotations or models, **swag init** again

# CI/CD Workflow

**Continuous Integration**

- Automating the process of testing and integrating code into a shared repo.
- Everytime you push to Github, a CI pipeline can:

1. Run unit tests
2. Check fromatting/linting
3. Build your Go project
4. Alter you if sth breaks

**The goal is to catch bugs early, ensure every commit is production ready**

# Continuous Deliver/Deployment

- This is about automating what happens after CI passes.

1. Deploy to a testing/staging server
2. Trigger an automated deployment to prod.
3. Build and push a docker image
4. Notify slack or discord
5. Tag releases

**Move from a working commit, to a delivered product with manual effort**

# Github Actions

- Github Actions is Github's native CI/CD tool. It uses workflows (writen in YAML) to define automated jobs that trigger events like:

1. push. 2. pull requests 3. release 4. manual dispatch

- workflows are defined in
  **.github/workflows/name.yml**

# Workflows

- A workflow is a YAML file, that contains one or more jobs.
- Each job contains steps(run commands, actions, scripts)
- All triggered by a specific event. (on:)

- workflow is like a recipe for Github actions

1. it defines when it run/triggers
2. What env to use (eg. ubuntu-latest)
3. What steps to take (eg. checkout code, build, test)

- Everything from CI to CD is in this file.

# Example

---

on: [push]

jobs:
build:
runs-on: ubuntu-latest
steps: - uses: actions/checkout@v3 - name: Set up Go
uses: actions/setup-go@v4
with:
go-version: '1.21' - name: Run Tests
run: go test ./...

---

# Example (in the .yml file for reference)

- **on** Triggers the workflow when you push to main or dev or even when you open a PR targeting main or dev.
- **jobs** Defines the job called **build-test** to run on an Ubuntu runner.
- **steps**

1. **actions/checkout** - Gets your repo code.
2. **setup-go** - Installs Go in the runner
3. **go mod tidy** - cleans up go.mod and go.sum
4. **go-build** - checks for compilation issues.
5. **go-test** - runs your unit/integration tests.
6. **go-vet** - finds suspicious code.
7. **gofmt** - Ensures all Go files are properly formatted.

- If go test or go build fails, the workflow will stop and mark the build as failed.

**gofmt**

- Go's built in formatter -- its the Prettier of Go.
  **gofmt -l .**
- shows unformatted files
  **gofmt -w .**
- Auto formats everything
- -w: write the formatted output back to the files.
- . format the current directory recursively.
