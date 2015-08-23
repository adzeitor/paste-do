package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

import "github.com/shurcooL/github_flavored_markdown"

func todoHandler(storage Storage) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("templates/todo.html")
		id := r.URL.Path[len("/todo/"):]
		val := storage.Get(id)
		markdown := string(github_flavored_markdown.Markdown(([]byte)(val.Content)))

		out := struct {
			ID              string
			Content         string
			Visits          uint64
			CreatedAt       time.Time
			UpdatedAt       time.Time
			MarkdownContent template.HTML
		}{
			id,
			(*val).Content,
			(*val).Visits,
			(*val).CreatedAt,
			(*val).UpdatedAt,
			template.HTML(markdown),
		}
		tmpl.Execute(w, out)
	}
	return fn
}

func postNewHandler(storage Storage) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("content")
		id := storage.New(body)
		http.Redirect(w, r, "/todo/"+id, http.StatusFound)
	}
	return fn
}

func postEditHandler(storage Storage) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/edit/"):]
		body := r.FormValue("content")
		storage.Edit(id, body)
		http.Redirect(w, r, "/todo/"+id, http.StatusFound)
	}
	return fn
}

func indexHandler(storage Storage) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, nil)
	}
	return fn
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var storage Storage
	dbBackendType := os.Getenv("BACKEND")
	if dbBackendType == "redis" {
		log.Println("loading redis backend...")
		u, _ := url.Parse(os.Getenv("REDIS_URL"))
		pass, _ := u.User.Password()
		host := u.Host
		storage = NewRedisStorage(host, pass)
	} else {
		log.Println("loading memory backend...")
		storage = NewMemoryStorage("???")
	}

	http.HandleFunc("/todo/", todoHandler(storage))
	http.HandleFunc("/edit/", postEditHandler(storage))
	http.HandleFunc("/new/", postNewHandler(storage))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", indexHandler(storage))

	log.Printf("listen %s...", port)
	http.ListenAndServe(":"+port, nil)
}
