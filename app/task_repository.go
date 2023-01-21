package ddosy

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(dbURL string) *TaskRepository {
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalln(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS Tasks (
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		StatusId INTEGER NOT NULL,
		CreatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		StartedAt DATETIME NULL,
		KilledAt DATETIME NULL,
		DoneAt DATETIME NULL,
		Request BLOB NOT NULL,
		Results TEXT DEFAULT ""
	);`

	log.Println("create tasks table")
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalln(err)
	}

	// TODO index on status, id and doneAt

	rep := &TaskRepository{
		db: db,
	}

	// todo backgroud worker to clean database

	// create database
	return rep
}

func (r *TaskRepository) Save(req ScheduleRequestWeb) (uint64, error) {
	query := `
		INSERT INTO Tasks (StatusId, Request)
		VALUES (?, ?)
		RETURNING Id;
    `

	bs, err := json.Marshal(req)
	if err != nil {
		log.Printf("error marshaling req %v %s\n", req, err)
		return 0, err
	}

	var id uint64
	if err = r.db.QueryRow(query, Scheduled, string(bs)).Scan(&id); err != nil {
		log.Printf("error saving to db %s\n", err)
		return 0, err
	}

	return id, nil
}

func (r *TaskRepository) Get(id uint64) (DatabaseTask, error) {
	query := `
	SELECT 
		Id, StatusId,
		CreatedAt, StartedAt, KilledAt, DoneAt, 
		Request, Results 
	FROM Tasks 
	WHERE Id = ?;`

	var t DatabaseTask
	var buf []byte
	err := r.db.QueryRow(query, id).Scan(
		&t.Id,
		&t.StatusId,
		&t.CreatedAt,
		&t.StartedAt,
		&t.KilledAt,
		&t.DoneAt,
		&buf,
		&t.Results,
	)
	if err != nil {
		log.Printf("error getting a task with id=%d %s\n", id, err)
		return t, err
	}

	if err = json.Unmarshal(buf, &t.Request); err != nil {
		log.Printf("cannot unmarshal req form db for id=%d\n", id)
		return t, err
	}

	return t, nil
}

func (r *TaskRepository) UpdateStatus(id uint64, newStatus TaskStatus) error {
	log.Printf("update status of id=%d to newStatusId=%d\n", id, newStatus)

	var query string
	switch newStatus {
	case Running:
		// update only if current status is Scheduled(1)
		query = `UPDATE Tasks
		SET StatusId = (CASE WHEN StatusId = 1 THEN 2 ELSE StatusId END),
		StartedAt = (CASE WHEN StatusId = 1 THEN CURRENT_TIMESTAMP ELSE StartedAt END)
		WHERE Id = ?;`
	case Killed:
		// update only if current status is Scheduled(1) OR Running(2)
		query = `UPDATE Tasks
		SET StatusId = (CASE WHEN StatusId IN (1, 2) THEN 3 ELSE StatusId END),
		KilledAt = (CASE WHEN StatusId IN (1, 2) THEN CURRENT_TIMESTAMP ELSE KilledAt END)
		WHERE Id = ?;`
	case Done:
		// update only if current status is Running(2)
		query = `UPDATE Tasks
		SET StatusId = (CASE WHEN StatusId = 2 THEN 4 ELSE StatusId END),
		DoneAt = (CASE WHEN StatusId = 2 THEN CURRENT_TIMESTAMP ELSE DoneAt END)
		WHERE Id = ?;`
	default:
		return nil
	}

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("error on update id=%d %s\n", id, err)
		return err
	}
	return nil
}

func (r *TaskRepository) UpdateProgress(id uint64, progress string) error {
	query := "UPDATE Tasks SET Results = Tasks.Results || '\n' || ? WHERE Id = ?;"
	_, err := r.db.Exec(query, progress, id)
	if err != nil {
		log.Printf("error on update id=%d %s\n", id, err)
		return err
	}
	return nil
}

// func (r *Repository) deleteWorker() {

// }

// func (r *Repository) InsertNew(task LoadTask) (uint64, error) {
// 	req, err := json.Marshal(task.req)
// 	if err != nil {
// 		log.Printf("error converting task to json %s\n", err)
// 		return 0, err
// 	}

// 	stmt, err := r.db.Prepare(
// 		`INSERT INTO TaskStatus (Request)
// 		VALUES (?)
// 		RETURNING Id;`)

// 	if err != nil {
// 		log.Printf("error creating prepared statment %s\n", err)
// 		return 0, err
// 	}

// 	var out uint64
// 	err = stmt.QueryRow(req).Scan(&out)
// 	if err != nil {
// 		log.Fatalf("error %s\n", err)
// 	}

// 	return out, nil
// }

// func (r *Repository) Close() {
// 	r.db.Close()
// }
