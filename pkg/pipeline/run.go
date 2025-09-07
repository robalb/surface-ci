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
	// pipeline ideas:
	// at the end of the discovery, resolve all domains to ips, one by one.
	// if a domain matches with an excluded ip, add it to the exclusions

	exclusions := MakeExclusion()
	exclusions.Insert(scopeExclusion)

	pipeline := Surface{}
	insert_safe(*knownSurface, exclusions, &pipeline)
	insert_safe(*scope, exclusions, &pipeline)

	logger.Info("pipeline - after insert", "domains", pipeline.Domains)

	// expand scope from urls
	{
		extractedDomains := URLExtractDomains(pipeline.URLs)
		insert_safe_string(extractedDomains, exclusions.Contains_domain, &pipeline.Domains)

		extractedIPs := URLExtractIPs(pipeline.URLs)
		insert_safe_string(extractedIPs, exclusions.Contains_ip, &pipeline.IPs)
	}

	// expand domains
	{
		// Filter subdomains of lower hierarchies before passing them to subfinder:
		// if the list contains bb.a.example.com, cc.a.example.com and a.example.com
		// we can assume that the whole a.example.com is in scope, and we can remove
		// all subdomains of a.example.com from the list since they would return the
		// same results
		filteredDomains, err := FilterSubdomains(pipeline.Domains)
		if err != nil {
			logger.Error("filterSubdomains fail", "error", err)
			return
		}
		logger.Info("pipeline - after filters", "domains", filteredDomains)

		outDomains, err := Subfinder(ctx, filteredDomains)
		if err != nil {
			logger.Error("subfinder fail", "error", err)
			return
		}

		insert_safe_string(outDomains, exclusions.Contains_domain, &pipeline.Domains)
		logger.Info("pipeline - subfinder", "domains", outDomains)
	}

	//fuzzy search domains
	{
		//fuzzy generate domain names, based on alterx and LLM prompts
		//insert into our dns pipeline only domains that resolve to something
		//even when domains resolve to something, make sure there are no wildcard dns 
		//before insering a bunch of allucinations. it can be tested by resolving random-letter subdomains
	}



}
