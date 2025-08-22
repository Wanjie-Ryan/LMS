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





