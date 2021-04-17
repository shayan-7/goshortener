package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shayan-7/goshortener/internal/db"
	"github.com/shayan-7/goshortener/internal/handlers"
	"github.com/shayan-7/goshortener/internal/models"
)

// TODO: Take ListenAddr and RedisAddr from CLI
var (
	ListenAddr = "localhost:8080"
	RedisAddr  = "localhost:6379"
	RedisChan  = "goshortener_stats"
)

func blue(v interface{}) string {
	return fmt.Sprintf("\x1b[34m%v\x1b[0m", v)
}

func main() {
	l := log.New(os.Stdout, "goshortener:", log.LstdFlags)
	r := db.NewRedis(RedisAddr)
	uh := handlers.NewURLHandler(l, r)

	err := db.InitDatabase()
	if err != nil {
		log.Fatalln("could not create database", err)
	}
	db.GlobalDB.AutoMigrate(&models.Member{})
	db.GlobalDB.AutoMigrate(&models.Stats{})

	switch true {
	case len(os.Args) == 2 && os.Args[1] == "serve":
		// Create new mux and register the handlers
		sm := mux.NewRouter()

		// Handle CORS
		cors := gHandlers.CORS(gHandlers.AllowedOrigins([]string{"*"}))

		routerGet := sm.Methods(http.MethodGet).Subrouter()
		routerGet.HandleFunc("/{id:.+}", uh.GetURLs)

		routerPost := sm.Methods(http.MethodPost).Subrouter()
		routerPost.HandleFunc("/", uh.AddURL)
		routerPost.HandleFunc("/signup", handlers.Signup)
		routerPost.HandleFunc("/login", handlers.Login)

		s := http.Server{
			Addr:         ListenAddr,
			Handler:      cors(sm),
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
			IdleTimeout:  12 * time.Second,
		}

		go func() {
			fmt.Printf(blue("I'm serving on %s\n"), ListenAddr)
			err := s.ListenAndServe()
			if err != nil {
				l.Fatalf("Server error: %v\n", err)
			}
		}()

		sigChan := make(chan os.Signal, 2)
		signal.Notify(sigChan, os.Interrupt)
		signal.Notify(sigChan, os.Kill)

		sig := <-sigChan
		l.Println("Recieved signal:", sig)

		tc, _ := context.WithTimeout(context.Background(), 7*time.Second)
		s.Shutdown(tc)

	case len(os.Args) == 2 && os.Args[1] == "subscribe":
		models.Subscribe(l, r, db.GlobalDB)

	default:
		fmt.Println("Only 'serve' and 'subscribe' is accepted as CLI argument")
	}
}
