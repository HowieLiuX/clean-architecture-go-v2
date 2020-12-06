package config

import (
	"database/sql"

	"github.com/eminetto/clean-architecture-go-v2/infrastructure/repository"
	"github.com/eminetto/clean-architecture-go-v2/usecase/book"
	"github.com/eminetto/clean-architecture-go-v2/usecase/loan"
	"github.com/eminetto/clean-architecture-go-v2/usecase/user"
	"go.uber.org/fx"
)

// ServiceConstructor used for registing service instance
var ServiceConstructor = fx.Options(
	fx.Provide(
		func(db *sql.DB) book.Repository {
			return repository.NewBookMySQL(db)
		},
		func(db *sql.DB) user.Repository {
			return repository.NewUserMySQL(db)
		},
		func(r book.Repository) book.UseCase {
			return book.NewService(r)
		},
		func(r user.Repository) user.UseCase {
			return user.NewService(r)
		},
		func(u user.UseCase, b book.UseCase) loan.UseCase {
			return loan.NewService(u, b)
		},
	),
)
