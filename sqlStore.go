package main

import (
	"database/sql"
)

func main() {
	// @NOTE: the real connection is not required for tests
	db, err := sql.Open("mysql", "root@/blog")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Now we can use this db throughout our code
}

// Appropriate fields required for your application
type User struct {
	ID        	int64  `json:"id"`
	FirstName 	string `json:"firstName"`
	LastName  	string `json:"lastName"`
}

// SQL Store interface allowing you to abstract the sql client
type MySQLStore struct {
	Client *sql.DB
}

// Constructs and returns a pointer to a MySQLStore struct
func NewMySQLStore(db *sql.DB) *MySQLStore {
	if db != nil {
		return &MySQLStore{
			Client: db,
		}
	}
	return nil
}


//GetByID returns the User with the given ID
func (s *MySQLStore) GetByID(id int64) (*User, error) {
	// Process request and return correct output

	return nil, nil
}


//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (s *MySQLStore) Insert(user *User) (*User, error) {
	// Process request and return correct output

	return nil, nil
}



