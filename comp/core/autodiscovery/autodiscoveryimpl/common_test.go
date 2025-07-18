// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package autodiscoveryimpl

import (
	"reflect"

	"github.com/DataDog/datadog-agent/comp/core/autodiscovery/integration"
	"github.com/DataDog/datadog-agent/comp/core/autodiscovery/listeners"
	filter "github.com/DataDog/datadog-agent/comp/core/workloadfilter/def"
)

type dummyService struct {
	ID              string
	ADIdentifiers   []string
	Hosts           map[string]string
	Ports           []listeners.ContainerPort
	Pid             int
	Hostname        string
	filterTemplates func(map[string]integration.Config)
}

// Equal returns whether the two dummyService are equal
func (s *dummyService) Equal(o listeners.Service) bool {
	return reflect.DeepEqual(s, o)
}

// GetServiceID returns the service entity name
func (s *dummyService) GetServiceID() string {
	return s.ID
}

// GetADIdentifiers returns dummy identifiers
func (s *dummyService) GetADIdentifiers() []string {
	return s.ADIdentifiers
}

// GetHosts returns dummy hosts
func (s *dummyService) GetHosts() (map[string]string, error) {
	return s.Hosts, nil
}

// GetPorts returns dummy ports
func (s *dummyService) GetPorts() ([]listeners.ContainerPort, error) {
	return s.Ports, nil
}

// GetTags returns the tags for this service
func (s *dummyService) GetTags() ([]string, error) {
	return nil, nil
}

// GetTagsWithCardinality returns the tags for this service
func (s *dummyService) GetTagsWithCardinality(_ string) ([]string, error) {
	return s.GetTags()
}

// GetPid return a dummy pid
func (s *dummyService) GetPid() (int, error) {
	return s.Pid, nil
}

// GetHostname return a dummy hostname
func (s *dummyService) GetHostname() (string, error) {
	return s.Hostname, nil
}

// IsReady returns if the service is ready
func (s *dummyService) IsReady() bool {
	return true
}

// HasFilter returns false
func (s *dummyService) HasFilter(_ filter.Scope) bool {
	return false
}

// GetExtraConfig isn't supported
func (s *dummyService) GetExtraConfig(_ string) (string, error) {
	return "", nil
}

// FilterTemplates calls filterTemplates, if not nil
func (s *dummyService) FilterTemplates(configs map[string]integration.Config) {
	if s.filterTemplates != nil {
		(s.filterTemplates)(configs)
	}
}
