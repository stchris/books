package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	_ "rsc.io/sqlite"
	"strconv"
	"strings"
)

var (
	usr, _ = user.Current()

	// DBPath is the path to the db
	DBPath = usr.HomeDir + "/.books/"

	// DBName is the db file name
	DBName = "books.db"
)

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

func insert(book *Book, db *sql.DB) int {
	_, err := db.Exec(
		`INSERT INTO books(title, author, isbn, comments) VALUES(?, ?, ?, ?)`,
		book.Title,
		book.Author,
		book.ISBN,
		book.Comments,
	)
	if err != nil {
		log.Printf("Error while Inserting: %s", err)
	}

	row := db.QueryRow("SELECT last_insert_rowid()")
	x := 0
	err = row.Scan(&x)
	if err != nil {
		log.Printf("Error while getting autoincrement value: %s", err)
	}
	return x
}

// get a Book slice if the title, author or comments contain the given query
func getBooks(query string, db *sql.DB) []Book {
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
		queryString = "SELECT id, title, author, isbn, comments FROM books"
	}
	rows, err := db.Query(queryString)
	if err != nil {
		log.Printf("Error while Selecting: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		book := new(Book)
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Comments)
		if err != nil {
			log.Fatalf("Error getting row data: %s", err)
		}
		books = append(books, *book)
	}

	return books
}

func getBookByID(id int, db *sql.DB) (*Book, error) {
	rows, err := db.Query(`SELECT * FROM books WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	book := Book{}
	err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Comments)
	return &book, err
}

func deleteBookByID(id int, db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM books WHERE id = ?`, id)
	return err
}

func initDb(dbPath string, dbName string) (*sql.DB, error) {
	os.Mkdir(dbPath, 0700)
	var dbFullPath = dbPath + dbName

	var db, err = sql.Open("sqlite3", dbFullPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS books(id INTEGER PRIMARY KEY AUTOINCREMENT, title VARCHAR(200), author VARCHAR(200), isbn VARCHAR(20), comments TEXT);")
	return db, err
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
	db, err := initDb(DBPath, DBName)
	if err != nil {
		log.Fatalf("Failed to initialize db connection: %v", err)
	}
	defer db.Close()
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		var books = getBooks("", db)
		json.NewEncoder(w).Encode(books)
	} else if r.Method == "POST" {
		r.ParseForm()
		title := r.PostFormValue("title")
		author := r.PostFormValue("author")
		isbn := r.PostFormValue("isbn")
		comments := r.PostFormValue("comments")
		book := Book{Title: title, Author: author, ISBN: isbn, Comments: comments}
		insert(&book, db)
		http.Redirect(w, r, "/", http.StatusFound)
	}
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

	var db, err = initDb(DBPath, DBName)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	if command == "ls" {
		query := ""
		if len(subArgs) > 0 {
			query = strings.Join(subArgs, " ")
		}
		books := getBooks(query, db)
		for _, b := range books {
			fmt.Printf("%v\n", b)
		}
	} else if command == "add" {
		author := prompt("Author: ")
		title := prompt("Title: ")
		comments := prompt("Comments: ")
		book := Book{0, title, author, "", comments}
		insert(&book, db)
	} else if command == "del" {
		idString := prompt("id: ")
		id, err := strconv.ParseInt(idString, 10, 0)
		if err != nil {
			log.Println("Invalid id (not a number)")
			os.Exit(1)
		}
		book, err := getBookByID(int(id), db)
		if err != nil {
			log.Println("Error fetching book with id ", idString)
			os.Exit(1)
		}
		var promptString = fmt.Sprintf("Confirm deleting of %v (y/N)? ", book)
		if strings.ToUpper(prompt(promptString)) == "Y" {
			log.Println("Deleting ", book)
			deleteBookByID(int(id), db)
		}
	} else if command == "web" {
		http.Handle("/", http.FileServer(http.Dir("web/")))
		http.HandleFunc("/api/book", webAPIBook)
		var url = "0.0.0.0:8765"
		log.Println("Web server listening at http://" + url)
		log.Println("Press ^C to stop")
		http.ListenAndServe(url, nil)
	} else if command == "help" {
		printUsage()
	}
}
