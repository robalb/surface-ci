// Package configFiles provides an abstraction layer between internal data structures and user-facing configuration files
package configfiles

import (
	"fmt"
	"path"

	"github.com/robalb/tinyasm/pkg/surface"
)


var (
	ScopeFileName  = "scope.yaml"
	ignoreFileName = "ignore-issues.yaml"
	asmconfigFileName = "asmconfig.yaml"
)

type ConfigFiles struct {
	Scope surface.Surface
	Exclusions surface.Surface
	//IgnoreIssues IgnoreIssues //TODO
	//Config Config //todo
}

func New(configFolder string) (*ConfigFiles, error){
	scopeFilePath := path.Join(configFolder, ScopeFileName)
	scopeFileData, err := parseScope(scopeFilePath)
	if err != nil {
		return nil, err
	}

	return &ConfigFiles{
		scopeFileData.Scope,
		scopeFileData.Exclusions,
	}, 
	nil
}

func (c *ConfigFiles) Summary() string{
	scope := fmt.Sprintf(
		"Elements in scope: {Domains[%d], IPs[%d], Endpoints[%d]}",
		len(c.Scope.Domains),
		len(c.Scope.IPs),
		len(c.Scope.URLs),
		)
	exclusions := fmt.Sprintf(
		"Elements excluded from scope: {Domains[%d], IPs[%d], Endpoints[%d]}",
		len(c.Scope.Domains),
		len(c.Scope.IPs),
		len(c.Scope.URLs),
		)
	return fmt.Sprintf("%s, %s ", scope, exclusions)
}
