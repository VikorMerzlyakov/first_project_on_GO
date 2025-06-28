package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// Структура todo
type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var todos []Todo
var nextID = 1

// сохраняю данные с сервера в файл
func saveTodosToFile(todos []Todo) error {
	file, err := os.Create("todos.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(todos)
}

// загружаю данные из файла
func loadTodosFromFile() ([]Todo, error) {
	file, err := os.Open("todos.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var todos []Todo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&todos)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

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

	err = saveTodosToFile(todos)
	if err != nil {
		http.Error(w, `{"error":"Не удалось сохранить данные"}`, http.StatusInternalServerError)
	}
}

// Обработчик PUT todos
func todoUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//получаем ID из URL
	id := r.URL.Path[len("/todos/"):]

	//структура для получения данных от клиента
	var updatedTodo struct {
		Title string `json:"title"`
		Done  bool   `json:"done"`
	}

	err := json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, `{"error": "Неверный формат"}`, http.StatusBadRequest)
		return
	}

	//Ищем задачу и обновляем
	found := false
	for i := range todos {
		if todos[i].ID == id {
			todos[i].Title = updatedTodo.Title
			todos[i].Done = updatedTodo.Done
			found = true
			break
		}
	}

	if !found {
		http.Error(w, `{"error": "Задача не  найдена"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(todos)

	err = saveTodosToFile(todos)
	if err != nil {
		http.Error(w, `{"error":"Не удалось сохранить данные"}`, http.StatusInternalServerError)
	}
}

// обработчик delete
func todoDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Получаем ID
	id := r.URL.Path[len("/todos/"):]

	//ищем индекс задачи
	index := -1
	for i := range todos {
		if todos[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}

	//удалчем задачу
	todos = append(todos[:index], todos[index+1:]...)

	//отправляем обновленный список
	json.NewEncoder(w).Encode(todos)

	err := saveTodosToFile(todos)
	if err != nil {
		http.Error(w, `{"error":"Не удалось сохранить данные"}`, http.StatusInternalServerError)
	}
}

func main() {
	//Загружаем данные из файла
	loadedTodos, err := loadTodosFromFile()
	if err == nil {
		todos = loadedTodos
		for _, t := range todos {
			id, _ := strconv.Atoi(t.ID)
			if id >= nextID {
				nextID = id + 1
			}
		}
	} else {
		fmt.Println("Не удалось загрузить задачи: ", err)
		todos = []Todo{}
		nextID = 1
	}
	// Регистрируем маршруты
	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			todoGet(w, r)
		} else if r.Method == "POST" {
			todoPost(w, r)
		} else if r.Method == "DELETE" {
			todoDelete(w, r)
		} else if r.Method == "PUT" {
			todoUpdate(w, r)
		} else {
			http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	})

	// Запускаем сервер
	fmt.Println("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
