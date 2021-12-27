package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func main() {
	initdb()
	exampleSelectUpdate()
	//exampleUpdateSelect()
	// exampleWithoutLock()
	// exampleLock()
}

func exampleSelectUpdate() {
	tx, errTx := db.Begin()
	if errTx != nil {
		log.Print("err tx ", errTx)
	}

	tx2, errTx2 := db.Begin()
	if errTx2 != nil {
		log.Print("err tx ", errTx2)
	}

	// select without tx
	qry := `
		SELECT
			title
		FROM testes;`

	rows, errQuery := tx.Query(qry)
	if errQuery != nil {
		log.Print("err query ", errQuery)
	}

	defer rows.Close()
	for rows.Next() {
		column := ""

		errScan := rows.Scan(&column)
		if errScan != nil {
			log.Print("err query ", column)
		}

		log.Print("result column ", column)
	}

	// update with tx
	qryu := `UPDATE testes
		SET title='2';`

	_, errUpdate := tx2.Exec(qryu)
	if errUpdate != nil {
		log.Print("err exec ", errUpdate)
	}

	tx.Rollback()
}

func exampleUpdateSelect() {
	tx, errTx := db.Begin()
	if errTx != nil {
		log.Print("err tx ", errTx)
	}

	// update with tx
	qryu := `UPDATE testes
	SET title='1';`

	_, errUpdate := tx.Exec(qryu)
	if errUpdate != nil {
		log.Print("err exec ", errUpdate)
	}

	// select without tx
	qry := `
	SELECT
		title
	FROM testes;`

	rows, errQuery := tx.Query(qry)
	if errQuery != nil {
		log.Print("err query ", errQuery)
	}

	defer rows.Close()
	for rows.Next() {
		column := ""

		errScan := rows.Scan(&column)
		if errScan != nil {
			log.Print("err query ", column)
		}

		log.Print("result column ", column)
	}

	tx.Rollback()
}

func exampleWithoutLock() {
	tx, errTx := db.Begin()
	if errTx != nil {
		log.Print("err tx ", errTx)
	}

	// update with tx
	qryu := `UPDATE testes
	SET title='1';`

	_, errUpdate := tx.Exec(qryu)
	if errUpdate != nil {
		log.Print("err exec ", errUpdate)
	}

	// select without tx
	qry := `
	SELECT
		title
	FROM testes;`

	rows, errQuery := db.Query(qry)
	if errQuery != nil {
		log.Print("err query ", errQuery)
	}

	defer rows.Close()
	for rows.Next() {
		column := ""

		errScan := rows.Scan(&column)
		if errScan != nil {
			log.Print("err query ", column)
		}

		log.Print("result column ", column)
	}

	tx.Rollback()
}

func exampleLock() {
	tx, errTx := db.Begin()
	if errTx != nil {
		log.Print("err tx ", errTx)
	}

	// update with tx
	qryu := `UPDATE testes
	SET title='1';`

	_, errUpdate := tx.Exec(qryu)
	if errUpdate != nil {
		log.Print("err exec ", errUpdate)
	}

	// insert without tx - lock
	qryi := `INSERT INTO testes
	(title) 
	VALUES ('TESTE');`

	_, errInsert := db.Exec(qryi)
	if errInsert != nil {
		log.Print("err exec ", errInsert)
	}

	tx.Rollback()
}

func initdb() {
	dba, err := sql.Open("mysql", "root:r00t@/debt")
	if err != nil {
		log.Panic("error open connection database ", err)
	}

	dba.SetMaxOpenConns(5)
	dba.SetMaxIdleConns(2)
	dba.SetConnMaxLifetime(time.Minute * 5)

	if err = dba.Ping(); err != nil {
		log.Panic("error ping connection database ", err)
	}

	db = dba
}
