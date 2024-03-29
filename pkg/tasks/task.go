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
	GetAllTasks(ctx context.Context) ([]*Task, error)                      // /tasks
	GetCreatedTasks(ctx context.Context, username string) ([]*Task, error) // owner
	GetMyTasks(ctx context.Context, username string) ([]*Task, error)      // /my
	Add(ctx context.Context, task *Task) error                             // /new XXX YYY ZZZ
	Assign(ctx context.Context, taskId uint64, username string) error      // assign_$ID`
	Unassign(ctx context.Context, taskId uint64) error                     // unassign_$ID
	Complete(ctx context.Context, taskId uint64) error                     // resolve_$ID`
}
