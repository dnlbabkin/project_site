package main

import (
  "fmt"
  "net/http"
  "html/template"
  "github.com/gorilla/mux"
  "database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type article struct {
  Id uint16
  Title, Anons, FullText string
}

var posts = []article {}
var showPst = article {}

func checkError(w http.ResponseWriter, err error) {
  if (err != nil) {
    fmt.Fprintf(w, err.Error())
  }
}

func index(w http.ResponseWriter, r *http.Request) {
  tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
  checkError(w, err)

  db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
  checkError(w, err)

  defer db.Close()

  res, err := db.Query("SELECT * FROM `articles`")
  checkError(w,   err)

  posts = []article{}

  for res.Next() {
    var article article
    err = res.Scan(&article.Id, &article.Title, &article.Anons, &article.FullText)

    checkError(w, err)

    posts = append(posts, article)
  }

  defer res.Close()

  tmpl.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
  tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
  checkError(w, err)

  tmpl.ExecuteTemplate(w, "create", nil)
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
  title := r.FormValue("title")
  anons := r.FormValue("anons")
  fullText := r.FormValue("full_text")

  if (title == "" || anons == "" || fullText == "") {
    fmt.Fprintf(w, "Э, чепушила, не все данные ты ввел")
  } else {
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
    checkError(w, err)

    defer db.Close()

    insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')", title, anons, fullText))
    checkError(w, err)

    defer insert.Close()

    http.Redirect(w, r, "/", http.StatusSeeOther)
  }
}

func showPost(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  w.WriteHeader(http.StatusOK)

  tmpl, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
  checkError(w, err)

  db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
  checkError(w, err)

  defer db.Close()

  res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
  checkError(w, err)

  showPst = article{}

  for res.Next() {
    var article article
    err = res.Scan(&article.Id, &article.Title, &article.Anons, &article.FullText)

    checkError(w, err)

    showPst = article
  }

  tmpl.ExecuteTemplate(w, "show", showPst)
}

func handleFunc() {
  rtr := mux.NewRouter()

  rtr.HandleFunc("/", index).Methods("GET")
  rtr.HandleFunc("/create", create).Methods("GET")
  rtr.HandleFunc("/save_article", saveArticle).Methods("POST")
  rtr.HandleFunc("/post/{id:[0-9]+}", showPost).Methods("GET")

  http.Handle("/", rtr)

  http.ListenAndServe(":8080", nil)
}

func main() {
  handleFunc()
}
