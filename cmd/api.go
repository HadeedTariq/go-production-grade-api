package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/HadeedTariq/go-production-grade-api/internal/adapters/postgresql/sqlc"
	"github.com/HadeedTariq/go-production-grade-api/internal/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	config config
	db     *pgxpool.Pool
}

type config struct {
	addr     string
	dbConfig dbConfig
}

type dbConfig struct {
	dsn string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // import for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})
	r.Route("/auth", func(authRoute chi.Router) {
		authService := auth.NewService(repo.New(app.db), app.db)
		authHandler := auth.NewHandler(authService)
		authRoute.Post("/register", authHandler.RegisterUser)
		authRoute.Get("/verification", authHandler.VerifyUser)
	})

	return r
}

func (app *application) run(handler http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}
