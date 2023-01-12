package infrastructure

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/scanner/app/interfaces"
	"go.uber.org/zap"
)

// Dispatch is handle routing
func Dispatch(logger *zap.Logger, sqlHandler interfaces.SQLHandler) {
	repoController := interfaces.NewRepoController(sqlHandler, logger)
	scanController := interfaces.NewScanController(sqlHandler, logger)
	r := chi.NewRouter()
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/api", func(r chi.Router) {
		r.Get("/repos", repoController.Index)
		r.Route("/repo", func(r chi.Router) {
			r.Post("/", repoController.Create)
			r.Get("/{repoID}", repoController.Show)
			r.Put("/{repoID}", repoController.Update)
			r.Delete("/{repoID}", repoController.Delete)
			r.Post("/{repoID}/scan", scanController.Scan)
		})

		r.Route("/scan", func(r chi.Router) {
			r.Get("/results", scanController.Index)
			r.Get("/result/{resultID}", scanController.Show)
		})
	})

	logger.Info(fmt.Sprintf("Starting server at port %s\n", os.Getenv("SERVER_PORT")))

	if err := http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), r); err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}

}
