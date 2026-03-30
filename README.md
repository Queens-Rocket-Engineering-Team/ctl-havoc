# HAVOC
**High Availability Validation Orchestration and Control**

HAVOC is a chaos engineering service for the QRET launch control stack. It injects controlled faults into the Docker Compose services running on the [coordination server](https://github.com/Queens-Rocket-Engineering-Team/prop-teststand/) and verifies that the stack recovers correctly — before a hotfire, not during one.

# Why This Exists
The launch control stack is a distributed system. The control server, redis, mediamtx, network, all have failure boundaries between them. Any one of them can fail: a container crashes, the network degrades, redis runs out of memory. We want to know how the system fails under these conditions and if it can recover cleanly.

Manual testing doesn't cover this well. You can manually kill a container and look at the logs, but you can't do that systematically across every part of the system for every failure mode before a hotfire, and even then it is difficult to show that the recovery was correct.

HAVOC automates this. It defines faults as structured scenarios with explicit hypotheses: "if mediamtx is killed, the control server stays healthy and streams resume within 30 seconds", and verifies them against real observable signals from the running stack.

# Architecture

HAVOC runs as its own process on the control server, separate from the Docker Compose stack. It is composed of:
- **Faults** - are "injected" into the stack, simulating errors that we want to test the stack against (e.g. a docker container crashes)

- **Signals** - are "asserted" against the stack, checking that the stack meets well-defined conditions (e.g. GET /health returns a 200 status code)

- **Scheduler** - selects which scenarios are run based on its configured mode (sequential, random, or targeted), injects faults and asserts signals. Waits for the stack to be healthy between scenarios before starting up the next one.


# Scenarios Spec Format
Scenarios are defined in YAML. HAVOC loads them at startup and the scheduler selects from them at runtime.
```yaml
compose_project: "prop-teststand" # The Docker Compose project name 

scenarios:
  - name: "Server Pause"
    description: "Simulate a server pause by Docker pause command"
    fault:
      type: docker-pause # The type of fault to inject into the stack
      service: "server" # Example of a fault-specific parameter, docker container related faults need a service/container to run against
      duration: 10s # How long the fault will last before timing out and cleaning up
    during_signals: # These signals run during the fault (for non-instantaneous faults) to check behaviour under degraded conditions
      - type: http # This is the type of signal to assert against
        url: "http://localhost:8000/health" # Example of a signal-specific parameter, http needs a url to request
        assertion: http-timeout # This is the condition we expect as the fault is active. In this case, we expect the server to not respond when paused
        within: 2s # Timeout for the signal
    after_signals:
      - type: http 
        url: "http://localhost:8000/health"
        assertion: http-status-200 # We expect the service to recover and return 200 OK after the pause
        within: 2s # We expect the 200 status code within 2 seconds of the fault ending
```

# Fault Types
TODO

# Signal Types
TODO