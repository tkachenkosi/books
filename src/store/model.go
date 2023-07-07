// служебные функции для работы с базой postgres
package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"my/booksdb/src/ini"
	"sync"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	Db     *postgres // ссылка на глобальный экземпляр для всей программы
	pgOnce sync.Once
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
	// ErrUnZipBook    = errors.New("unzip book")
)

// формирует строку подключения к базе
func pgConf() string {
	if ini.Parser("postgres") {
		return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s pool_max_conns=%s", ini.Val("user"), ini.Val("passwd"), ini.Val("host"), ini.Val("port"), ini.Val("db"), ini.Val("conns"))
	} else {
		panic("Нет настроек для подключение к postgres")
	}
}

// подключение
func NewPG(ctx context.Context) error {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, pgConf())
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}

		Db = &postgres{db}
	})

	return nil
}

// тестирование соединения
func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}
