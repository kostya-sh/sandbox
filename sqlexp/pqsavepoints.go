package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := sql.Open("postgres", "postgresql://ksh:ksh@localhost/mytest")
	ck(err)

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
