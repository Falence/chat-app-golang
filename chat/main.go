package main

import (
	"flag"
	"log"
	"net/http"

	// "os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	// "github.com/falence/go-blueprints/trace"
)

// temple represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application")
	flag.Parse() // parse the flags

	// setup gomniauth
	gomniauth.SetSecurityKey("my-auth-key")
	gomniauth.WithProviders(
		google.New(
			"310169518174-jjqlfedt8217illgaqp54o6029hljeid.apps.googleusercontent.com",
			"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",
			"http://localhost:8080/auth/callback/google",
		),
	)

	r := newRoom()
	// r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// get the room going
	go r.run()
	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
