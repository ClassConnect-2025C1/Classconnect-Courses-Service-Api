package main

import (
	"fmt"
	"net/http"
	"templateGo/controller"
)

func main() {
	mux := controller.SetupRoutes()
	port := "8080"
	fmt.Printf("Servidor escuchando en el puerto %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
