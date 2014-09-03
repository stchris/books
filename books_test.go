package main

import (
    "os"
    "testing"
)

func TestBasic(t *testing.T) {
    var conn, _ = initDb("./", "test.db")
    var bookCount = len(getBooks("", conn))
    if  bookCount != 0 {
        t.Errorf("Expected 0 books, found %v", bookCount)
    }

    var b = Book{}
    insert(&b, conn)

    bookCount = len(getBooks("", conn))
    if bookCount != 1 {
        t.Errorf("Expected 1 books, found %v", bookCount)
    }

    os.Remove("./test.db")
}
