package signals

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type SignalHTTP struct {
	URL        string
	Assertion  AssertionType
	Value      string
	Within     time.Duration
	httpClient http.Client
}

func NewHTTP(url string, assert AssertionType, value string, within time.Duration) *SignalHTTP {
	httpClient := http.Client{
		Timeout: within,
	}

	return &SignalHTTP{
		URL:        url,
		Assertion:  assert,
		Value:      value,
		Within:     within,
		httpClient: httpClient,
	}
}

func (s *SignalHTTP) Assert() (bool, error) {
	slog.Info("Asserting HTTP signal", "url", s.URL, "assertion", s.Assertion, "value", s.Value, "within", s.Within)

	// Make an HTTP request to the URL and check the response against the assertion
	resp, err := s.httpClient.Get(s.URL)

	if err != nil {
		// If its a timeout error, that is a valid result and not an actual error
		if isTimeout := os.IsTimeout(err); isTimeout {
			if s.Assertion == AssertHTTPTimeout {
				slog.Info("HTTP request timed out as expected", "url", s.URL, "assertion", s.Assertion, "value", s.Value, "within", s.Within)
				return true, nil
			}

			slog.Info("HTTP request timed out", "url", s.URL, "assertion", s.Assertion, "value", s.Value, "within", s.Within)
			return false, nil
		}

		return false, fmt.Errorf("Failed to make HTTP request to %s: %v", s.URL, err)
	}

	defer resp.Body.Close()

	switch s.Assertion {
	case AssertHTTPStatus200:
		return resp.StatusCode == http.StatusOK, nil
	case AssertHTTPStatusNot200:
		return resp.StatusCode != http.StatusOK, nil
	case AssertHTTPBodyContains:
		// Read the response body and check if it contains the specified value
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		bodyString := string(bodyBytes)
		return strings.Contains(bodyString, s.Value), nil
	case AssertHTTPTimeout:
		// If we got a response, then the assertion fails since we expected a timeout
		return false, nil
	default:
		return false, fmt.Errorf("Unknown assertion type: %v", s.Assertion)
	}
}

func (s *SignalHTTP) Name() string {
	return fmt.Sprintf("http/%s/%s/%s", s.Assertion, s.Value, s.URL)
}

func (s *SignalHTTP) Timeout() time.Duration {
	return s.Within
}

func (s *SignalHTTP) Close() error {
	// No resources to close for the HTTP signal
	return nil
}

func (s *SignalHTTP) Type() SignalType {
	return SignalTypeHTTP
}
