package api

import (
	_ "embed"
	"fmt"
	"net/http"

	"go-boilerplate/config"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ServePublicServer(cfg config.ServerConfig) {
	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
	)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		w.Write([]byte("boilerplate"))
	})

	httpPort := fmt.Sprintf(":%d", cfg.Port)
	go func() {
		if err := http.ListenAndServe(httpPort, r); err != nil {
			slog.Error("ServePublicServer http.ListenAndServe failed", "error", err)
			panic(err)
		}
	}()

	slog.Info("[HTTP]", "message", fmt.Sprintf("server is running at port: %d", cfg.Port))
}
