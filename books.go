package main

import (
	"bufio"
	"code.google.com/p/gosqlite/sqlite"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
)

var usr, _ = user.Current()
// path to db
var DBPATH = usr.HomeDir + "/.books/"
// db file name
var DBNAME = "books.db"

// Book bundles data related to a book
type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	ISBN     string `json:"isbn"`
	Comments string `json:"comments"`
}

func (book Book) String() string {
	return fmt.Sprintf(
		"{#%v: '%v' by '%v', comments: %v}",
		book.ID, book.Title, book.Author, book.Comments,
	)
}

func insert(book *Book, conn *sqlite.Conn) int {
	insertSql := fmt.Sprintf(
		`INSERT INTO books(title, author, isbn, comments) VALUES('%v', '%v', '%v', '%v');`,
		book.Title,
		book.Author,
		book.ISBN,
		book.Comments,
	)

	err := conn.Exec(insertSql)
	if err != nil {
		fmt.Printf("Error while Inserting: %s", err)
	}

	selectStmt, err := conn.Prepare("select last_insert_rowid();")
	if err != nil {
		fmt.Printf("Error while getting autoincrement value: %s", err)
	}

x := 0
	if selectStmt.Next() {
		selectStmt.Scan(&x)
	}

	return x
}

// try to parse a Book from a DB statement
func getBookFromStmt(stmt *sqlite.Stmt) *Book {
	book := new(Book)

	err := stmt.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Comments)
	if err != nil {
		fmt.Printf("Error while getting row data: %s\n", err)
		os.Exit(1)
	}

	return book
}

// get a Book slice if the title, author or comments contain the given query
func getBooks(query string, conn *sqlite.Conn) []Book {
	var books []Book
	var queryString string
	if query != "" {
		queryString = fmt.Sprintf(
			`SELECT * FROM books WHERE title LIKE "%%%v%%"
		UNION
		SELECT * FROM books WHERE author LIKE "%%%v%%"
		UNION
		SELECT * FROM books WHERE comments LIKE "%%%v%%"`,
			query,
			query,
			query)
	} else {
		queryString = "SELECT * FROM books"
	}
	selectStmt, err := conn.Prepare(queryString)
	err = selectStmt.Exec()
	if err != nil {
		fmt.Printf("Error while Selecting: %v", err)
	}

	for selectStmt.Next() {
		book := getBookFromStmt(selectStmt)
		books = append(books, *book)
	}

	return books
}

func getBookByID(id int, conn *sqlite.Conn) (*Book, error) {
	var book *Book
	var queryString = fmt.Sprintf(`SELECT * FROM books WHERE id = %v`, id)
	stmt, err := conn.Prepare(queryString)
	err = stmt.Exec()
	if err != nil {
		return book, err
	}
	stmt.Next()
	book = getBookFromStmt(stmt)
	return book, nil
}

func deleteBookByID(id int, conn *sqlite.Conn) error {
	queryString := fmt.Sprintf(`DELETE FROM books WHERE id = %v`, id)
	stmt, err := conn.Prepare(queryString)
	err = stmt.Exec()
	if err != nil {
		return err
	}
	stmt.Next()
	return nil
}

func initDb(dbPath string, dbName string) (*sqlite.Conn, error) {
	os.Mkdir(dbPath, 0700)
	var db = dbPath + dbName

	var conn, dberr = sqlite.Open(db)
	if dberr != nil {
		return nil, dberr
	}

	conn.Exec("CREATE TABLE books(id INTEGER PRIMARY KEY AUTOINCREMENT, title VARCHAR(200), author VARCHAR(200), isbn VARCHAR(20), comments TEXT);")

	return conn, nil
}

func prompt(text string) string {
	fmt.Print(text)
	bio := bufio.NewReader(os.Stdin)
	line, _, err := bio.ReadLine()
	if err != nil {
		fmt.Println("Error during input: ", err)
	}
	return string(line)
}

func webAPIBook(w http.ResponseWriter, r *http.Request) {
	conn, _ := initDb(DBPATH, DBNAME)
	defer conn.Close()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var books = getBooks("", conn)
	json.NewEncoder(w).Encode(books)
}

func printUsage() {
	fmt.Println("USAGE: books COMMAND argument1 argument2 ...")
	fmt.Println("Available commands:")
	fmt.Println("\tls - list books, pass search terms as arguments")
	fmt.Println("\tadd - add book (prompts for title, author, comments)")
	fmt.Println("\tdel - delete book (prompts for id)")
	fmt.Println("\tweb - starts the built-in web server on port 8765")
	fmt.Println("\thelp - display this text")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}
	command := args[0]
	subArgs := args[1:]

	var conn, err = initDb(DBPATH, DBNAME)
	if err != nil {
		fmt.Println("Error initializing database ", err)
		os.Exit(1)
	}
	defer conn.Close()

	if command == "ls" {
		query := ""
		if len(subArgs) > 0 {
			query = strings.Join(subArgs, " ")
		}
		books := getBooks(query, conn)
		for _, b := range books {
			fmt.Printf("%v\n", b)
		}
	} else if command == "add" {
		author := prompt("Author: ")
		title := prompt("Title: ")
		comments := prompt("Comments: ")
		book := Book{0, title, author, "", comments}
		insert(&book, conn)
	} else if command == "del" {
		idString := prompt("id: ")
		id, err := strconv.ParseInt(idString, 10, 0)
		if err != nil {
			fmt.Println("Invalid id (not a number)")
			os.Exit(1)
		}
		book, err := getBookByID(int(id), conn)
		if err != nil {
			fmt.Println("Error fetching book with id ", idString)
			os.Exit(1)
		}
		var promptString = fmt.Sprintf("Confirm deleting of %v (y/N)? ", book)
		if strings.ToUpper(prompt(promptString)) == "Y" {
			fmt.Println("Deleting ", book)
			deleteBookByID(int(id), conn)
		}
	} else if command == "web" {
		http.Handle("/", http.FileServer(http.Dir("web/")))
		http.HandleFunc("/api/book", webAPIBook)
		var url = "127.0.0.1:8765"
		fmt.Println("Web server listening at http://" + url)
		fmt.Println("Press ^C to stop")
		http.ListenAndServe(url, nil)
	} else if command == "help" {
		printUsage()
	}
}
