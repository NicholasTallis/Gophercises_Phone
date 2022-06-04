package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

//Updated!!!
// Phone represents the phone_numbers table in the DB
type Phone struct {
	Number string
	ID     int
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

type DB struct {
	db *sql.DB
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Seed2(data, ids []string) error {
	if len(data) == len(ids) {
		for i := range data {
			if _, err := insertPhone2(db.db, data[i], ids[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func insertPhone2(db *sql.DB, phone string, groupid string) (int, error) {
	statement := `INSERT INTO public.numbers(num,groupid) VALUES($1,$2) RETURNING groupid`
	var id int
	err := db.QueryRow(statement, phone, groupid).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (db *DB) AllPhones() ([]Phone, error) {
	rows, err := db.db.Query("SELECT num, groupid FROM numbers ORDER BY groupid")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.Number, &p.ID); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	row := db.db.QueryRow("SELECT * FROM numbers WHERE num=$1", number)
	err := row.Scan(&p.Number, &p.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func (db *DB) UpdatePhone(p *Phone, newNum string) error {
	statement := `UPDATE numbers SET num=$2 WHERE num=$1`
	_, err := db.db.Exec(statement, p.Number, newNum)
	return err
}

func (db *DB) DeletePhone(p *Phone) error {
	statement := `DELETE FROM numbers WHERE num=$1`
	_, err := db.db.Exec(statement, p.Number)
	return err
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumbersTable(db)
	if err != nil {
		return err
	}
	return db.Close()
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
    CREATE TABLE IF NOT EXISTS numbers (
      num varchar(40),
      groupid integer
    )`
	_, err := db.Exec(statement)
	return err
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

// We don't use this right now.
func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE id=$1", id)
	err := row.Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
}
func Foo() {
	println("Test for db package")
}
func Foo3() {
	println("Test3 for db package")
	println("Test4!")
}
