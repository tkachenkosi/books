package store

import (
	"bufio"
	"context"
	// "encoding/csv"
	// "unicode/utf8"
	"errors"
	"fmt"
	// "io"
	"log"
	"my/booksdb/src/ini"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TBook struct {
	id_fzip    int
	book_autor string
	book_name  string
	book_date  time.Time
	code       string
}

var (
	path_for_inp    string
	path_for_zip    string
	file_list_books string
)

// импорт списка файлов с описание (*.inp)
func (pg *postgres) ImpInp() {

	print("Импорт файлов описаний")

	// читаем секцию import, ключь list
	if ini.Parser("import", "list") {
		file_list_books = ini.Get()
	}

	_, err := pg.db.Exec(context.Background(), `TRUNCATE TABLE fzip`)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(file_list_books)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		println(scanner.Text())
		_, err := pg.db.Exec(context.Background(), `INSERT INTO fzip (fzip_name) VALUES ($1)`, strings.TrimSpace(scanner.Text()))
		if err != nil {
			log.Fatal(err)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

// импорт списка файлов с описание (*.inp - открываются как обычные txt и потом строки разбиваются через spit)
func (pg *postgres) ImpBooks2() {
	var (
		id   int
		fzip string
		// book_date time.Time
	)

	// выдиляет только первого автора
	t := func(s string) string {
		if i := strings.Index(s, ":"); i > 0 {
			if j := strings.Index(s[i:], ":"); j > 0 {
				// если много авторов
				return strings.ReplaceAll(s[:i], ",", " ") + ", " + strings.ReplaceAll(s[i:j], ",", " ") + ", + ... "
			}
			return strings.ReplaceAll(s[:i], ",", " ")
		}
		return strings.ReplaceAll(s, ",", " ")
	}

	println("Импорт описаний книг")

	if ini.Parser("import") {
		path_for_inp = ini.Val("inp")
		ini.DelVals()
	}

	_, err := pg.db.Exec(context.Background(), `TRUNCATE TABLE books`)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := pg.db.Query(context.Background(), "SELECT fzip_id, TRIM(fzip_name) FROM fzip")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		// цикл по списку файлов
		if err := rows.Scan(&id, &fzip); err != nil {
			log.Fatal(err)
			return
		}

		// file, err := os.Open("/media/n9/books-t/inpx/" + fzip)
		file, err := os.Open(path_for_inp + fzip)
		if err != nil {
			panic(err.Error())
		}
		defer file.Close()

		println("Импорт", fzip)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {

			l := strings.Split(scanner.Text(), fmt.Sprintf("%c", 0x4))
			// println(scanner.Text())
			// println(l[0], fmt.Sprintf("%s %s", l[2], l[3]), l[5], l[10])

			book_date, err := time.Parse("2006-01-02", l[10])
			if err != nil {
				panic(err.Error())
			}

			_, err = pg.db.Exec(context.Background(), `INSERT INTO books (id_fzip,book_autor,book_name,book_name2,book_date,code,code2,book_lang) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, id, t(l[0]), l[2], l[3], book_date, l[5], l[6], l[11])
			if err != nil {
				log.Fatal(err)
			}

		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		file.Close()
	} // цикл по файлам

}

// распаковка книги из архива
func UnZip(fzip, code string) string {
	const FOLDER_BOOK = "books"

	if len(code) <= 3 || len(fzip) <= 5 {
		log.Println("Нет данных для распаковки", fzip, code)
		return "Ошибка. нет параметров для распаковки."
	}

	if ini.Parser("import", "zip") {
		path_zip := ini.Get()
		file_zip := path_zip + strings.ReplaceAll(fzip, ".inp", ".zip")

		if _, err := os.Stat(file_zip); errors.Is(err, os.ErrNotExist) {
			log.Println("Нет файла", err)
			return "Нет файла c архивом: " + file_zip
		}

		app := "unzip"
		ar1 := file_zip
		ar2 := code + ".fb2"
		ar3 := "-d"
		ar4 := path_zip + FOLDER_BOOK

		if _, err := os.Stat(path_zip + FOLDER_BOOK + "/" + code + ".fb2"); errors.Is(err, os.ErrNotExist) {

			cmd := exec.Command(app, ar1, ar2, ar3, ar4)

			if errors.Is(cmd.Err, exec.ErrDot) {
				cmd.Err = nil
			}
			if err := cmd.Run(); err != nil {
				log.Println(app, ar1, ar2, ar3, ar4)
				log.Fatal(err)
			}

		} else {
			return "Такой файл уже существует:\n" + path_zip + FOLDER_BOOK + "/" + code + ".fb2"
		}
	}
	return ""
}
