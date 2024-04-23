package tasks

import "context"

type Task struct {
	ID          uint64
	Owner       string
	Executor    string
	Description string
	Completed   bool
	Assigned    bool
}

type TasksRepo interface {
	GetAllTasks(ctx context.Context) ([]*Task, error)
	GetCreatedTasks(ctx context.Context, username string) ([]*Task, error)
	GetMyTasks(ctx context.Context, username string) ([]*Task, error)
	// возвращает id вставленной задачи
	Add(ctx context.Context, task *Task) (uint64, error)
	Assign(ctx context.Context, taskId uint64, username string) error
	Unassign(ctx context.Context, taskId uint64) error
	Complete(ctx context.Context, taskId uint64) error
}
