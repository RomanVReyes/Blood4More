package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"roman-sangre/internal/database"
	"roman-sangre/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: No se encontró el archivo .env, usando variables del sistema")
	}

	// 2. Conectar a la base de datos
	database.ConnectDB()

	// 3. Configurar servidor (Ruta de prueba)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Petición recibida en ROOT. Path: %s", r.URL.Path)

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "templates/index.html")
	})

	http.HandleFunc("/donante", handlers.RedirectDonorEntry)
	http.HandleFunc("/donante/auth", handlers.ShowDonorAuth)
	http.HandleFunc("/donante/login", handlers.ShowDonorLogin)
	http.HandleFunc("/donante/registro", handlers.ShowDonorRegister)
	http.HandleFunc("/donante/dashboard", handlers.AuthMiddleware(handlers.ShowDonorDashboard))
	http.HandleFunc("/donante/logout", handlers.LogoutHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Sistema operando correctamente")
	})

	fmt.Printf("Servidor iniciado en http://localhost:%s\n", port)

	// 4. Iniciar el servidor
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error iniciando el servidor:", err)
	}
}
