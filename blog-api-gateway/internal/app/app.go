package app

import (
	"github.com/mtvy/blog-api-gateway/config"
	"github.com/mtvy/blog-api-gateway/internal/handler"
	"github.com/mtvy/blog-api-gateway/internal/logger"
	"github.com/mtvy/blog-api-gateway/internal/repository"
	"github.com/mtvy/blog-api-gateway/internal/usecase"
	"github.com/mtvy/blog-api-gateway/migrations"
	"github.com/pkg/errors"
)

func Run() error {
	cfg, err := config.Parse()
	if err != nil {
		return errors.Wrap(err, "parse cfg")
	}

	logger.Register(cfg.Log)

	repo := repository.NewPostProvider()
	uc := usecase.NewPostProvider(repo)
	handle := handler.New(uc)

	if err := migrations.Migrate(repo); err != nil {
		return errors.Wrap(err, "migrations")
	}

	if err := getRouter(handle).Listen(cfg.GetHTTPEndpoint()); err != nil {
		return errors.Wrap(err, "server listen")
	}
	return nil
}
