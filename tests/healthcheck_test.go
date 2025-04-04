package test

import (
	"io/ioutil"
	"net/http"
	"templateGo/controller"
	"testing"
	"time"
)

// Para evitar correr múltiples instancias en tests concurrentes,
// se asume que los tests se ejecutan de manera secuencial.
func TestHealthcheck(t *testing.T) {
	mux := controller.SetupRoutes()
	port := "8080"

	// Iniciamos ListenAndServe en una gorutina para no bloquear el test.
	go func() {
		err := http.ListenAndServe(":"+port, mux)
		if err != nil {
			t.Errorf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Esperar un momento para que el servidor se levante.
	time.Sleep(100 * time.Millisecond)

	// Test "/healthcheck" endpoint
	resp, err := http.Get("http://localhost:8080/healthcheck")
	if err != nil {
		t.Fatalf("Error haciendo GET en '/healthcheck': %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Para '/healthcheck' se esperaba status 200 y se obtuvo: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Para '/healthcheck' se esperaba Content-Type 'application/json' y se obtuvo: %s", contentType)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error leyendo la respuesta de '/healthcheck': %v", err)
	}
	expected := `{"status":"ok"}`
	if string(bodyBytes) != expected {
		t.Errorf("La respuesta de '/healthcheck' no coincide, se esperaba: %s y se obtuvo: %s", expected, string(bodyBytes))
	}
}
