package entrypoints

import (
	"context"
	"io"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/robalb/tinyasm/pkg/configfiles"
	"github.com/robalb/tinyasm/pkg/datafiles"
	"github.com/robalb/tinyasm/pkg/envconfig"
	"github.com/robalb/tinyasm/pkg/pipeline"
)

func Asm(
	ctx context.Context,
	stdout io.Writer,
	stderr io.Writer,
	args []string,
	getenv func(string) string,
) error {
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT,  // ctr-C from the terminal
		syscall.SIGTERM, // terminate signal from Docker / kubernetes / CI pipelines
	)
	defer cancel()

	// Initialize logging.
	// This program will run in Containerize environemnts or CI pipelines, where
	// log output must contain human-readable, useful information about the
	// program runtime.
	// The first steps here will therefore focus on logging contextual information on
	// what program is starting, the version, and the configuration parameters in use
	logger := slog.New(slog.NewTextHandler(stdout, nil))
	logger.Info("Starting TinyASM")

	// Read the configuration that can be set in ENV variables
	envConfig, err := envconfig.New(args, getenv, logger)
	if err != nil {
		logger.Error("Failed to parse all the environment variables", "error", err)
		return err
	}
	logger.Info("Config folder", "path", envConfig.ConfigFolder)
	logger.Info("Data folder", "path", envConfig.DataFolder)

	// Read all the configuration files, based on the path set in the ENV variables
	configFiles, err := configfiles.New(envConfig.ConfigFolder)
	if err != nil {
		logger.Error("Failed to parse all the configuration files", "error", err)
		return nil
	}
	logger.Info("file configuration values", "summary", configFiles.Summary())

	// Read all the data files, based on the path set in the ENV variables
	dataFiles, fileMissing, err := datafiles.New(envConfig.DataFolder)
	if err != nil {
		logger.Error("Failed to access or parse the data folder content", "error", err)
		return nil
	}
	logger.Info("data folder values", "summary", dataFiles.Summary())
	if fileMissing {
		logger.Warn("Some files in the data folder were missing or empty.",
			"suggestion", "If this is not the fist execution, make sure the data folder is being saved properly.")
	}

	pipeline.RunSurfaceDiscovery(
		ctx,
		logger,
		&dataFiles.KnownSurface,
		&configFiles.Scope,
		&configFiles.Exclusions,
	)

	return nil
}
