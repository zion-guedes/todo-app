package usecases

import "github.com/zion-guedes/todo-app/entities"

type TodosRepository interface {
	GetAllTodos() ([]entities.Todo, error)
}
