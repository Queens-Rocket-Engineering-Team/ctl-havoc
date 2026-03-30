package faults

import (
	"context"
	"time"
)

type FaultType string

const (
	FaultTypeDockerKill           FaultType = "docker-kill"            // Kill a process with SIGKILL
	FaultTypeDockerPause          FaultType = "docker-pause"           // Pause a Docker container
	FaultTypeDockerComposeRestart FaultType = "docker-compose-restart" // Restart the entire server stack
	FaultTypeNetworkDelay         FaultType = "network-delay"          // Emulate network latency with tc netem
	FaultTypeNetworkLoss          FaultType = "network-loss"           // Emulate packet loss with tc netem
	FaultTypePacketFlood          FaultType = "packet-flood"           // Flood the UDP stream with packets
	FaultTypePacketBlock          FaultType = "packet-block"           // Block the UDP stream with iptables
)

type Fault interface {
	Inject(ctx context.Context) error  // Inject the fault into the system
	Cleanup(ctx context.Context) error // Clean up the fault and restore the system to normal
	Name() string                      // Get a human-readable name for the fault
	Duration() time.Duration           // Get the duration for which the fault should be active
	Close() error                      // Close any resources used by the fault (e.g. Docker client connections)
}
