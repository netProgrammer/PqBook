package main

import (
	_"github.com/lib/pq"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type Book struct  {
	isbn string
	title string
	author string
	price float32
}

var db *sql.DB

func init()  {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=user password=pass dbname=Stoe sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main()  {
	http.HandleFunc("/books", booksIndex)
	http.ListenAndServe(":3000", nil)
}

func booksIndex(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	defer rows.Close()

	bks := make([]*Book, 0)
	for rows.Next(){
		bk:= new(Book)
		err:= rows.Scan(&bk.isbn, &bk.title, &bk.author, &bk.price)

		if err != nil {
			log.Fatal(err)
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, bk := range bks {
		fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.isbn, bk.title, bk.author, bk.price)
	}
}