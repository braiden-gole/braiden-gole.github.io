// -----------------------------------------------------------------------------
// Package: main
// File: main.go
//
// Description:
//
//	This file serves as an executable demonstration driver for the consensus
//	package. It instantiates a follower node, populates initial log entries,
//	simulates an incoming InstallSnapshot RPC from a leader, triggers log
//	compaction, and verifies state machine synchronization.
//
// -----------------------------------------------------------------------------
package main

import (
	"fmt"

	raft "github.com/braiden-gole/reliable-replicated-redundant-and-fault-tolerant/reliablereplicatedredundantandfaulttolerant"
)

func main() {
	// Initialize a new follower node instance with identifier 2.
	followerNode, err := raft.NewRaftNode(2)

	// Handle any errors encountered during node creation.
	if err != nil {
		fmt.Printf("Failed to create node: %v\n", err)

		return
	}

	// Populate the node's initial log entries prior to compaction.
	followerNode.LogEntries = []raft.LogEntry{
		{Command: raft.StateMachineCommand{Key: "userAlice", Value: "Alice"}, ReceivedTerm: 1, LogIndex: 1},
		{Command: raft.StateMachineCommand{Key: "userBob", Value: "Bob"}, ReceivedTerm: 1, LogIndex: 2},
	}

	// Display the initial log length before the snapshot is processed.
	fmt.Printf("Initial log length before snapshot: %d\n", len(followerNode.LogEntries))

	// Construct the incoming snapshot arguments payload simulating a leader transmission.
	snapshotArguments := raft.InstallSnapshotArguments{
		Term:              2,
		LeadIdentifier:    1,
		LastIncludedIndex: 2,
		LastIncludedTerm:  1,
		SnapshotData: map[string]string{
			"userAlice":   "alice@gmail.com",
			"userBob":     "bob@gmail.com",
			"userCharlie": "charlie@gmail.com",
		},
	}

	// Initialize the reply payload structure for the RPC response.
	var replyPayload raft.InstallSnapshotReply

	// Execute the snapshot installation handler on the follower node.
	executionError := followerNode.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Terminate execution if the snapshot installation fails.
	if executionError != nil {
		fmt.Printf("Error installing snapshot: %v\n", executionError)

		return
	}

	// Print the updated commit index and the truncated log length.
	fmt.Printf("Node current commit index: %d\n", followerNode.RetrieveCommitIndex())
	fmt.Printf("Log length after log compaction/truncation: %d\n", len(followerNode.LogEntries))

	// Retrieve updated key-value data from the synchronized state machine.
	alice, existsAlice := followerNode.RetrieveStateMachineData("userAlice")
	charlie, existsCharlie := followerNode.RetrieveStateMachineData("userCharlie")

	// Print the state machine value for userAlice if present.
	if existsAlice {
		fmt.Printf("State machine value for userAlice: %s\n", alice)
	}

	// Print the state machine value for userCharlie if present.
	if existsCharlie {
		fmt.Printf("State machine value for userCharlie: %s\n", charlie)
	}
}
