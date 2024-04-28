package service

import (
	"context"
)

type TasksStorage interface {
	GetAllTasks(ctx context.Context) ([]*Task, error)
	GetCreatedTasks(ctx context.Context, username string) ([]*Task, error)
	GetMyTasks(ctx context.Context, username string) ([]*Task, error)
	// возвращает id вставленной задачи
	Add(ctx context.Context, task *Task) (uint64, error)
	Assign(ctx context.Context, taskId uint64, username string) error
	Unassign(ctx context.Context, taskId uint64) error
	Complete(ctx context.Context, taskId uint64) error
}

type TasksService struct {
	repo TasksStorage
}

func NewTasksService(repo TasksStorage) *TasksService {
	return &TasksService{
		repo: repo,
	}
}

func (s *TasksService) GetAllTasks(ctx context.Context) ([]*Task, error) {
	tasks, err := s.repo.GetAllTasks(ctx)
	return tasks, err
}

func (s *TasksService) GetCreatedTasks(ctx context.Context, username string) ([]*Task, error) {
	tasks, err := s.repo.GetCreatedTasks(ctx, username)
	return tasks, err
}

func (s *TasksService) GetMyTasks(ctx context.Context, username string) ([]*Task, error) {
	tasks, err := s.repo.GetMyTasks(ctx, username)
	return tasks, err
}
// возвращает id вставленной задачи
func (s *TasksService) Add(ctx context.Context, task *Task) (uint64, error) {
	id, err := s.repo.Add(ctx, task)
	return id, err
}

func (s *TasksService) Assign(ctx context.Context, taskId uint64, username string) error {
	err := s.repo.Assign(ctx, taskId, username)
	return err
}

func (s *TasksService) Unassign(ctx context.Context, taskId uint64) error {
	err := s.repo.Unassign(ctx, taskId)
	return err
}

func (s *TasksService) Complete(ctx context.Context, taskId uint64) error {
	err := s.repo.Complete(ctx, taskId)
	return err
}
