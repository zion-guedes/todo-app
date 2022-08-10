package ui

import "github.com/zion-guedes/todo-app/entities"

type Service interface {
	GetAllTodos() ([]entities.Todo, error)
}
