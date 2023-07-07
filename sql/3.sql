-- инициализация базы для поиска книг
-- DROP TABLE IF EXISTS users;
CREATE TABLE users (user_id SERIAL PRIMARY KEY, user_login VARCHAR(40) NOT NULL, user_passwd VARCHAR(40) NOT NULL, user_name VARCHAR(50), user_email VARCHAR(100) NOT NULL, user_balance INTEGER NOT NULL);
CREATE INDEX users_name_idx ON users (user_name);
-- DROP TABLE IF EXISTS fzip;
CREATE TABLE fzip (fzip_id SERIAL PRIMARY KEY, fzip_name VARCHAR(40) NOT NULL);
-- DROP TABLE IF EXISTS books;
CREATE TABLE books (book_id SERIAL PRIMARY KEY, id_fzip INTEGER NOT NULL, book_autor VARCHAR(200) NOT NULL, book_name VARCHAR(200) NOT NULL, book_name2 VARCHAR(300), book_date DATE NOT NULL, code VARCHAR(10) NOT NULL, code2 VARCHAR(10), book_lang VARCHAR(5), book_rem VARCHAR(100));
CREATE INDEX books_fzip_idx ON books (id_fzip);
CREATE INDEX books_autor_idx ON books (book_autor);
CREATE INDEX books_name_idx ON books (book_name);
CREATE INDEX books_code_idx ON books (code);
