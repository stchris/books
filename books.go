package main

import (
	"code.google.com/p/gosqlite/sqlite"
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

// Book bundles data related to a book
type Book struct {
	ID       int
	Title    string
	Author   string
	ISBN     string
	Comments string
}

var conn = initDb()

func (book Book) String() string {
	return fmt.Sprintf(
		"{#%v: '%v' by '%v', comments: %v}",
		book.ID, book.Title, book.Author, book.Comments,
	)
}

func insert(book *Book) int {
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

func getBooks(query string) []Book {
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
		book := new(Book)

		err = selectStmt.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Comments)
		if err != nil {
			fmt.Printf("Error while getting row data: %s\n", err)
			os.Exit(1)
		}

		books = append(books, *book)
	}

	return books
}

func getBookByID(id int) Book {
	var book Book
	var queryString = fmt.Sprintf(`SELECT * FROM books WHERE id = %v`, id)
	stmt, err := conn.Prepare(queryString)
	err = stmt.Exec()
	if err != nil {
		fmt.Println("Error getting book with id ", id, ": ", err)
		os.Exit(1)
	}
	stmt.Next()
	err = stmt.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Comments)
	if err != nil {
		fmt.Printf("Error while getting row data: %s\n", err)
		os.Exit(1)
	}
	return book
}

func deleteBookByID(id int) {
	queryString := fmt.Sprintf(`DELETE FROM books WHERE id = %v`, id)
	stmt, err := conn.Prepare(queryString)
	err = stmt.Exec()
	if err != nil {
		fmt.Printf("Error while deleting Book")
		os.Exit(1)
	}
	stmt.Next()
}

func initDb() *sqlite.Conn {
	var usr, err = user.Current()
	if err != nil {
		fmt.Println("Error getting user's home dir: ", err)
	}
	os.Mkdir(usr.HomeDir + "/.books", 0700)
	var db = usr.HomeDir + "/.books/books.db"

	var conn, dberr = sqlite.Open(db)
	if dberr != nil {
		fmt.Println("Error opening the database file: ", dberr)
		os.Exit(1)
	}

	conn.Exec("CREATE TABLE books(id INTEGER PRIMARY KEY AUTOINCREMENT, title VARCHAR(200), author VARCHAR(200), isbn VARCHAR(20), comments TEXT);")

	return conn
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

func printUsage() {
	fmt.Println("USAGE: books COMMAND argument1 argument2 ...")
	fmt.Println("Available commands:")
	fmt.Println("\tls - list books, pass search terms as arguments")
	fmt.Println("\tadd - add book (prompts for title, author, comments")
	fmt.Println("\tdel - delete book (prompts for id)")
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

	initDb()
	defer conn.Close()

	if command == "ls" {
		query := ""
		if len(subArgs) > 0 {
			query = strings.Join(subArgs, " ")
		}
		books := getBooks(query)
		for _, b := range books {
			fmt.Printf("%v\n", b)
		}
	} else if command == "add" {
		author := prompt("Author: ")
		title := prompt("Title: ")
		comments := prompt("Comments: ")
		book := Book{0, title, author, "", comments}
		insert(&book)
	} else if command == "del" {
		idString := prompt("id: ")
		id, err := strconv.ParseInt(idString, 10, 0)
		if err != nil {
			fmt.Println("Invalid id (not a number)")
			os.Exit(1)
		}
		book := getBookByID(int(id))
		var promptString = fmt.Sprintf("Confirm deleting of %v (y/N)? ", book)
		if strings.ToUpper(prompt(promptString)) == "Y" {
			fmt.Println("Deleting ", book)
			deleteBookByID(int(id))
		}
	} else if command == "help" {
		printUsage()
	}
}
