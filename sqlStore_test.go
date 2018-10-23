package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestMySQLStore_GetByID(t *testing.T) {
	//create a new sql mock
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating sql mock: %v", err)
	}

	//ensure it's closed at the end of the test
	defer db.Close()

	// Initialize a user struct we will use as a test variable
	expectedUser := &User{
		ID: 1,
		FirstName: "John",
		LastName: "Doe",
	}

	// Initialize a MySQLStore struct to allow us to interface with the SQL client
	store := NewMySQLStore(db)

	// Create a row with the appropriate fields in your SQL database
	// Add the actual values to the row
	row := sqlmock.NewRows([]string{"id", "firstname", "lastname"})
	row.AddRow(expectedUser.ID, expectedUser.FirstName, expectedUser.LastName)


	// Expecting a successful "query"
	// This tells our db to expect this query (id) as well as supply a certain response (row)
	mock.ExpectQuery("Select * From users Where id=?").
		WithArgs(expectedUser.ID).WillReturnRows(row)

	// Since we know our query is successful, we want to test whether there happens to be
	// any expected error that may occur.
	user, err := store.GetByID(expectedUser.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Again, since we are assuming that our query is successful, we can test for when our
	// function doesn't work as expected.
	if err == nil && !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("User queried does not match expected user")
	}

	// Expecting a unsuccessful "query"
	// Attempting to search by an id that doesn't exist. This would result in a
	// sql.ErrNoRows error
	mock.ExpectQuery(regexp.QuoteMeta("Select * From users Where id=?")).
		WithArgs(-1).WillReturnError(sql.ErrNoRows)

	// Since we are expecting an error here, we create a condition opposing that to see
	// if our GetById is working as expected
	if 	_, err = store.GetByID(-1); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", sql.ErrNoRows)
	}

	// Attempting to trigger a DBMS querying error
	queryingErr := fmt.Errorf("DBMS error when querying")
	mock.ExpectQuery("Select * From users Where id=?").
		WithArgs(expectedUser.ID).WillReturnError(queryingErr)

	if 	_, err = store.GetByID(expectedUser.ID); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", queryingErr)
	}

	// This attempts to check if there are any expectations that we haven't met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}

}

func TestMySQLStore_Insert(t *testing.T) {
	//create a new sql mock
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating sql mock: %v", err)
	}
	//ensure it's closed at the end of the test
	defer db.Close()

	// Initialize a user struct we will use as a test variable
	inputUser := &User{
		ID: 2,
		FirstName: "John",
		LastName: "Doe",
	}

	// Initialize a MySQLStore struct to allow us to interface with the SQL client
	store := NewMySQLStore(db)

	// This tells our db to expect an insert query with certain arguments with a certain
	// return result
	mock.ExpectExec("insert into users(id, firstname, lastname) values(?,?,?)").
		WithArgs(inputUser.ID, inputUser.FirstName, inputUser.LastName).
		WillReturnResult(sqlmock.NewResult(2, 1))

	user, err := store.Insert(inputUser)


	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if err == nil && !reflect.DeepEqual(user, inputUser) {
		t.Errorf("User returned does not match input user")
	}

	// Inserting an invalid user
	invalidUser := &User{
		-1,
		nil,
		"Doe",
	}
	insertErr := fmt.Errorf("Error executing INSERT operation")
	mock.ExpectExec("insert into users(id, firstname, lastname) values(?,?,?)").
		WithArgs(invalidUser.ID, invalidUser.FirstName, invalidUser.LastName).
		WillReturnError(insertErr)

	if 	_, err = store.Insert(invalidUser); err == nil {
		t.Errorf("Expected error: %v but recieved nil", insertErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}

}