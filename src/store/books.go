package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

type TBooks struct {
	Book_id    int       `json:"book_id"`
	Book_autor string    `json:"book_autor"`
	Book_name  string    `json:"book_name"`
	Book_name2 string    `json:"book_name2"`
	Code       string    `json:"code"`
	Book_date  time.Time `json:"book_date"`
	Book_lang  string    `json:"book_lang"`
	Rem        string    `json:"rem"`
	Fzip_name  string    `json:"fzip_name"`
}

type TBookOne struct {
	Book_id    int    `json:"book_id"`
	Book_autor string `json:"book_autor"`
	Book_name  string `json:"book_name"`
	Fzip_name  string `json:"fzip_name"`
}

func (pg *postgres) GetBooksAutor(autor string) ([]TBooks, error) {
	rows, err := pg.db.Query(context.Background(), "SELECT book_id,book_autor,book_name,book_name2,code,book_date,book_lang,CASE WHEN book_rem IS NULL THEN '-' ELSE book_rem END AS rem,fzip_name FROM books LEFT JOIN fzip ON books.id_fzip = fzip.fzip_id WHERE LOWER(book_autor) LIKE CONCAT('%',$1::text,'%') LIMIT 500", autor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[TBooks])
	if err != nil {
		log.Println(err.Error())
	}

	return books, nil
}

func (pg *postgres) GetBooksName(name string) ([]TBooks, error) {
	rows, err := pg.db.Query(context.Background(), "SELECT book_id,book_autor,book_name,book_name2,code,book_date,book_lang,CASE WHEN book_rem IS NULL THEN '-' ELSE book_rem END AS rem,fzip_name FROM books LEFT JOIN fzip ON books.id_fzip = fzip.fzip_id WHERE LOWER(book_name) LIKE CONCAT('%',$1::text,'%') LIMIT 500", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[TBooks])
	if err != nil {
		log.Println(err.Error())
	}
	return books, nil
}

func (pg *postgres) GetBooks(autor, name string) ([]TBooks, error) {
	rows, err := pg.db.Query(context.Background(), "SELECT book_id,book_autor,book_name,book_name2,code,book_date,book_lang,CASE WHEN book_rem IS NULL THEN '-' ELSE book_rem END AS rem,fzip_name FROM books LEFT JOIN fzip ON books.id_fzip = fzip.fzip_id WHERE LOWER(book_autor) LIKE CONCAT('%',$1::text,'%') AND LOWER(book_name) LIKE CONCAT('%',$2::text,'%') LIMIT 500", autor, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[TBooks])
	if err != nil {
		log.Println(err.Error())
	}
	return books, nil
}

func (pg *postgres) GetOneBook(code string) TBookOne {
	row := pg.db.QueryRow(context.Background(), "SELECT book_id,book_autor,book_name,fzip_name FROM books LEFT JOIN fzip ON books.id_fzip = fzip.fzip_id WHERE code = $1", code)

	var book TBookOne
	if err := row.Scan(&book.Book_id, &book.Book_autor, &book.Book_name, &book.Fzip_name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("1:", err.Error(), ErrNotExist)
		}
		log.Println("2:", err.Error())
	}
	// println(book.Book_name)
	return book
}
