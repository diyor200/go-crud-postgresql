package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "Your_username"
	password = "Your_password"
	dbname   = "_your_dbname"
)

type Book struct {
	Id     int
	Title  string
	Author string
	Isbn   string
	Price  float64
}

// ma'lumotlarni console ga chiqarish
func displayAll(db *sql.DB) {
	fmt.Println(queryAllBooks(db))
}

// ma'lumot qo'shish uchun
func addData() Book {
	book := Book{}
	var (
		title, author, isbn string
		price               float64
	)

	fmt.Print("New title:")
	fmt.Scan(&title)
	fmt.Print("New author:")
	fmt.Scan(&author)
	fmt.Print("New isbn:")
	fmt.Scan(&isbn)
	fmt.Println("New price:")
	fmt.Scan(&price)
	book.Title = title
	book.Author = author
	book.Isbn = isbn
	book.Price = price
	return book
}

func main() {
	var command int
	// Create the connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open the connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	// Create the book table
	err = createTable(db)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("--Bazaga ulandi--")
	for {
		fmt.Println("Komandalar:\n1 - qo'shish\n2 - o'qish,\n3 - o'zgartirish\n4 - o'chirish.\n\n0 - to'xtatish")
		fmt.Scan(&command)

		switch command {
		case 1:
			// Insert data into the book table
			book := addData()
			err = insertBook(db, book)
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			// Query all books from the book table
			books, err := queryAllBooks(db)
			if err != nil {
				fmt.Println(err)
			}

			// Print the books
			fmt.Println("Books:")
			for _, book := range books {
				fmt.Printf("ID: %d, Title: %s, Author: %s, ISBN: %s, Price: %.2f\n", book.Id, book.Title, book.Author, book.Isbn, book.Price)
			}
		case 3:
			updateBook(db)
		case 4:
			delete(db)
		case 0:
			fmt.Println("--Xayr!--")
			os.Exit(0)
		}
	}
}
func createTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS book (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		isbn TEXT NOT NULL,
		price DECIMAL(10, 2) NOT NULL
	)`)
	if err != nil {
		return err
	}
	return nil
}

// updateBook() -> updates data on db
func updateBook(db *sql.DB) {
	fmt.Println(">>>Siz update commandasini tanladingiz<<<")
	var id int
	fmt.Println(queryAllBooks(db))
	fmt.Println("Id sini tanlang: ")
	fmt.Scan(&id)

	data := addData()

	_, err := db.Exec("UPDATE book SET title=$1, author=$2, isbn=$3, price=$4 WHERE id=$5",
		data.Title, data.Author, data.Isbn, data.Price, id,
	)
	if err != nil {
		panic(err)
	}

	list, err := queryOneBook(db, id)
	if err != nil {
		panic(err)
	}
	for _, i := range list {
		fmt.Println(">>> O'zgartirilgan ma'lumot <<<\n",
			i.Id, i.Title, i.Author, i.Isbn, i.Price)
	}
}

// insertBook() -> inserts data
func insertBook(db *sql.DB, book Book) error {
	_, err := db.Exec("INSERT INTO book (title, author, isbn, price) VALUES ($1, $2, $3, $4)", book.Title, book.Author, book.Isbn, book.Price)
	if err != nil {
		return err
	}
	fmt.Println(">>>Ma'lumotlar kiritildi<<<")
	return nil
}

// delete() -> deletes data
func delete(db *sql.DB) {
	var id int
	displayAll(db)
	fmt.Println(">>>O'chirish uchun id ni kiriting:<<<")
	fmt.Scan(&id)

	_, err := db.Exec("DELETE FROM book WHERE id=$1", id)
	if err != nil {
		panic(err)
	}
	fmt.Println("Ma'lumot o'chirildi")

}

// queryOneBook() -> gets one data
func queryOneBook(db *sql.DB, id int) ([]Book, error) {
	rows, err := db.Query("SELECT * FROM book WHERE id=$1", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		book := Book{}
		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Isbn, &book.Price)
		if err != nil {
			panic(err)
		}
		books = append(books, book)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return books, nil
}

// queryAllBooks() -> gets all data
func queryAllBooks(db *sql.DB) ([]Book, error) {
	rows, err := db.Query("SELECT * FROM book")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		book := Book{}
		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Isbn, &book.Price)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return books, nil
}
