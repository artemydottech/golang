package service

import (
	"testing"

	"tasks-crud/internal/domain"
	"tasks-crud/internal/repository"
)

func newTestTaskService() *TaskService {
	return NewTaskService(repository.NewInMemoryTaskRepository())
}

func TestGetAllTasksIsOrderedByID(t *testing.T) {
	svc := newTestTaskService()

	for _, title := range []string{"third", "fourth", "fifth"} {
		if _, err := svc.CreateTask(domain.CreateTaskRequest{Title: title}); err != nil {
			t.Fatalf("CreateTask(%q) failed: %v", title, err)
		}
	}

	for range 20 {
		tasks, err := svc.GetAllTasks()
		if err != nil {
			t.Fatalf("GetAllTasks failed: %v", err)
		}

		for i := 1; i < len(tasks); i++ {
			if tasks[i-1].ID > tasks[i].ID {
				t.Fatalf("tasks out of order: %d before %d", tasks[i-1].ID, tasks[i].ID)
			}
		}
	}
}

func TestCreateTaskTrimsAndRejectsBlankTitle(t *testing.T) {
	svc := newTestTaskService()

	task, err := svc.CreateTask(domain.CreateTaskRequest{Title: "  padded  "})
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	if task.Title != "padded" {
		t.Fatalf("title %q was not trimmed", task.Title)
	}

	if _, err := svc.CreateTask(domain.CreateTaskRequest{Title: "   "}); err == nil {
		t.Fatal("CreateTask accepted a blank title")
	}
}

func TestCreateTaskRejectsDuplicateTitle(t *testing.T) {
	svc := newTestTaskService()

	if _, err := svc.CreateTask(domain.CreateTaskRequest{Title: "unique"}); err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	if _, err := svc.CreateTask(domain.CreateTaskRequest{Title: "UNIQUE"}); err == nil {
		t.Fatal("CreateTask accepted a case-insensitive duplicate")
	}
}

func TestCreateTaskStampsTimestamps(t *testing.T) {
	svc := newTestTaskService()

	task, err := svc.CreateTask(domain.CreateTaskRequest{Title: "stamped"})
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	if task.CreatedAt.IsZero() {
		t.Error("created_at was left at the zero time")
	}
	if task.UpdatedAt.IsZero() {
		t.Error("updated_at was left at the zero time")
	}
}

func TestUpdateTaskTouchesUpdatedAtAndKeepsCreatedAt(t *testing.T) {
	svc := newTestTaskService()

	created, err := svc.CreateTask(domain.CreateTaskRequest{Title: "before"})
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	done := true
	updated, err := svc.UpdateTask(created.ID, domain.UpdateTaskRequest{Completed: &done})
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	if !updated.Completed {
		t.Error("completed flag was not applied")
	}
	if updated.Title != "before" {
		t.Errorf("title changed to %q when only completed was sent", updated.Title)
	}
	if !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Error("created_at was rewritten by an update")
	}
	if !updated.UpdatedAt.After(created.UpdatedAt) {
		t.Error("updated_at did not move forward")
	}
}

func TestUpdateTaskRejectsBlankTitleAndMissingID(t *testing.T) {
	svc := newTestTaskService()

	blank := "   "
	if _, err := svc.UpdateTask(1, domain.UpdateTaskRequest{Title: &blank}); err == nil {
		t.Error("UpdateTask accepted a blank title")
	}

	title := "ok"
	if _, err := svc.UpdateTask(9999, domain.UpdateTaskRequest{Title: &title}); err == nil {
		t.Error("UpdateTask accepted an unknown id")
	}

	if _, err := svc.UpdateTask(0, domain.UpdateTaskRequest{Title: &title}); err == nil {
		t.Error("UpdateTask accepted id 0")
	}
}

func TestDeleteTaskRemovesItOnce(t *testing.T) {
	svc := newTestTaskService()

	if err := svc.DeleteTask(1); err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	if err := svc.DeleteTask(1); err == nil {
		t.Fatal("DeleteTask succeeded twice for the same id")
	}

	if _, err := svc.GetTaskByID(1); err == nil {
		t.Fatal("GetTaskByID still returns a deleted task")
	}
}
