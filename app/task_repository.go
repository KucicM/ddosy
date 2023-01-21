package ddosy

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type TaskRepository struct {
	url string
}

func NewTaskRepository(dbURL string, truncate bool) *TaskRepository {
	log.Printf("creating db: %s\n", dbURL)
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if truncate {
		log.Println("cleaning table")
		_, err = db.Exec("DROP TABLE IF EXISTS Tasks")
		if err != nil {
			log.Fatalln(err)
		}
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
		url: dbURL,
		// db: db,
	}

	// todo backgroud worker to clean database

	return rep
}

func (r *TaskRepository) Save(req ScheduleRequestWeb) (uint64, error) {
	query := `
		INSERT INTO Tasks (StatusId, Request)
		VALUES (?, ?)
		RETURNING Id;
    `

	db, err := sql.Open("sqlite3", r.url)
	if err != nil {
		return 0, err
	}

	bs, err := json.Marshal(req)
	if err != nil {
		log.Printf("error marshaling req %v %s\n", req, err)
		return 0, err
	}

	var id uint64
	if err = db.QueryRow(query, Scheduled, string(bs)).Scan(&id); err != nil {
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

	db, err := sql.Open("sqlite3", r.url)
	if err != nil {
		return t, err
	}

	var buf []byte
	err = db.QueryRow(query, id).Scan(
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

	db, err := sql.Open("sqlite3", r.url)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, id)
	if err != nil {
		log.Printf("error on update id=%d %s\n", id, err)
		return err
	}
	return nil
}

func (r *TaskRepository) UpdateProgress(id uint64, progress string) error {
	query := "UPDATE Tasks SET Results = Tasks.Results || '\n' || ? WHERE Id = ?;"

	db, err := sql.Open("sqlite3", r.url)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, progress, id)
	if err != nil {
		log.Printf("error on update id=%d %s\n", id, err)
		return err
	}
	return nil
}

// returns task id, task original request and error
// if no tasks are found id will be 0 and error = nil
// in case of errors, id = 0 and error != nil
// this will NOT update the status, why? IDK sql
func (r *TaskRepository) GetNext() (uint64, ScheduleRequestWeb, error) {
	query := `
	SELECT Id, Request 
	FROM Tasks
	WHERE StatusId = 1 -- scheduled
	ORDER BY Id
	LIMIT 1`

	var req ScheduleRequestWeb

	db, err := sql.Open("sqlite3", r.url)
	if err != nil {
		return 0, req, err
	}

	var id uint64
	var buf []byte
	err = db.QueryRow(query).Scan(&id, &buf)

	if err == sql.ErrNoRows {
		return 0, req, nil
	}

	if err != nil {
		log.Printf("error fetching from db %s\n", err)
		return 0, req, err
	}

	if err := json.Unmarshal(buf, &req); err != nil {
		return 0, req, err
	}

	return id, req, nil
}
