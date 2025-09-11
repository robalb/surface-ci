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
	dnsCache := NewDNSCache()

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
		// Remove subdomains of lower hierarchies before passing them to subfinder:
		// if the list contains bb.a.example.com, cc.a.example.com and a.example.com
		// we can assume that the whole a.example.com is in scope, and we can remove
		// all subdomains of a.example.com from the list since they would return the
		// same results
		filteredDomains, err := TrimSubdomains(pipeline.Domains)
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

	wildcards := DnsxFilterWildcards(pipeline.Domains, dnsCache)

	//fuzzy search domains
	{
		//fuzzy generate domain names, based on alterx and LLM prompts
		//insert into our dns pipeline only domains that resolve to something.
		//even when domains resolve to something, make sure there are no wildcard dns
		//to avoid false positives. it can be tested by resolving random strings
		unfuzzableDomains := SelectSubdomains(pipeline.Domains, wildcards)
		fuzzableDomains := Subtract(pipeline.Domains, unfuzzableDomains)

		fuzzDomains, err := Alterx(fuzzableDomains)
		if err != nil {
			logger.Error("alterx fail", "error", err)
			return
		}
		logger.Info("pipeline - after fuzz", "domains", fuzzDomains)

		// exclude from our validation all fuzz domains that are part of wildcards domains,
		// since they are untestable
		logger.Info("pipeline - widcards", "domains", wildcards)
		untestableFuzzDomains := SelectSubdomains(fuzzDomains, wildcards)
		logger.Info("pipeline - untestable", "domains", untestableFuzzDomains)
		fuzzDomains = Subtract(fuzzDomains, untestableFuzzDomains)
		logger.Info("pipeline - testable", "domains", fuzzDomains)

		// filter domains that resolve to an ip
		validFuzzed := DnsxFilterActive(fuzzDomains, dnsCache)
		logger.Info("pipeline - fuzz active dns", "domains", validFuzzed)
		insert_safe_string(validFuzzed, exclusions.Contains_domain, &pipeline.Domains)

	}

	// expand domains, ips, urls, list into active urls 
	{
		// we must run httpx two times: one with a surface set that only contains wildcard domains,
		// and one with a surface set that does not contain wildcard domains.
		// in NOwildcard mode, an http response is considered "discovered surface", and its url is put into the surface urls.
		// in wildcard mode, all children of a specific wildcard are tested together, and only the domain 
		// that receive a response deviating from the median response will be considered "discovered surface"
		results, err := Httpx(pipeline, 2)
		if err != nil {
			logger.Error("alterx fail", "error", err)
			return
		}
		logger.Info("httpx results", "results", results)

	}

}
