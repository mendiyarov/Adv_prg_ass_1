package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// Глобальная переменная для базы данных
var DB *gorm.DB

// Структура пользователя
type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Email string `gorm:"unique"`
}

// Подключение к базе данных
func connectDatabase() {
	dsn := "user=postgres password=Dias140506 dbname=mydb2 port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to database successfully")

	// Автоматическая миграция
	DB.AutoMigrate(&User{})
}

// Обработчик для получения списка пользователей
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User
	DB.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Обработчик для создания нового пользователя
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request payload"})
		return
	}
	DB.Create(&user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Обработчик для обновления пользователя
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request payload"})
		return
	}

	var existingUser User
	if err := DB.First(&existingUser, user.ID).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "User not found"})
		return
	}

	existingUser.Name = user.Name
	existingUser.Email = user.Email
	DB.Save(&existingUser)

	json.NewEncoder(w).Encode(existingUser)
}

// Обработчик для удаления пользователя
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request payload"})
		return
	}

	if err := DB.Delete(&User{}, user.ID).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "User not found"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func main() {
	// Подключение к базе данных
	connectDatabase()

	// Настройка маршрутов
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsersHandler(w, r)
		case http.MethodPost:
			createUserHandler(w, r)
		case http.MethodPut:
			updateUserHandler(w, r)
		case http.MethodDelete:
			deleteUserHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed"})
		}
	})

	// Обслуживание статических файлов
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Запуск сервера
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
