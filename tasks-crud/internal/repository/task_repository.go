package repository

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"tasks-crud/internal/domain"
)

type TaskRepository interface {
	GetAll() ([]domain.Task, error)
	GetByID(id int) (*domain.Task, error)
	Create(task *domain.Task) error
	Update(id int, task *domain.Task) error
	Delete(id int) error
}

type InMemoryTaskRepository struct {
	tasks     map[int]domain.Task
	currentID int
	mu        sync.RWMutex
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	repo := &InMemoryTaskRepository{
		tasks:     make(map[int]domain.Task),
		currentID: 1,
	}

	seededAt := time.Now()
	repo.tasks[1] = domain.Task{
		ID:        1,
		Title:     "Выучить основы Go",
		Completed: false,
		CreatedAt: seededAt,
		UpdatedAt: seededAt,
	}
	repo.tasks[2] = domain.Task{
		ID:        2,
		Title:     "Написать первое API",
		Completed: true,
		CreatedAt: seededAt,
		UpdatedAt: seededAt,
	}
	repo.currentID = 3

	return repo
}

func (r *InMemoryTaskRepository) GetAll() ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	taskList := make([]domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		taskList = append(taskList, task)
	}
	sort.Slice(taskList, func(i, j int) bool { return taskList[i].ID < taskList[j].ID })

	return taskList, nil
}

func (r *InMemoryTaskRepository) GetByID(id int) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with id %d not found", id)
	}

	return &task, nil
}

func (r *InMemoryTaskRepository) Create(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	task.ID = r.currentID
	task.CreatedAt = now
	task.UpdatedAt = now

	r.tasks[r.currentID] = *task

	r.currentID++

	return nil
}

func (r *InMemoryTaskRepository) Update(id int, updatedTask *domain.Task) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.tasks[id]
	if !exists {
		return fmt.Errorf("task with id %d not found", id)
	}

	r.tasks[id] = *updatedTask

	return nil
}

func (r *InMemoryTaskRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.tasks[id]
	if !exists {
		return fmt.Errorf("task with id %d not found", id)
	}

	delete(r.tasks, id)

	return nil
}
