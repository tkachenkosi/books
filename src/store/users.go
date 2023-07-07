package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"log"
)

// для автоматического сохраниение в струтуру результатов select
type TUser struct {
	User_id      int
	User_login   string
	User_name    string
	User_passwd  string
	User_email   string
	User_balance int
}

type SUser struct {
	Id      int
	Login   string
	Name    string
	Passwd  string
	Email   string
	Balance int
}

// добавить одну запись
func (pg *postgres) InsertUser(ctx context.Context, args pgx.NamedArgs) error {
	query := `INSERT INTO users (user_login, user_passwd, user_name, user_email, user_balance) VALUES (@login, @passwd, @name, @email, @balance)`

	// args := pgx.NamedArgs{arg}

	/* args := pgx.NamedArgs{
		"login":   "efim16",
		"passwd":  "123",
		"name":    "Слукин Ефим Макарович",
		"email":   "efim16@mail.ru",
		"balance": "0",
	} */

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

// пример добовление записи и возврата нового id
func (pg *postgres) AddUser(args pgx.NamedArgs) uint64 {
	var id uint64
	if err := pg.db.QueryRow(context.Background(), `INSERT INTO users (user_login, user_passwd, user_name, user_email, user_balance) VALUES (@login, @passwd, @name, @email, @balance) RETURNING user_id`, args).Scan(&id); err != nil {
		log.Println("Ошибка to INSERT:", err.Error())
	}
	return id
}

// изменить одну запись passwd md5
func (pg *postgres) UpdateUser(args pgx.NamedArgs) error {
	// query := `UPDATE users SET user_passwd = @passed WHERE user_id = @id`

	_, err := pg.db.Exec(context.Background(), `UPDATE users SET user_passwd = @passwd WHERE user_id = @id`, args)
	if err != nil {
		return err
	}

	return nil
}

// криптографическая хранения пароля
func (pg *postgres) SetUserPasswd(args pgx.NamedArgs) error {
	_, err := pg.db.Exec(context.Background(), `UPDATE users SET user_passwd = @passwd WHERE user_id = @id`, args)
	return err
}

// поиск одной записи по логину
func (pg *postgres) GetOneUser(login string) (*SUser, error) {
	// row := pg.db.QueryRow(context.Background(), "SELECT user_id,user_login,user_passwd,user_name,user_email,user_balance FROM users WHERE user_login LIKE $1 || '%'", login)
	row := pg.db.QueryRow(context.Background(), "SELECT user_id,user_login,user_passwd,user_name,user_email,user_balance FROM users WHERE LOWER(user_login) LIKE CONCAT($1::text,'%')", login)

	var user SUser
	if err := row.Scan(&user.Id, &user.Login, &user.Passwd, &user.Name, &user.Email, &user.Balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("1:", err.Error())
			return nil, ErrNotExist
		}
		log.Println("2:", err.Error())
		return nil, err
	}
	return &user, nil
}

