package entrypoints

import (
	"context"
	"io"
	"log/slog"

	"github.com/robalb/tinyasm/pkg/configfiles"
	"github.com/robalb/tinyasm/pkg/datafiles"
	"github.com/robalb/tinyasm/pkg/envconfig"
)

func Validate(
	ctx context.Context,
	stdout io.Writer,
	stderr io.Writer,
	args []string,
	getenv func(string) string,
) error {
	logger := slog.New(slog.NewTextHandler(stdout, nil))
	logger.Info("Validating TinyASM file configs")

	envConfig, err := envconfig.New(args, getenv, logger)
	if err != nil {
		logger.Error("Failed to parse all the environment variables", "error", err)
		return err
	}
	logger.Info("Config folder", "path", envConfig.ConfigFolder)
	logger.Info("Data folder", "path", envConfig.DataFolder)

	// Read all the configuration files, based on the path set in the ENV variables
	_, err = configfiles.New(envConfig.ConfigFolder)
	if err != nil {
		logger.Error("Failed to parse all the configuration files", "error", err)
		return err
	}
	logger.Info("Configuration files: OK")

	// Read all the data files, based on the path set in the ENV variables
	_, fileMissing, err := datafiles.New(envConfig.DataFolder)
	if err != nil {
		logger.Error("Failed to access or parse the data folder content", "error", err)
		return err
	}
	if fileMissing {
		logger.Warn("Some files in the data folder were missing or empty.",
			"suggestion", "If this is not the fist execution, make sure the data folder is being saved properly.")
	}
	logger.Info("Data files: OK")

	logger.Info("Everything is OK. quitting.")
	return nil
}
