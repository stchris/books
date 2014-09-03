package main

import (
    "code.google.com/p/gosqlite/sqlite"
    "os"
    "testing"
)

func queryAndExpect(query string, expected int, conn *sqlite.Conn, t *testing.T) {
    var bookCount = len(getBooks(query, conn))
    if bookCount != expected {
        t.Errorf("Expected %v books, found %v", expected, bookCount)
    }
}

func TestBasic(t *testing.T) {
    var conn, _ = initDb("./", "test.db")
    defer conn.Close()

    queryAndExpect("", 0, conn, t)

    var b = Book{}
    insert(&b, conn)

    queryAndExpect("", 1, conn, t)

    os.Remove("./test.db")
}

func TestQueries(t *testing.T) {
    var conn, _ = initDb("./", "test.db")
    defer conn.Close()

    var b1 = Book{ID:1, Title:"Foo", Author:"bar"}
    var b2 = Book{ID:2, Title:"Zebra", Author:"bar"}
    var b3 = Book{ID:3, Title:"Delta Spaces otherchars (subtitle goes here)", Author:"other"}
    insert(&b1, conn)
    insert(&b2, conn)
    insert(&b3, conn)

    queryAndExpect("", 3, conn, t)
    queryAndExpect("bar", 2, conn, t)
    queryAndExpect("other", 1, conn, t)
    queryAndExpect("a", 3, conn, t)
    queryAndExpect("NO such $tring", 0, conn, t)

    os.Remove("./test.db")
}
