package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=1234 dbname=todos_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключится к БД")
	}

	db.AutoMigrate(&Todo{})
	return db
}

// Структура todo
type Todo struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var nextID = 1

func generateID(db *gorm.DB) string {
	//var maxID int64 = 1
	var last Todo
	result := db.Order("id DESC").First(&last)
	if result.Error != nil {
		return "1"
	}

	id, err := strconv.ParseInt(last.ID, 10, 64)
	if err != nil || id < 1 {
		return "1"
	}

	return strconv.FormatInt(id+1, 10)
}

// Обработчик для "/todos" GET
func todoGet(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var todos []Todo
		db.Find(&todos)
		json.NewEncoder(w).Encode(todos)
	}
}

// Обработчик для "/todos" POST
func todoPost(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var newTodo struct {
			Title string `json:"title"`
		}

		err := json.NewDecoder(r.Body).Decode(&newTodo)
		if err != nil || newTodo.Title == "" {
			http.Error(w, `{"error": "Неверный формат или отсутствует поле title"}`, http.StatusBadRequest)
			return
		}

		todo := Todo{
			ID:    generateID(db),
			Title: newTodo.Title,
			Done:  false,
		}

		db.Create(&todo)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)
	}
}

// Обработчик PUT todos
func todoUpdate(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		id := r.URL.Path[len("/todos/"):]
		var updatedTodo struct {
			Title string `json:"title"`
			Done  bool   `json:"done"`
		}

		err := json.NewDecoder(r.Body).Decode(&updatedTodo)
		if err != nil {
			http.Error(w, `{"error": "Неверный формат запроса"}`, http.StatusBadRequest)
			return
		}

		var todo Todo
		result := db.First(&todo, "id = ?", id)
		if result.Error != nil {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			return
		}

		todo.Title = updatedTodo.Title
		todo.Done = updatedTodo.Done
		db.Save(&todo)

		json.NewEncoder(w).Encode(todo)
	}
}

// обработчик delete
func todoDelete(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/todos/"):]
		var todo Todo

		result := db.Where("id = ?", id).Delete(&todo)
		if result.Error != nil || result.RowsAffected == 0 {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}

func main() {
	db := connectDB()
	// Регистрируем маршруты
	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			todoGet(db)(w, r)
		} else if r.Method == "POST" {
			todoPost(db)(w, r)
		} else if r.Method == "DELETE" {
			todoDelete(db)(w, r)
		} else if r.Method == "PUT" {
			todoUpdate(db)(w, r)
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
