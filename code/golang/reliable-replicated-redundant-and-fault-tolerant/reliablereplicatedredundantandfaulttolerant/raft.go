// -----------------------------------------------------------------------------
// Package: reliablereplicatedredundantandfaulttolerant
// File: raft.go
//
// Description:
//
//	This file implements core components of the Raft consensus protocol in Go,
//	focusing on node initialization, thread-safe state access, and the robust
//	InstallSnapshot RPC handler complete with absolute log compaction.
//
//	The implementation prioritizes strict thread-safety via sync.Mutex,
//	defensive programming against malformed inputs, and clean encapsulation
//	by deep-copying external state machine and snapshot payloads.
//
// -----------------------------------------------------------------------------
package reliablereplicatedredundantandfaulttolerant

import (
	"errors"
	"sync"
)

// StateMachineCommand represents an individual key-value command destined
// to be applied to the replicated state machine.
type StateMachineCommand struct {
	Key   string
	Value string
}

// LogEntry encapsulates a single entry in the Raft replicated log,
// pairing the state machine command with the term in which it was received
// and its absolute global log index.
type LogEntry struct {
	Command      StateMachineCommand
	ReceivedTerm int
	LogIndex     int
}

// Snapshot captures a compacted point-in-time checkpoint of the Raft node's
// state, allowing obsolete log prefixes to be safely discarded to save memory.
type Snapshot struct {
	LastIncludedIndex int
	LastIncludedTerm  int
	StateData         map[string]string
}

// Node represents an individual participant in a Raft consensus cluster.
// All internal mutable fields are protected by a mutual exclusion lock to
// guarantee thread safety across concurrent RPC handlers and queries.
type Node struct {
	mutex            sync.Mutex
	nodeIdentifier   int
	currentTerm      int
	votedFor         int
	LogEntries       []LogEntry
	commitIndex      int
	lastApplied      int
	LatestSnapshot   Snapshot
	stateMachineData map[string]string
}

// NewRaftNode initializes and returns a new Raft consensus node instance.
// It validates that the provided node identifier is a positive integer and
// sets up initial default states, empty log slices, and active snapshot maps.
func NewRaftNode(nodeIdentifier int) (*Node, error) {
	// Node identifiers must be positive integers in a valid cluster configuration.
	if nodeIdentifier <= 0 {
		return nil, errors.New("node identifier must be a positive integer.")
	}

	return &Node{
		nodeIdentifier: nodeIdentifier,
		currentTerm:    0,
		votedFor:       -1,
		LogEntries:     make([]LogEntry, 0),
		commitIndex:    0,
		lastApplied:    0,
		LatestSnapshot: Snapshot{
			LastIncludedIndex: 0,
			LastIncludedTerm:  0,
			StateData:         make(map[string]string),
		},
		stateMachineData: make(map[string]string),
	}, nil
}

// InstallSnapshotArguments defines the RPC payload transmitted by a leader
// to a follower when instructing it to install a snapshot.
type InstallSnapshotArguments struct {
	Term              int
	LeadIdentifier    int
	LastIncludedIndex int
	LastIncludedTerm  int
	SnapshotData      map[string]string
}

// InstallSnapshotReply defines the response payload returned by a follower
// after processing an incoming InstallSnapshot RPC request.
type InstallSnapshotReply struct {
	Term int
}

// RetrieveCommitIndex provides thread-safe read access to the node's
// current commit index marker.
func (node *Node) RetrieveCommitIndex() int {
	node.mutex.Lock()

	defer node.mutex.Unlock()

	return node.commitIndex
}

// RetrieveStateMachineData provides thread-safe lookup for values stored
// within the node's replicated state machine.
func (node *Node) RetrieveStateMachineData(key string) (string, bool) {
	node.mutex.Lock()

	defer node.mutex.Unlock()

	value, exists := node.stateMachineData[key]

	return value, exists
}

// HandleInstallSnapshot processes an incoming InstallSnapshot RPC from a leader.
// It validates terms, checks bounds, discards redundant snapshots, compacts the
// log using absolute indices, and synchronizes the state machine via deep-copying.
func (node *Node) HandleInstallSnapshot(arguments InstallSnapshotArguments, reply *InstallSnapshotReply) error {
	// Guard against nil pointer dereference on the reply payload.
	if reply == nil {
		return errors.New("reply pointer cannot be nil")
	}

	node.mutex.Lock()

	defer node.mutex.Unlock()

	// Perform defensive bounds checks on incoming arguments.
	if arguments.Term < 0 {
		return errors.New("incoming snapshot term cannot be negative")
	}

	if arguments.LastIncludedIndex < 0 {
		return errors.New("incoming snapshot last included index cannot be negative")
	}

	if arguments.SnapshotData == nil {
		return errors.New("snapshot data map cannot be nil")
	}

	// Reject snapshots from outdated leader terms.
	if arguments.Term < node.currentTerm {
		reply.Term = node.currentTerm

		return errors.New("incoming snapshot term is older than current node term")
	}

	// Step down or update term if the incoming term is higher than local term.
	if arguments.Term > node.currentTerm {
		node.currentTerm = arguments.Term
		node.votedFor = -1
	}

	reply.Term = node.currentTerm

	// Ignore snapshots that are older than or equal to the node's existing snapshot.
	if arguments.LastIncludedIndex <= node.LatestSnapshot.LastIncludedIndex {
		return nil
	}

	var remainingEntries []LogEntry

	// Discard log entries superseded by the snapshot using their absolute LogIndex.
	for _, entry := range node.LogEntries {
		if entry.LogIndex > arguments.LastIncludedIndex {
			remainingEntries = append(remainingEntries, entry)
		}
	}

	node.LogEntries = remainingEntries

	// Deep copy snapshot state data to maintain strict encapsulation
	// and prevent external memory mutation side effects.
	clonedStateData := make(map[string]string, len(arguments.SnapshotData))

	for key, value := range arguments.SnapshotData {
		clonedStateData[key] = value
	}

	// Update the node's latest snapshot reference and state machine data.
	node.LatestSnapshot = Snapshot{
		LastIncludedIndex: arguments.LastIncludedIndex,
		LastIncludedTerm:  arguments.LastIncludedTerm,
		StateData:         clonedStateData,
	}

	node.stateMachineData = clonedStateData

	// Advance commit index to match the snapshot boundary if necessary.
	if node.commitIndex < arguments.LastIncludedIndex {
		node.commitIndex = arguments.LastIncludedIndex
	}

	// Advance last applied index to match the snapshot boundary if necessary.
	if node.lastApplied < arguments.LastIncludedIndex {
		node.lastApplied = arguments.LastIncludedIndex
	}

	return nil
}
