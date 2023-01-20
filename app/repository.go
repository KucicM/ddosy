package ddosy

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dbURL string) *Repository {
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalln(err)
	}

	query := `CREATE TABLE TaskStatus (
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		RequestAcceptedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		Request BLOB
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}

	// create database
	return &Repository{
		db: db,
	}
}

func (r *Repository) InsertNew(task LoadTask) (uint64, error) {
	req, err := json.Marshal(task.req)
	if err != nil {
		log.Printf("error converting task to json %s\n", err)
		return 0, err
	}

	stmt, err := r.db.Prepare(
		`INSERT INTO TaskStatus (Request) 
		VALUES (?) 
		RETURNING Id;`)

	if err != nil {
		log.Printf("error creating prepared statment %s\n", err)
		return 0, err
	}

	var out uint64
	err = stmt.QueryRow(req).Scan(&out)
	if err != nil {
		log.Fatalf("error %s\n", err)
	}

	return out, nil
}

func (r *Repository) UpdateStatus(id uint64, status string) {

}

func (r *Repository) Close() {
	r.db.Close()
}
