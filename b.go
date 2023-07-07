// изучение
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	// "strings"
	// "unicode/utf8"
	// "encoding/json"
	// "html/template"
	// "log"
	// "net/http"
	"my/booksdb/src/store"
	// "my/booksdb/src/tools"
	fiber "github.com/gofiber/fiber/v2"
	mustache "github.com/gofiber/template/mustache/v2"
	"my/booksdb/src/ini"
)

type Tfound struct {
	Autor string `json:"autor"`
	Name  string `json:"name"`
}

type TOneFound struct {
	Code string `json:"code"`
}

var (
	err         error
	count_query int
)

func main() {
	fmt.Println("Поиск книг")

	// настройка логирования
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	log.SetOutput(file)
	// log.SetOutput(os.Stdout)

	ctx := context.Background()

	err = store.NewPG(ctx)
	if err != nil {
		panic(err.Error())
	}

	err = store.Db.Ping(ctx)
	if err != nil {
		panic(err.Error())
	}

	// импорт файлов описаний
	// store.Db.ImpInp()
	// импорт книг
	// store.Db.ImpBooks2()

	// host := conf.Sec("web").GetStr("host")
	// port := conf.Sec("web").GetStr("port")

	engine := mustache.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views:                 engine,
		DisableStartupMessage: true,
	})

	app.Static("/", "./public")
	app.Get("/", app_main)
	app.Post("/found", app_found)
	app.Post("/result", app_result)
	app.Post("/count", app_count)
	app.Post("/book", app_book)

	if ini.Parser("web") {
		app.Listen(fmt.Sprintf("%s:%s", ini.Val("host"), ini.Val("port")))
	} else {
		app.Listen(":3000")
	}

	store.Db.Close()
	log.Println("End app")
}

func app_main(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{"title": "Тестовый сервер. (golang,fiber,postgres)"})
}

// количестов строк в запросе для ответом шаблоном
func app_count(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("%d", count_query))
}

// поиск одной книги и распаковка из архива
func app_book(c *fiber.Ctx) error {
	var (
		oneFound TOneFound
		bookOne  store.TBookOne
	)

	if err := c.BodyParser(&oneFound); err != nil {
		log.Println("Err =>", err.Error())
	}

	bookOne = store.Db.GetOneBook(strings.TrimSpace(oneFound.Code))
	ok := store.UnZip(bookOne.Fzip_name, oneFound.Code) // распаковка книги
	if len(ok) > 1 {
		bookOne.Book_name = ok     // текст ошибки от unzip
		bookOne.Book_autor = "err" // признак ошибки
	}
	bookOne.Fzip_name = strings.ReplaceAll(bookOne.Fzip_name, ".inp", ".zip")
	return c.JSON(bookOne)
}

// поиск книг и возврат json
func app_found(c *fiber.Ctx) error {
	var (
		found Tfound
		books []store.TBooks
	)

	if err := c.BodyParser(&found); err != nil {
		log.Println("Err =>", err.Error())
	}

	found.Autor = strings.ToLower(strings.TrimSpace(found.Autor))
	found.Name = strings.ToLower(strings.TrimSpace(found.Name))
	// found.Autor = strings.TrimSpace(found.Autor)
	// found.Name = strings.TrimSpace(found.Name)
	// println(found.Autor, found.Name, len(found.Autor), len(found.Name))

	if len(found.Autor) > 0 && len(found.Name) == 0 {
		books, err = store.Db.GetBooksAutor(found.Autor)
	} else if len(found.Autor) == 0 && len(found.Name) > 0 {
		books, err = store.Db.GetBooksName(found.Name)
	} else if len(found.Autor) > 0 && len(found.Name) > 0 {
		books, err = store.Db.GetBooks(found.Autor, found.Name)
	}

	if err != nil {
		log.Println("Err =>", err.Error())
	}
	// println("len", len(books))
	return c.JSON(books)
}

// поиск книг и возврат шаблон
func app_result(c *fiber.Ctx) error {
	var (
		found Tfound
		books []store.TBooks
	)

	if err := c.BodyParser(&found); err != nil {
		log.Println("Err =>", err.Error())
	}
	found.Autor = strings.ToLower(strings.TrimSpace(found.Autor))
	found.Name = strings.ToLower(strings.TrimSpace(found.Name))
	// found.Autor = strings.TrimSpace(found.Autor)
	// found.Name = strings.TrimSpace(found.Name)

	if len(found.Autor) > 0 && len(found.Name) == 0 {
		books, err = store.Db.GetBooksAutor(found.Autor)
	} else if len(found.Autor) == 0 && len(found.Name) > 0 {
		books, err = store.Db.GetBooksName(found.Name)
	} else if len(found.Autor) > 0 && len(found.Name) > 0 {
		books, err = store.Db.GetBooks(found.Autor, found.Name)
	}

	if err != nil {
		log.Println("Err =>", err.Error())
	}

	count_query = len(books)

	// попытка конвертации struct -> map
	mb, _ := json.Marshal(&books)
	// var m []map[string]interface{}
	var m []map[string]any
	_ = json.Unmarshal(mb, &m)
	// println(m)

	return c.Render("result", fiber.Map{"line": m})
}
