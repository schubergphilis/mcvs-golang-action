package main

import (
	"log/slog"
	"os"

	"github.com/schubergphilis/mcvs-golang-action/internal/app/mcvs_golang_action_bootstrap/application"
	"github.com/schubergphilis/mcvs-golang-action/internal/app/mcvs_golang_action_bootstrap/data"
	"github.com/schubergphilis/mcvs-golang-action/internal/app/mcvs_golang_action_bootstrap/presentation"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	storer := &data.FileStorer{Logger: logger}
	executor := &application.Downloader{Storer: storer, Logger: logger}
	presenter := &presentation.CLIPresenter{Executor: executor, Logger: logger}

	if err := presenter.Run(); err != nil {
		logger.Error("fatal", "err", err)
		os.Exit(1)
	}
}
