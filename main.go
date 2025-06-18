package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Структура todo
type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var todos []Todo
var nextID = 1

// Обработчик для "/todos" GET
func todoGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Обработчик для "/todos" POST
func todoPost(w http.ResponseWriter, r *http.Request) {
	//Устанавливаю заголовок JSON
	w.Header().Set("Content-Type", "application/json")

	// Структура для получения данных от клиента
	var newTodo struct {
		Title string `json:"title"`
	}

	//Декодирую JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil || newTodo.Title == "" {
		http.Error(w, `{"error": "Неверный формат или отсутствует поле title"}`, http.StatusBadRequest)
		return
	}

	//Создаю новую задачу
	todo := Todo{
		ID:    fmt.Sprintf("%d", nextID),
		Title: newTodo.Title,
		Done:  false,
	}

	//Добавляем в список
	todos = append(todos, todo)
	nextID++

	//Возвращаем созданную задачу с кодом 201 Created
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func main() {
	// Регистрируем маршруты
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			todoGet(w, r)
		} else if r.Method == "POST" {
			todoPost(w, r)
		} else {
			http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	})

	// Запускаем сервер
	fmt.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
