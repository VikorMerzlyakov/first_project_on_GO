package main

import (
	"testing"
)

func TestGenerateID(t *testing.T) {
	db := connectDB()
	defer db.Where("1 = 1").Delete(&Todo{})

	id := generateID(db)
	if id != "1" {
		t.Errorf("Ожидалось ID=1, получено %s", id)
	}

	var todo Todo
	todo = Todo{ID: id, Title: "Test", Done: false}
	db.Create(&todo)

	id2 := generateID(db)
	if id2 != "2" {
		t.Errorf("Ожидалось ID=2, получено %s", id2)
	}
}
