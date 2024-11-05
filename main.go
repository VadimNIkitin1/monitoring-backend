package main

import (
	"database/sql"
	"log"
	"net/http"

	"goApp/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/monitoring.db")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v\n", err)
	}
	defer db.Close()

	log.Println("Сервер запускается на порту 8080...")

	http.HandleFunc("/servers_list", handlers.GetServerList(db, 1))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v\n", err)
	}
}
