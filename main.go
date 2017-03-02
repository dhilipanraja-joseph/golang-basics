package main

import (
  "fmt"
  "net/http"
  "html/template"
  "encoding/json"
  "strings"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Todo struct {
  Id    bson.ObjectId   `bson:"_id,omitempty" json:"_id"`
  Task  string          `bson:"task" json:"task"`
}

func rootRoute(w http.ResponseWriter, r *http.Request) {
  templates := template.Must(template.ParseFiles("public/index.html"))
  templates.ExecuteTemplate(w, "index.html", nil)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
  s := dbSession()
  defer s.Close()
  getJSONresp(w)
}

// POST METHOD - to Create New Todo
func addTodo(w http.ResponseWriter, r *http.Request) {
  s := dbSession()
  r.ParseForm()
  task := strings.Join(r.Form["Task"], "")
  todo := Todo{Task:task}
  defer r.Body.Close()
  defer s.Close()
  db := s.DB("goTodos").C("todos")
  err := db.Insert(todo)
  if err != nil {
    panic(err)
  }
  getJSONresp(w)
}

// DELETE METHOD - TO Delete A Todo By ID
func deleteTodo(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id := bson.ObjectIdHex(vars["id"])
  s := dbSession()
  defer s.Close()
  db := s.DB("goTodos").C("todos")
  err := db.Remove(bson.M{ "_id": id })
  if err != nil {
    panic(err)
  }
  getJSONresp(w)
}

// Generate JSON Response with all Todos - Getall
func getJSONresp(w http.ResponseWriter) {
  var t []Todo
  s:= dbSession()
  err := s.DB("goTodos").C("todos").Find(bson.M{}).All(&t)
  if err != nil {
    panic(err)
  }
  resB, er := json.MarshalIndent(t, "", "   ")
  if er != nil {
    panic(er)
  }
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(resB)
  w.WriteHeader(http.StatusOK)

}

func dbSession() *mgo.Session {
  s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
  return s.Copy()
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", rootRoute).Methods("GET") // Routes

  s := r.PathPrefix("/todos").Subrouter() // Create SubRoutes
  s.HandleFunc("", getTodos).Methods("GET") // "/todos"
  s.HandleFunc("", addTodo).Methods("POST") // Create todos
  s.HandleFunc("/{id}", deleteTodo).Methods("DELETE")

  fmt.Println(http.ListenAndServe(":8000", r))  // Server Listener
}
