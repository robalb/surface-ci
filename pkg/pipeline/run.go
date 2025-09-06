package pipeline

import (
	"context"
	"log/slog"
)

func RunSurfaceDiscovery(
	ctx context.Context,
	logger *slog.Logger,
	knownSurface *Surface,
	scope *Surface,
	scopeExclusion *Surface,
) {

	//we start with an empty pipeline
	pipeline := Surface{}

	exclusions := MakeExclusion()
	exclusions.Insert(scopeExclusion)

	insert_safe(*knownSurface, exclusions, &pipeline)
	insert_safe(*scope, exclusions, &pipeline)

	logger.Info("pipeline", "domains", pipeline.Domains)

	outDomains, err := Subfinder(ctx, pipeline.Domains)
	if err != nil {
		logger.Error("subfinder fail", "error", err)
		return
	}
	logger.Info("subfinder out: ", "domains", outDomains)

}
