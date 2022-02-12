package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// set the active Avatar implementation
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application")
	flag.Parse() // parse the flags

	// setup gomniauth
	gomniauth.SetSecurityKey(os.Getenv("GOMNIAUTH_KEY"))
	gomniauth.WithProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_REDIRECT_URL"),
		),
		github.New(
			os.Getenv("GITHUB_CLIENT_ID"),
			os.Getenv("GITHUB_CLIENT_SECRET"),
			os.Getenv("GITHUB_REDIRECT_URL"),
		),
		facebook.New(
			os.Getenv("FACEBOOK_CLIENT_ID"),
			os.Getenv("FACEBOOK_CLIENT_SECRET"),
			os.Getenv("FACEBOOK_REDIRECT_URL"),
		),
	)

	// r := newRoom(UseAuthAvatar)
	// r := newRoom(UseGravatar)
	r := newRoom(avatars)
	// r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.Handle("/avatars/", http.StripPrefix(
		"/avatars",
		http.FileServer(http.Dir("./avatars"))),
	)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	
	// get the room going
	go r.run()
	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
