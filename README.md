# books
Программа для работы с архивом книг из сайта флибуста.
База конвертирована в postgres. Web server написан на golang. 

Памятка по утилитам:
psql -h 192.168.1.22 -U user1 -d booksdb
psql -h 192.168.1.22 -U user1 -d booksdb -f 3.sql
rg -i -tgo bcrypt = рекрусивный поиск только в файлах ".go" фразу "bcrypt"
GOOS=linux GOARCH=amd64 go build main.go -o ./hello_world
