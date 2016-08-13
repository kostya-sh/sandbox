package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func savepoints(db *sql.DB) {
	tx, err := db.Begin()
	ck(err)

	_, err = tx.Exec("insert into t1 values ('go1')")
	ck(err)
	_, err = tx.Exec("savepoint s1")
	ck(err)
	_, err = tx.Exec("insert into t1 values ('go2')")
	ck(err)
	_, err = tx.Exec("rollback to savepoint s1")
	ck(err)

	tx.Commit()
	ck(err)
}

func preparedStatementsDateTrunc(db *sql.DB) {
	_, err := db.Prepare("select date_trunc('hour', $1)")
	//	_, err := db.Prepare("select date_trunc('hour', cast($1 as timestamp))")
	ck(err)
}

func queryDateTrunc(db *sql.DB) {
	_, err := db.Query("select date_trunc('hour', $1)", time.Now())
	ck(err) // 2016/08/13 23:41:13 pq: function date_trunc(unknown, unknown) is not unique
}

func main() {
	db, err := sql.Open("postgres", "postgresql://ksh:ksh@localhost/mytest")
	ck(err)

	// savepoints(db)
	// preparedStatementsDateTrunc(db)
	queryDateTrunc(db)
}
