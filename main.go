package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	// Read and parse the YAML configuration file
	// scenarios.yaml for now in CWD

	// Read the YAML file
	yamlData, err := os.ReadFile("scenarios.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to read YAML file: %v", err))
	}

	config, err := parseConfig(yamlData)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse YAML config: %v", err))
	}

	// For simplicity, we'll just run the first scenario in the config
	if len(config.Scenarios) == 0 {
		panic("No scenarios defined in the config")
	}

	scenario, err := unmarshalYAMLToScenario(config.Scenarios[0], config.ComposeProject)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert YAML scenario to Scenario struct: %v", err))
	}

	scheduler := Scheduler{
		Scenarios: []Scenario{*scenario},
		Mode:      "sequential",
	}

	scheduler.Run(context.Background())
}
