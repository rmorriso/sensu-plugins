package main

import "testing"

func TestInit(t *testing.T) {
	config, err := Init("pagerduty.yaml")
	if err != nil {
		t.Errorf("Failed opening pagerduty.yaml: %s\n", err)
	}

	if config.ServiceKey != "REPLACE-WITH-YOUR-PAGERDUTY-SERVICE-KEY" {
		t.Error("Failed to set config.ServiceKey from YAML")
	}
}
