package tasks

import (
	"HW4/internal/config"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrConnectingMySQL    = errors.New("error of connection mysql db")
	ErrPingMySQL          = errors.New("error of ping mysql db")
	ErrCreatingTableMySQL = errors.New("error of creating tasks table")
)

// TODO ВЫДЕЛИ В ФУНКЦИЮ, код повторяется

type TasksRepoMySQL struct {
	DB *sql.DB
}

func NewTasksRepoMySQL(ctx context.Context, config *config.MySQLDb) *TasksRepoMySQL {
	cfg := mysql.Config{
		User:              config.User,
		Passwd:            config.Password,
		Addr:              config.Host + ":" + config.Port,
		DBName:            config.Name,
		InterpolateParams: true,
	}

	connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		log.Fatalf("error: %s, Description: %s", err, ErrConnectingMySQL)
	}

	db := sql.OpenDB(connector)

	// db.SetMaxOpenConns(10)
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error: %s, Description: %s", err, ErrPingMySQL)
	}

	query := `CREATE TABLE IF NOT EXISTS tasks(tasks_id int primary key auto_increment, owner text, 
	executor text, description text, completed bool, assigned bool)`

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("Error %s, Description: %s", err, ErrCreatingTableMySQL)
	}

	return &TasksRepoMySQL{DB: db}
}

func (repo *TasksRepoMySQL) GetAllTasks(ctx context.Context) ([]*Task, error) {
	tasks := []*Task{}
	rows, err := repo.DB.QueryContext(ctx, "SELECT id, owner, executor, description, completed, assigned FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("select mysql error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		err = rows.Scan(&task.ID, &task.Owner, &task.Executor, &task.Description, &task.Completed, &task.Assigned)
		if err != nil {
			return nil, fmt.Errorf("scanning mysql error: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (repo *TasksRepoMySQL) Add(ctx context.Context, task *Task) error {
	_, err := repo.DB.ExecContext(ctx,
		"INSERT INTO tasks (`owner`, `executor`, `description`, `completed`, `assigned`) VALUES (?, ?, ?, ?, ?)",
		task.Owner,
		task.Executor,
		task.Description,
		task.Completed,
		task.Assigned,
	)
	if err != nil {
		return fmt.Errorf("insert mysql error: %w", err)
	}
	return nil
}

func (repo *TasksRepoMySQL) GetCreatedTasks(ctx context.Context, username string) ([]*Task, error) {
	tasks := []*Task{}
	rows, err := repo.DB.QueryContext(ctx, "SELECT id, owner, executor, description, completed, assigned FROM tasks WHERE owner=?", username)
	if err != nil {
		return nil, fmt.Errorf("select mysql error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		err = rows.Scan(&task.ID, &task.Owner, &task.Executor, &task.Description, &task.Completed, &task.Assigned)
		if err != nil {
			return nil, fmt.Errorf("scanning mysql error: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (repo *TasksRepoMySQL) GetMyTasks(ctx context.Context, username string) ([]*Task, error) {
	tasks := []*Task{}
	rows, err := repo.DB.QueryContext(ctx, "SELECT id, owner, executor, description, completed, assigned FROM tasks WHERE executor=?", username)
	if err != nil {
		return nil, fmt.Errorf("select mysql error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		err = rows.Scan(&task.ID, &task.Owner, &task.Executor, &task.Description, &task.Completed, &task.Assigned)
		if err != nil {
			return nil, fmt.Errorf("scanning mysql error: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (repo *TasksRepoMySQL) Assign(ctx context.Context, taskId uint64, username string) error {
	_, err := repo.DB.QueryContext(ctx, "UPDATE tasks SET `executor` = ?, `assigned` = 1, WHERE id = ?", username, taskId)
	if err != nil {
		return fmt.Errorf("update mysql error: %w", err)
	}
	return nil
}

func (repo *TasksRepoMySQL) Unassign(ctx context.Context, taskId uint64) error {
	_, err := repo.DB.QueryContext(ctx, "UPDATE tasks SET `assigned` = 0, WHERE id = ?", taskId)
	if err != nil {
		return fmt.Errorf("update mysql error: %w", err)
	}
	return nil
}

func (repo *TasksRepoMySQL) Complete(ctx context.Context, taskId uint64) error {
	_, err := repo.DB.QueryContext(ctx, "UPDATE tasks SET `completed` = 1, WHERE id = ?", taskId)
	if err != nil {
		return fmt.Errorf("update mysql error: %w", err)
	}
	return nil
}
