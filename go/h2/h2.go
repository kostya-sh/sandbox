package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	// &disable_prepared_binary_result=yes
	connURL = "postgres://sa:sa@localhost:5435/test?sslmode=disable&binary_parameters=no&disable_prepared_binary_result=yes"
)

var (
	db *sql.DB
)

func exec(sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatalf("Failed to execute %q: %v", sql, err)
	}
}

func initDB() {
	var err error
	db, err = sql.Open("postgres", connURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	exec("DROP TABLE IF EXISTS Fortune")

	exec(`CREATE TABLE Fortune (
           id integer NOT NULL,
           message varchar(2048) NOT NULL,
           PRIMARY KEY  (id))`)

	exec(`INSERT INTO Fortune (id, message) VALUES (1, 'fortune: No such file or directory')`)
	exec(`INSERT INTO Fortune (id, message) VALUES (2, 'A computer scientist is someone who fixes things that aren''t broken.')`)
	exec(`INSERT INTO Fortune (id, message) VALUES (3, 'After enough decimal places, nobody gives a damn.')`)
	exec(`INSERT INTO Fortune (id, message) VALUES (4, 'A bad random number generator: 1, 1, 1, 1, 1, 4.33e+67, 1, 1, 1')`)
	exec(`INSERT INTO Fortune (id, message) VALUES (5, 'A computer program does what you tell it to do, not what you want it to do.')`)
	exec(`INSERT INTO Fortune (id, message) VALUES (6, 'Emacs is a nice operating system, but I prefer UNIX. — Tom Christaensen')`)
	exec(`INSERT INTO Fortune (id, message) VALUES (7, 'Any program that runs right is obsolete.')`)
	exec(`INSERT INTO Fortune (id, message) VALUES (8, 'A list is only as strong as its weakest link. — Donald Knuth')`)
	exec(`INSERT INTO Fortune (id, message) VALUES (9, 'Feature: A bug with seniority.')`)
	exec(`INSERT INTO Fortune (id, message) VALUES (10, 'Computers make very fast, very accurate mistakes.')`)
	exec(`INSERT INTO Fortune (id, message) VALUES (11, '<script>alert("This should not be displayed in a browser alert box.");</script>')`)
	exec(`INSERT INTO Fortune (id, message) VALUES (12, 'フレームワークのベンチマーク')`)

	// worldSelectPrepared1, err = db1.Prepare(worldSelect1)
	// if err != nil {
	// 	log.Fatalf("Failed to prepare %q: %s", worldSelect1, err)
	// }
	// worldUpdatePrepared1, err = db1.Prepare(worldUpdate1)
	// if err != nil {
	// 	log.Fatalf("Failed to prepare %q: %s", worldUpdate1, err)
	// }
	// fortuneSelectPrepared1, err = db.Prepare(fortuneSelect1)
	// if err != nil {
	// 	log.Fatalf("Failed to prepare %q: %s", fortuneSelect1, err)
	// }
	//
	// fortuneSelectPrepared2, err = db.Prepare(fortuneSelect2)
	// if err != nil {
	// 	log.Fatalf("Failed to prepare %q: %s", fortuneSelect1, err)
	// }
}

func main() {
	initDB()

	fetchAllFortunes()
	fmt.Println()

	fetchOneFortune()
	fmt.Println()
}

func fetchAllFortunes() {
	stmt, err := db.Prepare("SELECT id, message FROM Fortune")
	if err != nil {
		log.Fatalf("Failed to prepare (fetchAllFortunes): %v", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatalf("Failed to query (fetchAllFortunes): %v", err)
	}
	defer rows.Close()

	for rows.Next() { //Fetch rows
		var id int
		var m string
		if err = rows.Scan(&id, &m); err != nil {
			log.Fatalf("Failed to scan (fetchAllFortunes)): %s", err)
		}
		fmt.Printf("%d\t%s\n", id, m)
	}

	if rows.Err() != nil {
		log.Fatalf("Error (fetchAllFortunes)): %s", rows.Err())
	}
}

func fetchOneFortune() {
	stmt, err := db.Prepare("SELECT id, message FROM Fortune where id = $1")
	if err != nil {
		log.Fatalf("Failed to prepare (fetchAllFortunes): %v", err)
	}

	rows, err := stmt.Query(12)
	if err != nil {
		log.Fatalf("Failed to query (fetchAllFortunes): %v", err)
	}
	defer rows.Close()

	for rows.Next() { //Fetch rows
		var id int
		var m string
		if err = rows.Scan(&id, &m); err != nil {
			log.Fatalf("Failed to scan (fetchAllFortunes)): %s", err)
		}
		fmt.Printf("%d\t%s\n", id, m)
	}

	if rows.Err() != nil {
		log.Fatalf("Error (fetchAllFortunes)): %s", rows.Err())
	}
}
