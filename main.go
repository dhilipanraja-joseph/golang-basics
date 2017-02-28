package main

import (
  "fmt"
  "net/http"
  "html/template"
  "encoding/json"
  //"io/ioutil"

  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Todo struct {
  Id    bson.ObjectId   `bson:"_id,omitempty" json:"-"`
  task  string          `bson:"task" json:"task"`
}

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

func addTodo(w http.ResponseWriter, r *http.Request) {
    s := dbSession()
    var d Todo
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&d)
    // r.ParseForm()
    // for key, _ := range r.Form {
    //   json.Unmarshal([]byte(key),&d)
    // }
    // body, err := ioutil.ReadAll(r.Body)
    // json.Unmarshal(body,&d)
    // fmt.Println(json.NewDecoder(r.Body))
    // err := d.Decode(&todo)
    // if err != nil {
    //   fmt.Println("cannot add:",err,todo)
    //  return
    // }
    fmt.Println(decoder,d, err)
    defer r.Body.Close()
    defer s.Close()
    // session := s.Copy()
    // session.DB("goTodos").C("todos").Insert(todo)
    // w.WriteHeader(http.StatusOK)
}

func dbSession() *mgo.Session {
  s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
  return s
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", rootRoute).Methods("GET") // Routes

  s := r.PathPrefix("/todos").Subrouter() // Create SubRoutes
  s.HandleFunc("", getTodos).Methods("GET") // "/todos"
  s.HandleFunc("", addTodo).Methods("POST") // Create todos
  s.HandleFunc("/{keys}", keyHandler).Methods("GET") // "/todos/{keys}"

  fmt.Println(http.ListenAndServe(":8000", r))  // Server Listener
}
