package usecases_test

import (
	"fmt"
	"testing"

	"github.com/gomagedon/expectate"
	"github.com/zion-guedes/todo-app/entities"
	"github.com/zion-guedes/todo-app/usecases"
)

var dummyTodos = []entities.Todo{
	{
		Title:       "Task 1",
		Description: "Descrpt task 1",
		IsCompleted: true,
	},
	{
		Title:       "Task 2",
		Description: "Descrpt task 2",
		IsCompleted: false,
	},
	{
		Title:       "Task 3",
		Description: "Descrpt task 3",
		IsCompleted: true,
	},
}

type BadTodoRepo struct{}

type MockTodoRepo struct {
}

func (BadTodoRepo) GetAllTodos() ([]entities.Todo, error) {
	return nil, fmt.Errorf("Something went wrong")
}

func (MockTodoRepo) GetAllTodos() ([]entities.Todo, error) {
	return dummyTodos, nil
}

func TestGetTodos(t *testing.T) {
	// Test
	t.Run("Returns ErrInternal when TodosRepo error", func(t *testing.T) {
		expect := expectate.Expect(t)
		repo := new(BadTodoRepo)
		todos, err := usecases.GetTodos(repo)

		expect(err).ToBe(usecases.ErrInternal)
		//expect(todos).ToBe(nil)
		if todos != nil {
			t.Fatalf("Expected todos to be nil")
		}
	})

	// Test
	t.Run("Returns todos from TodosRepo", func(t *testing.T) {
		expect := expectate.Expect(t)

		repo := new(MockTodoRepo)

		todos, err := usecases.GetTodos(repo)

		expect(err).ToBe(nil)
		expect(todos).ToEqual(dummyTodos)
	})
}
