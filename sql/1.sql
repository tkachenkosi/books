-- пример удаление таблици с индексами
-- drop table books cascade;

-- alter table fzip alter column fzip_name type VARCHAR(40);

-- alter table books alter column book_autor type VARCHAR(200);
-- alter table books alter column book_name type VARCHAR(300);

-- TRUNCATE TABLE fzip;
-- TRUNCATE TABLE books;


-- alter table books add column book_lang varchar(5);
-- alter table books add column book_rem varchar(100);


ALTER SEQUENCE fzip_fzip_id RESTART WITH 0;
ALTER SEQUENCE books_book_id RESTART WITH 0;