// добавить группу записей
func (pg *postgres) AddUsers(ctx context.Context) error {

	args := [4]pgx.NamedArgs{
		{"login": "anton54", "passwd": "123", "name": "Старков Антон Егорович", "email": "anton54@mail.ru", "balance": "1"},
		{"login": "maryamna97", "passwd": "123", "name": "Бурова Марьямна Власовна", "email": "maryamna97@rambler.ru", "balance": "1"},
		{"login": "maksim1987", "passwd": "123", "name": "Зууфин Максим Давидович", "email": "maksim1987@hotmail.com", "balance": "1"},
		{"login": "aleksandra5305", "passwd": "123", "name": "Белоногова Александра Тимофеевна", "email": "aleksandra5305@yandex.ru", "balance": "1"}}

	query := `INSERT INTO users (user_login, user_passwd, user_name, user_email, user_balance) VALUES (@login, @passwd, @name, @email, @balance)`

	for _, v := range args {
		_, err := pg.db.Exec(ctx, query, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// получение всех записей
func (pg *postgres) GetUsersAll() ([]SUser, error) {
	rows, err := pg.db.Query(context.Background(), "SELECT user_id,user_login,user_passwd,user_name,user_email,user_balance FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []SUser
	for rows.Next() {
		var user SUser
		if err := rows.Scan(&user.Id, &user.Login, &user.Passwd, &user.Name, &user.Email, &user.Balance); err != nil {
			return nil, err
		}
		all = append(all, user)
	}
	return all, nil
}

// второй вариант получение всех записей по условию. Возврат в структуру.
func (pg *postgres) GetUsers(login string) ([]TUser, error) {
	rows, err := pg.db.Query(context.Background(), "SELECT user_id,user_login,user_passwd,user_name,user_email,user_balance FROM users WHERE LOWER(user_login) LIKE CONCAT($1::text,'%')", login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// RowToStructByName[TUser] - сопостовляет поля по именам
	// RowToStructByPos[TUser] - сопостовляет поля по позиции (думаю это быстрее)
	// CollectOneRow - возвращает первую строку, работает совмесно QueryRow
	// CollectRows - выполняет итерацию по строкам и выполняет функция заполнения структуры
	// RpwTo - возвращает массив значение одного поля в переменную указанного типа (int, string)
	// RowToMap - возвращает одну запись в виде мар
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[TUser])
	if err != nil {
		log.Println(err.Error())
	}
	return users, nil
}

/*
// подключение старый вариант
func NewPG(ctx context.Context) (*postgres, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, PGConf())
		if err != nil {
			panic(err.Error())
			// return fmt.Errorf("unable to create connection pool: %w", err)
			// return fmt.Errorf("unable to create connection pool:", err)
		}

		pgInstance = &postgres{db}
	})

	return pgInstance, nil
}


func RowToMap(row CollectableRow) (map[string]any, error)

var a, b int
_, err = pgx.ForEachRow(rows, []any{&a, &b}, func() error {


func (pg *postgres) GetOneUser(login string) (*SUser, error) {
	// row := r.db.QueryRow(ctx, "SELECT * FROM websites WHERE name = $1", name)
	row := pg.db.QueryRow(context.Background(), "SELECT user_id,user_login,user_name,user_passwd,user_email,user_balance FROM users WHERE user_login = $1", login)

	var user SUser
	if err := row.Scan(&user.Id, &user.Login, &user.Name, &user.Passwd, &user.Email, &user.Balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExist
		}
		return nil, err
	}
	return &user, nil
}
*/

/*
type mapRowScanner map[string]any

func (rs *mapRowScanner) ScanRow(rows pgx.Rows) error {
	values, err := rows.Values()
	if err != nil {
		return err
	}

	*rs = make(mapRowScanner, len(values))

	for i := range values {
		(*rs)[string(rows.FieldDescriptions()[i].Name)] = values[i]
	}

	return nil
}
*/

// err = db.InsertUser(ctx)
// err = db.InsertUser(ctx, &pgx.NamedArgs{"login": "apollinariya2532", "passwd": "123", "name": "Головнина Аполлинария Егоровна", "email": "apollinariya2532@outlook.com", "balance": "0"})
// err = db.InsertUser(ctx, pgx.NamedArgs{"login": "filipp56", "passwd": "123", "name": "Кантонистов Филипп Павлович", "email": "filipp56@yandex.ru", "balance": "0"})
// err = db.InsertUser(ctx, pgx.NamedArgs{"login": "filipp56", "passwd": "123", "name": "Кантонистов Филипп Павлович", "email": "filipp56@yandex.ru", "balance": "0"})

/* err = db.InsertUser(ctx, pgx.NamedArgs{"login": "sa", "passwd": MD5("13456"), "name": "Курганский Игорь Ефимович", "email": "kurg@yandex.ru", "balance": "2"})
if err != nil {
	panic(err.Error())
} */

// изменить пароль
// err = db.UpdateUser(pgx.NamedArgs{"passwd": MD5("123"), "id": 5})
// err = db.UpdateUser(pgx.NamedArgs{"passwd": PW5("apollinariya2532", "123"), "id": 2})
/* err = db.UpdateUser(pgx.NamedArgs{"passwd": PW5("maksim1987", "abc"), "id": 6})
if err != nil {
	panic(err.Error())
} */

/* err = db.UpdateUser(pgx.NamedArgs{"passwd": PasswdEncrypt("123"), "id": 1})
if err != nil {
	panic(err.Error())
}
err = db.UpdateUser(pgx.NamedArgs{"passwd": PasswdEncrypt("123"), "id": 2})
if err != nil {
	panic(err.Error())
}
err = db.UpdateUser(pgx.NamedArgs{"passwd": PasswdEncrypt("123"), "id": 3})
if err != nil {
	panic(err.Error())
}
err = db.UpdateUser(pgx.NamedArgs{"passwd": PasswdEncrypt("123"), "id": 4})
if err != nil {
	panic(err.Error())
}
err = db.UpdateUser(pgx.NamedArgs{"passwd": PasswdEncrypt("abc"), "id": 5})
if err != nil {
	panic(err.Error())
}
err = db.UpdateUser(pgx.NamedArgs{"passwd": PasswdEncrypt("abc"), "id": 7})
if err != nil {
	panic(err.Error())
}
err = db.UpdateUser(pgx.NamedArgs{"passwd": PasswdEncrypt("123"), "id": 8})
if err != nil {
	panic(err.Error())
} */

// _ = db.AddUsers(ctx)i

// id := db.AddUser(pgx.NamedArgs{"login": "admin", "passwd": PasswdEncrypt("123"), "name": "Шевелев Стас Валерьевич", "email": "shevel@yandex.ru", "balance": "2"})
// println("Добавлен пользователь. Новый ID =", id)

// println(db.PGConf())

/*


	u, err := store.Db.GetUsersAll()
	if err != nil {
		panic(err.Error())
	}
	println("len =", len(u))
	println(u[1].Login, u[1].Name, u[1].Email, u[1].Id)
	println("***")
	for i, v := range u {
		println(i, v.Id, v.Login, v.Name, v.Email)
	}
	println("*** версия 5 ***")

	// не четкий поиск
	r, err := store.Db.GetUsers("m")
	if err != nil {
		panic(err.Error())
	}
	for _, v := range r {
		println(v.User_id, v.User_login, v.User_name, v.User_email)
	}

	println("***")

	p := tools.PasswdEncrypt("abc")
	println(p, len(p))

	println("***")

	o, err := store.Db.GetOneUser("ap")
	if err != nil {
		panic(err.Error())
	}
	println("UserOne =", o.Login, o.Name, o.Passwd)



*/
