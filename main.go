package main

import (
	"bufio"
	"database/sql"
	"os"
	"strings"

	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Book struct {
	title  string
	author string
}

func isConnected() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "*******",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "for_libs",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

}

func addBook(book Book) {
	for_id := db.QueryRow("SELECT * FROM library ORDER BY id DESC LIMIT 1")
	var id int
	var title, author string

	err := for_id.Scan(&id, &title, &author)

	if err != nil {
		log.Fatalf("Error - %v", err)
	}

	res, err := db.Exec("INSERT INTO library (id, title, author) VALUES (?, ?, ?)", id+1, book.title, book.author)
	if err != nil {
		log.Fatal("Could not execute the command")
	}

	_, err = res.LastInsertId()
	if err != nil {
		log.Fatal("Could not get last id")
	}

	fmt.Printf("The book %v by %v was added with id %v", title, author, id+1)
}

func removeBook(book Book) {
	_, err := db.Exec("DELETE FROM library WHERE title=?", book.title)
	if err != nil {
		log.Fatal("Could not execute command")
	}

	fmt.Println("Successfuly removed")
}

func checkBook(book Book) bool {
	row := db.QueryRow("SELECT * FROM library WHERE title=?", book.title)

	var id, title, author string
	err := row.Scan(&id, &title, &author)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Fatalf("Error - %v", err)
	}

	return true
}

func main() {
	isConnected()

	for {
		readerBook := bufio.NewReader(os.Stdin)
		readerChoice := bufio.NewReader(os.Stdin)
		readerAuthor := bufio.NewReader(os.Stdin)

		var new_book Book

		fmt.Println("Add, check or remove the book? To quit type 'quit'")

		choice, err := readerChoice.ReadString('\n')
		if err != nil {
			log.Fatal("Could not read the input")
		}
		if choice == "quit\n" {
			os.Exit(1)
		}

		fmt.Print("What book ?")

		new_book.title, err = readerBook.ReadString('\n')
		if err != nil {
			log.Fatal("Could not read the book")
		}

		fmt.Print("Write author:")

		new_book.author, err = readerAuthor.ReadString('\n')
		if err != nil {
			log.Fatal("Could not read the author")
		}

		new_book.title = strings.TrimSuffix(new_book.title, "\n")
		new_book.author = strings.TrimSuffix(new_book.author, "\n")

		switch choice {
		case "Add\n":
			addBook(new_book)
		case "Remove\n":
			removeBook(new_book)
		case "Check\n":
			fmt.Println(checkBook(new_book))
		default:
			log.Fatal("WRONG")
		}

	}

}
