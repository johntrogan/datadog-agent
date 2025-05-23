// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build windows

package modules

import (
	"github.com/DataDog/datadog-agent/pkg/eventmonitor"
	netconfig "github.com/DataDog/datadog-agent/pkg/network/config"
	"github.com/DataDog/datadog-agent/pkg/system-probe/api/module"
	"github.com/DataDog/datadog-agent/pkg/system-probe/config"
)

func init() { registerModule(EventMonitor) }

// EventMonitor - Event monitor Factory
var EventMonitor = &module.Factory{
	Name:             config.EventMonitorModule,
	ConfigNamespaces: eventMonitorModuleConfigNamespaces,
	Fn:               createEventMonitorModule,
}

func createProcessMonitorConsumer(_ *eventmonitor.EventMonitor, _ *netconfig.Config) error {
	return nil
}
