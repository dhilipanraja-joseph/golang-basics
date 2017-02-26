package main

import (
  "fmt"
  "net/http"
  "html/template"
  "github.com/gorilla/mux"
)

func rootRoute(w http.ResponseWriter, r *http.Request) {
  templates := template.Must(template.ParseFiles("public/index.html"))
  templates.ExecuteTemplate(w, "index.html", nil)
}

func keyHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "keys: %v\n", vars["keys"])
}

func getTodos(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "Get all todos from dataBase")
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", rootRoute).Methods("GET") // Routes

  s := r.PathPrefix("/todos").Subrouter() // Create SubRoutes
  s.HandleFunc("/", getTodos).Methods("GET") // "/todos"
  s.HandleFunc("/{keys}", keyHandler).Methods("GET") // "/todos/{keys}"

  fmt.Println(http.ListenAndServe(":8000", r))  // Server Listener
}
