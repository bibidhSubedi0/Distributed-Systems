# Distributed Key-Value Store Using Raft (Go)

## Project Overview

-> distributed key-value (KV) store that uses the Raft consensus algorithm to ensure fault-tolerant, strongly consistent data replication across multiple nodes. 

-> The KV store exposes a simple API (SET/GET/DELETE) and is resilient to node failures by electing a leader and replicating logs in a cluster.

---

## Goals

- Implement a highly available, strongly consistent key-value store.
- Learn and apply distributed systems concepts: consensus, replication, fault tolerance.
---

## Tools & Technologies

| Tool/Tech            | Purpose                                      |
|----------------------|----------------------------------------------|
| Go                   | Main programming language                    |
| gRPC                 | Node-to-node and client-to-cluster RPC       |
| Protocol Buffers     | Message and API schema definition            |
| GitHub Projects      | Task and milestone tracking                  |
| Markdown/README      | Documentation                                |
| draw.io/Excalidraw   | Architecture and protocol diagrams           |
| Go Testing Framework | Unit and integration tests                   |
| Git                  | Version control                              |
|            Docker    | Easy multi-node deployment                   |

---

## Core Concepts Required

- **Distributed Systems Basics**
  - CAP theorem (Consistency, Availability, Partition Tolerance)
  - Replication (leader/follower model)
- **Consensus Algorithms**
  - Raft roles: follower, candidate, leader
  - Leader election, log replication, commit mechanics
- **State Machines**
  - Deterministic application of commands (SET/GET/DELETE)
- **RPC and Networking**
  - gRPC concepts, protobuf schema, network communication
- **Concurrency in Go**
  - Goroutines, channels, mutexes for node logic
- **Persistence**
  - Write-ahead log, state snapshots (optional for MVP)
- **Testing Distributed Systems**
  - Simulating node failures, partitions, and recovery

---

## Project Structure

```
raft-kv/
├── main.go
├── raft/
│   ├── node.go            # Raft node logic (roles, election, RPC handlers)
│   ├── log.go             # Log entry management, commit logic
│   ├── state_machine.go   # KV store state machine (apply commands)
├── proto/
│   ├── raft.proto         # Raft RPC message definitions
│   ├── kv.proto           # Client API message definitions
├── client/
│   └── client.go          # CLI or API client for testing cluster
├── docs/
│   ├── PROJECT_OVERVIEW.md
│   ├── Progress.md        # Task and milestone checklist
│   └── diagrams/          # Architecture diagrams
├── tests/
│   ├── raft_test.go
│   ├── kv_test.go
└── README.md
```

---

## Implementation Steps

### 1. Project Initialization

- Set up Go module: `go mod init github.com/yourusername/raft-kv`
- Set up basic repo structure and initial README.

### 2. Define Data Structures & Protocols

- Design Raft log entry, node state, and state machine structs in Go.
- Write `.proto` files for Raft RPCs (RequestVote, AppendEntries) and client API (SET/GET/DELETE).

### 3. Implement Raft Node

- Follower logic: receive RPCs, apply committed entries.
- Candidate logic: election timeout, vote request/response.
- Leader logic: handle client commands, replicate logs, send heartbeats.

### 4. Implement Leader Election

- Election timeout handling (randomized).
- RequestVote RPC: voting logic, persistent term, last log index/term.
- Leader selection and role transition.

### 5. Implement Log Replication

- AppendEntries RPC: leader sends log entries to followers.
- Log consistency: truncate/replace follower logs as needed.
- Commit entries when replicated to majority; apply to state machine.

### 6. State Machine and KV Store

- Apply committed commands to in-memory map (SET/DELETE).
- Ensure deterministic command application order.

### 7. Persistent Log Storage (Optional for MVP)

- Write log entries to disk before commit.
- Implement periodic snapshots for state machine (for large logs).

### 8. Client API

- Expose gRPC endpoints for SET/GET/DELETE.
- Implement CLI client or simple web API for user interaction.
- Route all writes to leader; followers redirect or proxy requests.

### 9. Testing and Simulation

- Unit tests for Raft node, log logic, state machine.
- Integration tests: start multiple nodes, simulate leader election, log replication, failover.
- Simulate network partitions and node failures.

### 10. Documentation & Visualization

- Update README with usage, architecture, and design decisions.
- Document API endpoints and message formats.
- Create diagrams showing node roles, data flow, election process.

### 11. Advanced Features (Optional)

- Dynamic cluster membership: add/remove nodes.
- Log compaction and snapshotting.
- Monitoring, metrics, and dashboards.

---

## Sample Task List (docs/Progress.md)

````markdown name=docs/Progress.md
# Project Progress

## Milestone 1: Raft Node Bootstrapping
- [x] Initialize Raft node struct
- [x] Follower state
- [ ] Election timeout
- [ ] Candidate/leader transitions

## Milestone 2: Leader Election
- [ ] RequestVote RPC
- [ ] Vote counting, leader selection

## Milestone 3: Log Replication
- [ ] Log struct
- [ ] AppendEntries RPC
- [ ] Commit logic

## Milestone 4: KV Store Logic
- [ ] State machine: SET/GET/DELETE
- [ ] Apply committed entries

## Milestone 5: Client API
- [ ] gRPC endpoints
- [ ] CLI/web client

## Milestone 6: Persistence & Testing
- [ ] Disk log
- [ ] Snapshotting
- [ ] Unit/integration tests