package main

import (
	"context"
	"os"

	"github.com/robalb/tinyasm/internal/entrypoints"
)

func main() {
	ctx := context.Background()
	err := entrypoints.Validate(ctx, os.Stdout, os.Stderr, os.Args, os.Getenv)
	if err != nil {
		os.Exit(1)
	}
}
