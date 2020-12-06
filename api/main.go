package main

import (
	"database/sql"
	"fmt"
	"log"

	"go.uber.org/fx"

	"context"

	"github.com/codegangsta/negroni"
	"github.com/eminetto/clean-architecture-go-v2/api/middleware"
	"github.com/eminetto/clean-architecture-go-v2/api/router"
	"github.com/eminetto/clean-architecture-go-v2/config"
	"github.com/eminetto/clean-architecture-go-v2/pkg/metric"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// NewSQLDB create and open database
func NewSQLDB(lc fx.Lifecycle) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}

	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				log.Println("Start NewSqlDB")
				return nil
			},
			OnStop: func(context.Context) error {
				log.Println("Stop NewSqlDB")
				db.Close()
				return nil
			},
		},
	)

	return db, err
}

// NewPrometheusService wrap NewPrometheusService for adding hooks
func NewPrometheusService(lc fx.Lifecycle) (metric.Service, error) {

	metricService, err := metric.NewPrometheusService()
	appMetric := metric.NewCLI("search")

	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				log.Println("Start NewPrometheusService")
				appMetric.Started()
				return nil
			},
			OnStop: func(context.Context) error {
				log.Println("Stop NewPrometheusService")
				appMetric.Finished()
				err = metricService.SaveCLI(appMetric)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	)
	return metricService, err
}

// core options for fx
func opts() fx.Option {
	return fx.Options(
		fx.Provide(
			NewSQLDB,
			NewPrometheusService,
			mux.NewRouter,
			func(metric metric.Service) negroni.Negroni {
				return *negroni.New(
					negroni.HandlerFunc(middleware.Cors),
					negroni.HandlerFunc(middleware.Metrics(metric)),
					negroni.NewLogger(),
				)
			},
		),

		config.ServiceConstructor,
		config.HandlerInvoker,

		fx.Invoke(router.Register),
	)
}

func main() {
	fx.New(opts()).Run()
}
