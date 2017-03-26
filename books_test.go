package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
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

func TestAPIgetOk(t *testing.T) {
	DBName = "testgetok.db"
	defer os.Remove(DBName)

	req, err := http.NewRequest("GET", "/api/book", nil)
	if err != nil {
		t.Fatalf("Error setting up http request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webAPIBook)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPIpostAndGet(t *testing.T) {
	DBName = "testpostandgetok.db"
	defer os.Remove(DBName)

	postReq := httptest.NewRequest("POST", "/api/book", nil)
	getReq := httptest.NewRequest("GET", "/api/book", nil)

	rr := httptest.NewRecorder()

	webAPIBook(rr, postReq)
	resp := rr.Result()
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	webAPIBook(rr, getReq)
	resp = rr.Result()
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPImethodFail(t *testing.T) {
	DBName = "testmethodfail.db"
	defer os.Remove(DBName)

	req, err := http.NewRequest("PATCH", "/api/book", nil)
	if err != nil {
		t.Fatalf("Error setting up http request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webAPIBook)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
