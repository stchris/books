package main

import (
	"database/sql"
	"os"
	_ "rsc.io/sqlite"
	"testing"
)

func queryAndExpect(query string, expected int, db *sql.DB, t *testing.T) {
	var bookCount = len(getBooks(query, db))
	if bookCount != expected {
		t.Errorf("Expected %v books, found %v", expected, bookCount)
	}
}

func TestBasic(t *testing.T) {
	var db, err = initDb("./", "test.db")
	if err != nil {
		t.Fatalf("Error on db init: %v", err)
	}
	defer db.Close()

	queryAndExpect("", 0, db, t)

	var b = Book{}
	insert(&b, db)

	queryAndExpect("", 1, db, t)

	os.Remove("./test.db")
}

func TestQueries(t *testing.T) {
	var db, err = initDb("./", "test.db")
	if err != nil {
		t.Fatalf("Error on db init: %v", err)
	}
	defer db.Close()

	var b1 = Book{ID: 1, Title: "Foo", Author: "bar"}
	var b2 = Book{ID: 2, Title: "Ze'bra", Author: "bar"}
	var b3 = Book{ID: 3, Title: "Delta Spaces otherchars (subtitle goes here)", Author: "other"}
	insert(&b1, db)
	insert(&b2, db)
	insert(&b3, db)

	queryAndExpect("", 3, db, t)
	queryAndExpect("bar", 2, db, t)
	queryAndExpect("other", 1, db, t)
	queryAndExpect("a", 3, db, t)
	queryAndExpect("NO such $tring", 0, db, t)

	os.Remove("./test.db")
}
