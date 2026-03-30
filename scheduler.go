package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Queens-Rocket-Engineering-Team/ctl-havoc/faults"
	"github.com/Queens-Rocket-Engineering-Team/ctl-havoc/signals"
)

type Scenario struct {
	Name          string
	Description   string
	Fault         faults.Fault
	DuringSignals []signals.Signal
	AfterSignals  []signals.Signal
}

func (s *Scenario) Run(ctx context.Context) error {
	println("Running scenario:", s.Name)
	println("Description:", s.Description)

	// Inject the fault
	if err := s.Fault.Inject(ctx); err != nil {
		return fmt.Errorf("Failed to inject fault: %v", err)
	}

	if s.Fault.Duration() > 0 {
		faultCtx, cancel := context.WithTimeout(ctx, s.Fault.Duration())
		defer cancel()

		// Assert during signals
		for _, sig := range s.DuringSignals {
			go func(sig signals.Signal) {
				// We want to assert the signal within the context of the fault duration
				// so we create a child context with the appropriate timeout for each signal assertion
				sigCtx, cancelSig := context.WithTimeout(faultCtx, sig.Timeout())
				defer cancelSig()

				ok, err := sig.Assert()
				if err != nil {
					slog.Error("Error asserting signal during fault injection", "signal", sig.Name(), "error", err)
				} else {
					slog.Info("During signal assertion result", "signal", sig.Name(), "result", ok)
				}

				<-sigCtx.Done() // Wait for the signal assertion to complete or timeout

				if sigCtx.Err() == context.DeadlineExceeded {
					slog.Warn("Signal assertion timed out during fault injection", "signal", sig.Name())
				}
			}(sig)
		}

		<-faultCtx.Done() // Wait for the fault duration to elapse

		// Cleanup the fault
		if err := s.Fault.Cleanup(ctx); err != nil {
			return fmt.Errorf("Failed to cleanup fault: %v", err)
		}
	}

	// Assert after signals
	for _, sig := range s.AfterSignals {
		ok, err := sig.Assert()
		if err != nil {
			return fmt.Errorf("Error asserting signal after fault cleanup: %v", err)
		} else {
			slog.Info("After signal assertion result", "signal", sig.Name(), "result", ok)
		}
	}

	slog.Info("Completed scenario", "name", s.Name)
	return nil
}

type Scheduler struct {
	Scenarios []Scenario
	Mode      string // "sequential" or "random"
}

func (s *Scheduler) Run(ctx context.Context) {
	// Currently only supports sequential mode
	for _, scenario := range s.Scenarios {
		scenario.Run(ctx)
	}
}
