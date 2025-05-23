// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package collector

import (
	"maps"
	"sync"

	checkid "github.com/DataDog/datadog-agent/pkg/collector/check/id"
)

// collectorErrors holds the error objects
type collectorErrors struct {
	loader map[string]map[string]string // check Name -> loader -> error
	run    map[checkid.ID]string        // check ID -> error
	m      sync.RWMutex
}

// newCollectorErrors returns an instance holding autoconfig errors stats
func newCollectorErrors() *collectorErrors {
	return &collectorErrors{
		loader: make(map[string]map[string]string),
		run:    make(map[checkid.ID]string),
	}
}

// setLoaderError will safely set the error for that check and loader to the LoaderErrorStats
func (ce *collectorErrors) setLoaderError(checkName string, loaderName string, err string) {
	_, found := ce.loader[checkName]
	if !found {
		ce.loader[checkName] = make(map[string]string)
	}

	ce.loader[checkName][loaderName] = err
}

// removeLoaderErrors removes the errors for a check (usually when successfully loaded)
func (ce *collectorErrors) removeLoaderErrors(checkName string) {
	delete(ce.loader, checkName)
}

// GetLoaderErrors will safely get the errors regarding loaders
func (ce *collectorErrors) getLoaderErrors() map[string]map[string]string {
	ce.m.RLock()
	defer ce.m.RUnlock()

	errorsCopy := make(map[string]map[string]string)

	for check, loaderErrors := range ce.loader {
		errorsCopy[check] = make(map[string]string)
		maps.Copy(errorsCopy[check], loaderErrors)
	}

	return errorsCopy
}

func (ce *collectorErrors) setRunError(checkID checkid.ID, err string) {
	ce.m.Lock()
	defer ce.m.Unlock()

	ce.run[checkID] = err
}

func (ce *collectorErrors) getRunErrors() map[checkid.ID]string {
	ce.m.RLock()
	defer ce.m.RUnlock()

	runCopy := make(map[checkid.ID]string)
	maps.Copy(runCopy, ce.run)

	return runCopy
}
