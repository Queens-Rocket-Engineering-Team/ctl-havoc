package main

import (
	"fmt"
	"time"

	"github.com/Queens-Rocket-Engineering-Team/ctl-havoc/faults"
	"github.com/Queens-Rocket-Engineering-Team/ctl-havoc/signals"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ComposeProject string         `yaml:"compose_project"`
	Scenarios      []YAMLScenario `yaml:"scenarios"`
}

type YAMLScenario struct {
	Name          string       `yaml:"name"`
	Description   string       `yaml:"description"`
	Fault         YAMLFault    `yaml:"fault"`
	DuringSignals []YAMLSignal `yaml:"during_signals"`
	AfterSignals  []YAMLSignal `yaml:"after_signals"`
}

type YAMLFault struct {
	Type     faults.FaultType `yaml:"type"`
	Duration string           `yaml:"duration,omitempty"` // Optional for instantaneous faults

	// Fault-specific parameters
	Service string `yaml:"service"` // For Docker faults, the service name to target
}

type YAMLSignal struct {
	Type      signals.SignalType    `yaml:"type"`
	Assertion signals.AssertionType `yaml:"assertion"`
	Within    string                `yaml:"within"`

	// Signal-specific parameters
	URL string `yaml:"url,omitempty"` // For HTTP signals
}

func (yf *YAMLFault) ToFault(composeProject string) (faults.Fault, error) {
	duration, err := time.ParseDuration(yf.Duration)
	if err != nil {
		return nil, fmt.Errorf("invalid fault duration: %v", err)
	}

	switch yf.Type {
	case faults.FaultTypeDockerPause:
		// For simplicity, using hardcoded project and service names here
		return faults.NewDockerPause(composeProject, yf.Service, duration)
	case faults.FaultTypeDockerKill:
		return faults.NewDockerKill(composeProject, yf.Service)
	default:
		return nil, fmt.Errorf("unsupported fault type: %s", yf.Type)
	}
}

func (ys *YAMLSignal) ToSignal() (signals.Signal, error) {
	switch ys.Type {
	case signals.SignalTypeHTTP:
		timeout, err := time.ParseDuration(ys.Within)
		if err != nil {
			return nil, err
		}

		return signals.NewHTTP(ys.URL, ys.Assertion, "", timeout), nil
	default:
		return nil, fmt.Errorf("unsupported signal type: %s", ys.Type)
	}
}

func parseConfig(yamlData []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %v", err)
	}
	return &config, nil
}

func unmarshalYAMLToScenario(yamlScenario YAMLScenario, composeProject string) (*Scenario, error) {
	fault, err := yamlScenario.Fault.ToFault(composeProject)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML fault to Fault interface: %v", err)
	}

	var duringSignals []signals.Signal
	for _, ys := range yamlScenario.DuringSignals {
		sig, err := ys.ToSignal()
		if err != nil {
			return nil, fmt.Errorf("failed to convert YAML signal to Signal interface: %v", err)
		}
		duringSignals = append(duringSignals, sig)
	}

	var afterSignals []signals.Signal
	for _, ys := range yamlScenario.AfterSignals {
		sig, err := ys.ToSignal()
		if err != nil {
			return nil, fmt.Errorf("failed to convert YAML signal to Signal interface: %v", err)
		}
		afterSignals = append(afterSignals, sig)
	}

	return &Scenario{
		Name:          yamlScenario.Name,
		Description:   yamlScenario.Description,
		Fault:         fault,
		DuringSignals: duringSignals,
		AfterSignals:  afterSignals,
	}, nil
}
