package signals

import "time"

type SignalType string

const (
	SignalTypeHTTP SignalType = "http"
)

type AssertionType string

const (
	AssertHTTPStatus200    AssertionType = "http-status-200"
	AssertHTTPStatusNot200 AssertionType = "http-status-not-200"
	AssertHTTPBodyContains AssertionType = "http-body-contains"
	AssertHTTPTimeout      AssertionType = "http-timeout"
)

type Signal interface {
	Assert() (bool, error)  // Assert that the signal condition is met
	Type() SignalType       // Get the type of the signal
	Name() string           // Get a human-readable name for the signal
	Timeout() time.Duration // Get the duration for which the signal assertion should be valid
	Close() error           // Close any resources used by the signal (e.g. HTTP client connections)
}
