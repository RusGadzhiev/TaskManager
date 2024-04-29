package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/RusGadzhiev/TaskManager/internal/config"
	"github.com/RusGadzhiev/TaskManager/internal/service"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrConnectingMySQL    = errors.New("error of connection mysql db")
	ErrPingMySQL          = errors.New("error of ping mysql db")
	ErrCreatingTableMySQL = errors.New("error of creating Tasks table")
)

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

	db.SetMaxOpenConns(10)
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error: %s, Description: %s", err, ErrPingMySQL)
	}

	query := `
		CREATE TABLE IF NOT EXISTS Tasks (
					id 			INT PRIMARY KEY AUTO_INCREMENT,
					owner 		TEXT, 
					executor 	TEXT,
					description TEXT,
					completed 	BOOL,
					assigned 	BOOL
		);

		CREATE INDEX IF NOT EXISTS idx_owner ON Tasks USING hash(
			owner
		);

		CREATE INDEX IF NOT EXISTS idx_executor ON links USING hash(
			executor
		);
	`

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("Error %s, Description: %s", err, ErrCreatingTableMySQL)
	}

	return &TasksRepoMySQL{DB: db}
}

func (repo *TasksRepoMySQL) Add(ctx context.Context, task *service.Task) (uint64, error) {
	res, err := repo.DB.ExecContext(ctx,
		"INSERT INTO Tasks (`owner`, `executor`, `description`, `completed`, `assigned`) VALUES (?, ?, ?, ?, ?)",
		task.Owner,
		task.Executor,
		task.Description,
		task.Completed,
		task.Assigned,
	)
	if err != nil {
		return 0, fmt.Errorf("insert mysql error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insert (last inserted ID) mysql error: %w", err)
	}
	return uint64(id), nil
}

func (repo *TasksRepoMySQL) GetAllTasks(ctx context.Context) ([]*service.Task, error) {
	return repo.getSomeTasks(ctx, service.FilterAllTasks, nil)
}

func (repo *TasksRepoMySQL) GetCreatedTasks(ctx context.Context, username string) ([]*service.Task, error) {
	return repo.getSomeTasks(ctx, service.FilterCreatedTasks, map[string]string{service.UserName: username})
}

func (repo *TasksRepoMySQL) GetMyTasks(ctx context.Context, username string) ([]*service.Task, error) {
	return repo.getSomeTasks(ctx, service.FilterMyTasks, map[string]string{service.UserName: username})
}

func (repo *TasksRepoMySQL) Assign(ctx context.Context, taskId uint64, username string) error {
	return repo.updateSth(ctx, service.FilterAssign, map[string]interface{}{service.TaskId: taskId, service.UserName: username})
}

func (repo *TasksRepoMySQL) Unassign(ctx context.Context, taskId uint64) error {
	return repo.updateSth(ctx, service.FilterUnassign, map[string]interface{}{service.TaskId: taskId})
}

func (repo *TasksRepoMySQL) Complete(ctx context.Context, taskId uint64) error {
	return repo.updateSth(ctx, service.FilterComplete, map[string]interface{}{service.TaskId: taskId})
}

func (repo *TasksRepoMySQL) getSomeTasks(ctx context.Context, filter string, args map[string]string) ([]*service.Task, error) {
	var rows *sql.Rows
	var err error
	switch filter {
	case service.FilterAllTasks:
		rows, err = repo.DB.QueryContext(ctx, "SELECT id, owner, executor, description, completed, assigned FROM Tasks")
	case service.FilterMyTasks:
		rows, err = repo.DB.QueryContext(ctx, "SELECT id, owner, executor, description, completed, assigned FROM Tasks WHERE executor=?", args[service.UserName])
	case service.FilterCreatedTasks:
		rows, err = repo.DB.QueryContext(ctx, "SELECT id, owner, executor, description, completed, assigned FROM Tasks WHERE owner=?", args[service.UserName])
	}
	if err != nil {
		return nil, fmt.Errorf("select mysql error: %w", err)
	}
	defer rows.Close()

	Tasks := []*service.Task{}
	for rows.Next() {
		Task := &service.Task{}
		err = rows.Scan(&Task.ID, &Task.Owner, &Task.Executor, &Task.Description, &Task.Completed, &Task.Assigned)
		if err != nil {
			return nil, fmt.Errorf("scanning mysql error: %w", err)
		}
		Tasks = append(Tasks, Task)
	}
	return Tasks, nil
}

func (repo *TasksRepoMySQL) updateSth(ctx context.Context, filter string, args map[string]interface{}) error {
	var err error
	switch filter {
	case service.FilterAssign:
		_, err = repo.DB.QueryContext(ctx, "UPDATE Tasks SET `executor` = ?, `assigned` = 1 WHERE id = ?", args[service.UserName], args[service.TaskId])
	case service.FilterUnassign:
		_, err = repo.DB.QueryContext(ctx, "UPDATE Tasks SET `executor` = \"\", `assigned` = 0 WHERE id = ?", args[service.TaskId])
	case service.FilterComplete:
		_, err = repo.DB.QueryContext(ctx, "UPDATE Tasks SET `completed` = 1 WHERE id = ?", args[service.TaskId])
	}
	if err != nil {
		return fmt.Errorf("update mysql error: %w", err)
	}
	return nil
}
