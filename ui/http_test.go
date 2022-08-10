package ui_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gomagedon/expectate"
	"github.com/zion-guedes/todo-app/entities"
	"github.com/zion-guedes/todo-app/ui"
)

// MockService

type MockService struct {
	err   error
	todos []entities.Todo
}

func (s MockService) GetAllTodos() ([]entities.Todo, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.todos, nil
}

// dummyTodos

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

// HTTPTest

type HTTPTest struct {
	name           string
	service        *MockService
	inputMethod    string
	inputURL       string
	expectedStatus int
	expectedTodos  []entities.Todo
}

func TestHTTP(t *testing.T) {
	tests := getTests()
	tests = append(tests, getDisallowedMethodsTest()...)

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			testHTTP(t, test)
		})
	}
}

// Main Test

func testHTTP(t *testing.T, test HTTPTest) {
	expect := expectate.Expect(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(test.inputMethod, test.inputURL, nil)

	server := ui.NewHTTP()
	server.UseService(test.service)
	server.ServeHTTP(w, r)

	var body []entities.Todo
	json.NewDecoder(w.Result().Body).Decode(&body)

	expect(w.Result().StatusCode).ToBe(test.expectedStatus)
	expect(body).ToEqual(test.expectedTodos)
}

func getTests() []HTTPTest {
	return []HTTPTest{
		{
			name:           "Random error 500 status",
			service:        &MockService{err: fmt.Errorf("Somenthing was wrong")},
			inputMethod:    "GET",
			inputURL:       "http://checkr.com/todos",
			expectedStatus: 500,
		},
		{
			name:           "Wrong path app 404 status",
			service:        &MockService{todos: dummyTodos},
			inputMethod:    "GET",
			inputURL:       "http://checkr.com/foo",
			expectedStatus: 404,
		},
		{
			name:           "Wrong path app 404 status",
			service:        &MockService{todos: dummyTodos},
			inputMethod:    "GET",
			inputURL:       "http://checkr.com/bar",
			expectedStatus: 404,
		},
		{
			name:           "Returns todos from service.",
			service:        &MockService{todos: dummyTodos},
			inputMethod:    "GET",
			inputURL:       "http://checkr.com/todos/",
			expectedStatus: 200,
			expectedTodos:  dummyTodos,
		},
	}
}

func getDisallowedMethodsTest() []HTTPTest {
	tests := []HTTPTest{}

	disallowedMethods := []string{
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
	}

	for _, method := range disallowedMethods {
		tests = append(tests, HTTPTest{
			name:           fmt.Sprintf("Method %s gives 405 status", method),
			service:        &MockService{todos: dummyTodos},
			inputURL:       "http://checkr.com/todos/",
			inputMethod:    method,
			expectedStatus: http.StatusMethodNotAllowed,
		})
	}

	return tests
}
