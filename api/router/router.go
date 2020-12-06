package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/eminetto/clean-architecture-go-v2/config"
	context2 "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	R          *mux.Router
}

// Register
func Register(p Params) error {

	http.Handle("/", p.R)
	http.Handle("/metrics", promhttp.Handler())
	p.R.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(config.API_PORT),
		Handler:      context2.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}

	p.Lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				logger.Println("Starting server.")
				go srv.ListenAndServe()
				// if err != nil {
				// 	log.Fatal(err.Error())
				// }

				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)

				// Block until a signal is received.
				go func() {
					s := <-c

					logger.Println(
						"Got signal.",
						zap.String("signal", s.String()),
					)

					if err := p.Shutdowner.Shutdown(); err != nil {
						logger.Fatal("Could not shutdown.", zap.Error(err))
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Println("Shutting down server.")
				return srv.Shutdown(ctx)
			},
		},
	)
	return nil
}
