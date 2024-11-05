package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Server struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Status            string `json:"status"`
	StatusDescription string `json:"status_description"`
}

type Response struct {
	Name      string   `json:"name"`
	Result    []Server `json:"result"`
	Timestamp float64  `json:"timestamp"`
}

func GetServerList(db *sql.DB, ownerID int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var ownerName string
		err := db.QueryRow("SELECT name FROM owners WHERE id = ?", ownerID).Scan(&ownerName)
		if err != nil {
			log.Println("Ошибка при получении имени владельца:", err)
			http.Error(w, "Не удалось получить данные владельца", http.StatusInternalServerError)
			return
		}

		var servers []Server
		rows, err := db.Query("SELECT id, name, status FROM servers WHERE owner_id = ?", ownerID)
		if err != nil {
			log.Println("Ошибка сканирования сервера:", err)
			http.Error(w, "Не удалось проанализировать данные сервера", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var server Server
			if err := rows.Scan(&server.ID, &server.Name, &server.Status); err != nil {
				log.Println("Ошибка сканирования сервера:", err)
				http.Error(w, "Не удалось проанализировать данные сервера", http.StatusInternalServerError)
				return
			}
			servers = append(servers, server)
		}

		response := Response{
			Name:      ownerName,
			Result:    servers,
			Timestamp: float64(time.Now().UnixNano()) / 1e9,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		log.Println("Список серверов успешно получен и отправлен.", ownerID)
	}

}
