package main

import (
	"encoding/json"
	"log"
	"luma/internal/api"
	"luma/internal/database"
	"luma/internal/repository"
	"net/http"
	"time"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatalf("Error al iniciar SQLite: %v", err)
	}
	defer db.Close()

	log.Println("SQLite conectada correctamente")

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Error al aplicar migraciones: %v", err)
	}

	log.Println("Esquema SQLite inicializado correctamente")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]string{
			"status":  "ok",
			"service": "luma",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error al enviar respuesta: %v", err)
		}
	})

	folderRepository := repository.NewFolderRepository(db)
	taskRepository := repository.NewTaskRepository(db)

	api.RegisterFolderRoutes(mux, folderRepository)
	api.RegisterFolderDetailRoutes(mux, folderRepository)
	api.RegisterTaskRoutes(mux, taskRepository)
	api.RegisterTaskDetailRoutes(mux, taskRepository)

	server := &http.Server{
		Addr:              "127.0.0.1:8081",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Luma API iniciada en http://127.0.0.1:8081")

	log.Fatal(server.ListenAndServe())
}
