package tasks

import (
	"github.com/RusGadzhiev/TaskManager/internal/config"
	"github.com/RusGadzhiev/TaskManager/internal/users"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrConnectingMySQL    = errors.New("error of connection mysql db")
	ErrPingMySQL          = errors.New("error of ping mysql db")
	ErrCreatingTableMySQL = errors.New("error of creating tasks table")
)

const (
	FilterAllTasks     = "AllTasks"
	FilterMyTasks      = "MyTasks"
	FilterCreatedTasks = "CreatedTasks"
	FilterAssign       = "Assign"
	FilterUnassign     = "Unassign"
	FilterComplete     = "Complete"
	TaskId             = "taskId"
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
	// строка ниже нужна чтобы база успела подняться, строка не нужна если контейнер task_manager в режиме restart always
	// time.Sleep(5 * time.Second) // выяснить почему без этого не работает пинг

	db.SetMaxOpenConns(10)
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error: %s, Description: %s", err, ErrPingMySQL)
	}

	query := `CREATE TABLE IF NOT EXISTS tasks(id int primary key auto_increment, owner text, 
	executor text, description text, completed bool, assigned bool)`
	// добавь индекс
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("Error %s, Description: %s", err, ErrCreatingTableMySQL)
	}

	return &TasksRepoMySQL{DB: db}
}

func (repo *TasksRepoMySQL) Add(ctx context.Context, task *Task) (uint64, error) {
	res, err := repo.DB.ExecContext(ctx,
		"INSERT INTO tasks (`owner`, `executor`, `description`, `completed`, `assigned`) VALUES (?, ?, ?, ?, ?)",
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

func (repo *TasksRepoMySQL) GetAllTasks(ctx context.Context) ([]*Task, error) {
	return repo.getSomeTasks(ctx, FilterAllTasks, nil)
}

func (repo *TasksRepoMySQL) GetCreatedTasks(ctx context.Context, username string) ([]*Task, error) {
	return repo.getSomeTasks(ctx, FilterCreatedTasks, map[string]string{users.UserName: username})
}

func (repo *TasksRepoMySQL) GetMyTasks(ctx context.Context, username string) ([]*Task, error) {
	return repo.getSomeTasks(ctx, FilterMyTasks, map[string]string{users.UserName: username})
}

func (repo *TasksRepoMySQL) Assign(ctx context.Context, taskId uint64, username string) error {
	return repo.updateSth(ctx, FilterAssign, map[string]interface{}{TaskId: taskId, users.UserName: username})
}

func (repo *TasksRepoMySQL) Unassign(ctx context.Context, taskId uint64) error {
	return repo.updateSth(ctx, FilterUnassign, map[string]interface{}{TaskId: taskId})
}

func (repo *TasksRepoMySQL) Complete(ctx context.Context, taskId uint64) error {
	return repo.updateSth(ctx, FilterComplete, map[string]interface{}{TaskId: taskId})
}

func (repo *TasksRepoMySQL) getSomeTasks(ctx context.Context, filter string, args map[string]string) ([]*Task, error) {
	var rows *sql.Rows
	var err error
	switch filter {
	case FilterAllTasks:
		rows, err = repo.DB.QueryContext(ctx, "SELECT id, owner, executor, description, completed, assigned FROM tasks")
	case FilterMyTasks:
		rows, err = repo.DB.QueryContext(ctx, "SELECT id, owner, executor, description, completed, assigned FROM tasks WHERE executor=?", args[users.UserName])
	case FilterCreatedTasks:
		rows, err = repo.DB.QueryContext(ctx, "SELECT id, owner, executor, description, completed, assigned FROM tasks WHERE owner=?", args[users.UserName])
	}
	if err != nil {
		return nil, fmt.Errorf("select mysql error: %w", err)
	}
	defer rows.Close()

	tasks := []*Task{}
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

func (repo *TasksRepoMySQL) updateSth(ctx context.Context, filter string, args map[string]interface{}) error {
	var err error
	switch filter {
	case FilterAssign:
		_, err = repo.DB.QueryContext(ctx, "UPDATE tasks SET `executor` = ?, `assigned` = 1 WHERE id = ?", args[users.UserName], args[TaskId])
	case FilterUnassign:
		_, err = repo.DB.QueryContext(ctx, "UPDATE tasks SET `executor` = \"\", `assigned` = 0 WHERE id = ?", args[TaskId])
	case FilterComplete:
		_, err = repo.DB.QueryContext(ctx, "UPDATE tasks SET `completed` = 1 WHERE id = ?", args[TaskId])
	}
	if err != nil {
		return fmt.Errorf("update mysql error: %w", err)
	}
	return nil
}
