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
