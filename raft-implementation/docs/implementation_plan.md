# Implementation Plans for Raft KV Store Issues

---

## Raft Node Bootstrapping

### #1 Basic setup with raft nodes
**Steps:**
[DONE] 1. Define Raft node struct (ID, address, role, term, log, votedFor, neighbors).
[DONE] 2. Set up configuration for cluster (static node list or config file).
[DONE] 3. Initialize nodes with default state (follower, term = 0, empty log).
[DONE] 4. Create networking setup (gRPC server for RPCs).
5. Start node main loop: listen for RPCs, handle timers.

### #2 Follower state
**Steps:**
1. Implement follower logic:  
   - Receive RPCs (AppendEntries, RequestVote).
   - Reset election timeout on valid leader communication.
2. Add state transitions:  
   - Follower → candidate on timeout.
   - Follower stays follower on valid leader heartbeat.

### #3 Election timeout
**Steps:**
1. Implement randomized election timeout (e.g., 150–300ms).
2. On timeout, transition to candidate and start election.
3. Use goroutines/timers for timeout handling.

---

## Leader Election

### #4 Candidate/leader transitions
**Steps:**
1. Candidate state:  
   - Send RequestVote RPCs.
   - Vote for self, reset election timeout.
2. Count votes, become leader if majority reached.
3. Transition back to follower if AppendEntries received from valid leader.

### #5 RequestVote RPC
**Steps:**
1. Define RequestVote RPC in proto file.
2. Implement RequestVote handler:  
   - Compare terms, log up-to-date check.
   - Grant vote if not voted and candidate log is at least as up-to-date.
3. Implement client for sending/voting.

### #6 Vote counting, leader selection
**Steps:**
1. Track votes received per election term.
2. Become leader if votes > N/2.
3. Handle split vote: re-election after timeout.

---

## Log Replication

### #7 Log struct
**Steps:**
1. Define log entry struct (index, term, command, committed flag).
2. Store log in node (slice/array).
3. Implement methods: append entry, get entry, truncate log.

### #8 AppendEntries RPC
**Steps:**
1. Define AppendEntries RPC (proto: leader term, prevLogIndex/Term, entries[], leaderCommit).
2. Implement handler:  
   - Check log consistency (compare prevLogIndex/Term).
   - If matching, append new entries.
   - Reply success/failure.
3. Leader: send heartbeats (empty entries) periodically.

### #9 Commit logic
**Steps:**
1. Track commitIndex for each node.
2. Leader advances commitIndex when entry is stored on majority.
3. Apply committed entries to state machine.

---

## KV Store Logic

### #10 State machine: SET/GET/DELETE
**Steps:**
1. Define state machine (map[string]string).
2. Implement Apply(command) to mutate state.
3. Handle SET (add/update), DELETE (remove), GET (read-only).

### #11 Apply committed entries
**Steps:**
1. On commitIndex advance, apply all un-applied entries ≤ commitIndex.
2. Ensure deterministic application order.
3. Update state machine, respond to client requests.

---

## Client API

### #12 gRPC endpoints
**Steps:**
1. Define client API (proto: SET/GET/DELETE request/response).
2. Implement gRPC server on leader node for client requests.
3. Route all SET/DELETE to leader; followers redirect or proxy.
4. Handle GET on any node (if eventual consistency is acceptable).
5. Implement client CLI or test script for demo.

---

## General Tips & Best Practices

- Use Go channels/goroutines for concurrency.
- Log every state transition and RPC for debugging.
- Keep configuration and error handling robust.
- Write tests for each module before moving to next.
- Document API endpoints and message formats.

---

## Example Architecture Diagram (Recommended)

- Draw: nodes, client, leader election, log replication, state machine, API calls.

---

## References

- [Raft Paper](https://raft.github.io/raft.pdf)
- [gRPC Go Docs](https://grpc.io/docs/languages/go/)
- [Raft Visual Demo](https://thesecretlivesofdata.com/raft/)

---

*Use this plan to tackle each issue methodically. Update your GitHub Issues with checklists and link relevant code/docs as you progress!*