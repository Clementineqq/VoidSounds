package main

import (
	"log"
	"net/http"

	"voidsounds/internal/components"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ping", pingHandler)

	log.Println("🚀 VoidSounds запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	component := components.Home()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<p class="text-emerald-400 mt-8">✅ Всё работает! HTMX активен.</p>`))
}
