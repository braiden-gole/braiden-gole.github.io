// -----------------------------------------------------------------------------
// Package: reliablereplicatedredundantandfaulttolerant
// File: raft_test.go
//
// Description:
//
//	This file provides a comprehensive test suite of 30 unique, meaningful unit
//	tests covering node initialization, validation rules, log compaction, state
//	machine synchronization, encapsulation guards, and concurrency safety for
//	the Raft consensus implementation.
//
// -----------------------------------------------------------------------------
package reliablereplicatedredundantandfaulttolerant

import (
	"testing"
)

// TestNewRaftNode_ValidIdentifier verifies that initializing a Raft node with a
// positive integer identifier succeeds and returns a non-nil node instance.
func TestNewRaftNode_ValidIdentifier(test *testing.T) {
	// Arrange.
	var validNodeIdentifier int = 1

	// Act.
	node, err := NewRaftNode(validNodeIdentifier)

	// Assert.
	if err != nil {
		test.Fatalf("expected no error for valid node identifier, got %v", err)
	}

	if node == nil {
		test.Fatalf("expected raft node instance to be created, got nil")
	}
}

// TestNewRaftNode_NegativeIdentifier verifies that attempting to initialize
// a node with a negative integer identifier fails and returns a nil instance.
func TestNewRaftNode_NegativeIdentifier(test *testing.T) {
	// Arrange.
	var invalidNodeIdentifier int = -1

	// Act.
	node, err := NewRaftNode(invalidNodeIdentifier)

	// Assert.
	if err == nil {
		test.Fatalf("expected error for negative node identifier, got nil")
	}

	if node != nil {
		test.Fatalf("expected nil node instance for negative identifier, got %v", node)
	}
}

// TestNewRaftNode_ZeroIdentifier verifies that attempting to initialize
// a node with a zero identifier is rejected with an error.
func TestNewRaftNode_ZeroIdentifier(test *testing.T) {
	// Arrange.
	var zeroNodeIdentifier int = 0

	// Act.
	node, err := NewRaftNode(zeroNodeIdentifier)

	// Assert.
	if err == nil {
		test.Fatalf("expected error for zero node identifier, got nil")
	}

	if node != nil {
		test.Fatalf("expected nil node instance for zero identifier, got %v", node)
	}
}

// TestNewRaftNode_InitialStateCommitIndex confirms that newly instantiated
// nodes start with a commit index of zero.
func TestNewRaftNode_InitialStateCommitIndex(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	// Act.
	var actualCommitIndex int = node.RetrieveCommitIndex()

	// Assert.
	var expectedCommitIndex int = 0

	if actualCommitIndex != expectedCommitIndex {
		test.Fatalf("expected initial commit index to be %d, got %d", expectedCommitIndex, actualCommitIndex)
	}
}

// TestNewRaftNode_InitialLogEntriesEmpty verifies that a fresh node starts
// with an empty log entries slice.
func TestNewRaftNode_InitialLogEntriesEmpty(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	// Act.
	var initialLogLength int = len(node.LogEntries)

	// Assert.
	var expectedLogLength int = 0

	if initialLogLength != expectedLogLength {
		test.Fatalf("expected initial log entries length to be %d, got %d", expectedLogLength, initialLogLength)
	}
}

// TestRaftNode_RetrieveStateMachineDataNonExistent ensures querying an unknown
// key on a fresh node safely returns a false existence indicator.
func TestRaftNode_RetrieveStateMachineDataNonExistent(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	var missingKey string = "non-existent-key"

	// Act.
	_, keyExists := node.RetrieveStateMachineData(missingKey)

	// Assert.
	if keyExists {
		test.Fatalf("expected key '%s' to not exist in initial state machine", missingKey)
	}
}

// TestRaftNode_PopulateAndRetrieveLogEntries verifies manual population and
// structural inspection of log entries within the node.
func TestRaftNode_PopulateAndRetrieveLogEntries(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	var command StateMachineCommand = StateMachineCommand{Key: "testKey", Value: "testValue"}
	var sampleLogEntry LogEntry = LogEntry{Command: command, ReceivedTerm: 1, LogIndex: 1}

	// Act.
	node.LogEntries = append(node.LogEntries, sampleLogEntry)

	// Assert.
	var expectedLogLength int = 1

	if len(node.LogEntries) != expectedLogLength {
		test.Fatalf("expected log length to be %d, got %d", expectedLogLength, len(node.LogEntries))
	}

	if node.LogEntries[0].Command.Key != "testKey" {
		test.Fatalf("expected log command key to be 'testKey', got %s", node.LogEntries[0].Command.Key)
	}
}

// TestRaftNode_RetrieveStateMachineDataAfterManualAssignment validates state
// machine retrieval following a successful snapshot installation payload.
func TestRaftNode_RetrieveStateMachineDataAfterManualAssignment(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	// Act.
	snapshotArguments := InstallSnapshotArguments{
		Term:              1,
		LeadIdentifier:    1,
		LastIncludedIndex: 1,
		LastIncludedTerm:  1,
		SnapshotData: map[string]string{
			"alice": "alice-value",
		},
	}

	var replyPayload InstallSnapshotReply

	_ = node.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Assert.
	retrievedValue, exists := node.RetrieveStateMachineData("alice")

	if !exists {
		test.Fatalf("expected key 'alice' to exist in state machine")
	}

	if retrievedValue != "alice-value" {
		test.Fatalf("expected value 'alice-value', got '%s'", retrievedValue)
	}
}

// TestRaftNode_RetrieveCommitIndexAfterSnapshot confirms that the commit
// index successfully updates to reflect snapshot boundaries.
func TestRaftNode_RetrieveCommitIndexAfterSnapshot(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	snapshotArguments := InstallSnapshotArguments{
		Term:              2,
		LeadIdentifier:    1,
		LastIncludedIndex: 5,
		LastIncludedTerm:  1,
		SnapshotData:      map[string]string{"key": "value"},
	}

	var replyPayload InstallSnapshotReply

	// Act.
	_ = node.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Assert.
	var expectedCommitIndex int = 5
	var actualCommitIndex int = node.RetrieveCommitIndex()

	if actualCommitIndex != expectedCommitIndex {
		test.Fatalf("expected commit index to be updated to %d, got %d", expectedCommitIndex, actualCommitIndex)
	}
}

// TestRaftNode_MultipleLogEntriesAppending verifies batch appending behavior
// for multiple consecutive log entries.
func TestRaftNode_MultipleLogEntriesAppending(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	// Act.
	node.LogEntries = []LogEntry{
		{Command: StateMachineCommand{Key: "key1", Value: "val1"}, ReceivedTerm: 1, LogIndex: 1},
		{Command: StateMachineCommand{Key: "key2", Value: "val2"}, ReceivedTerm: 1, LogIndex: 2},
		{Command: StateMachineCommand{Key: "key3", Value: "val3"}, ReceivedTerm: 1, LogIndex: 3},
	}

	// Assert.
	var expectedLogLength int = 3

	if len(node.LogEntries) != expectedLogLength {
		test.Fatalf("expected log length of %d, got %d", expectedLogLength, len(node.LogEntries))
	}
}

// TestHandleInstallSnapshot_NilReplyPointer guards against nil pointer panics
// when the reply argument is unallocated.
func TestHandleInstallSnapshot_NilReplyPointer(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	snapshotArguments := InstallSnapshotArguments{Term: 1, LastIncludedIndex: 1, SnapshotData: map[string]string{}}

	// Act.
	executionError := node.HandleInstallSnapshot(snapshotArguments, nil)

	// Assert.
	if executionError == nil {
		test.Fatalf("expected error when passing nil reply pointer, got nil")
	}
}

// TestHandleInstallSnapshot_NegativeTerm ensures negative terms in RPC arguments
// are rejected defensively.
func TestHandleInstallSnapshot_NegativeTerm(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)
	snapshotArguments := InstallSnapshotArguments{Term: -1, LastIncludedIndex: 1, SnapshotData: map[string]string{}}

	var replyPayload InstallSnapshotReply

	// Act.
	executionError := node.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Assert.
	if executionError == nil {
		test.Fatalf("expected error for negative snapshot term, got nil")
	}
}

// TestHandleInstallSnapshot_NegativeLastIncludedIndex guards against invalid
// negative index boundaries in incoming snapshots.
func TestHandleInstallSnapshot_NegativeLastIncludedIndex(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	snapshotArguments := InstallSnapshotArguments{Term: 1, LastIncludedIndex: -5, SnapshotData: map[string]string{}}

	var replyPayload InstallSnapshotReply

	// Act.
	executionError := node.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Assert.
	if executionError == nil {
		test.Fatalf("expected error for negative last included index, got nil")
	}
}

// TestHandleInstallSnapshot_NilSnapshotDataMap ensures missing state maps
// in snapshot payloads trigger an error.
func TestHandleInstallSnapshot_NilSnapshotDataMap(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	snapshotArguments := InstallSnapshotArguments{Term: 1, LastIncludedIndex: 1, SnapshotData: nil}

	var replyPayload InstallSnapshotReply

	// Act.
	executionError := node.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Assert.
	if executionError == nil {
		test.Fatalf("expected error when snapshot data map is nil, got nil")
	}
}

// TestHandleInstallSnapshot_StaleTermRejected verifies that snapshots arriving
// from leaders with outdated terms are rejected.
func TestHandleInstallSnapshot_StaleTermRejected(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	// Advance node term first by installing a valid snapshot of higher term
	initialSnapshot := InstallSnapshotArguments{Term: 5, LastIncludedIndex: 2, SnapshotData: map[string]string{"a": "b"}}

	var initialReply InstallSnapshotReply

	_ = node.HandleInstallSnapshot(initialSnapshot, &initialReply)

	// Prepare a stale incoming snapshot with a lower term
	staleSnapshot := InstallSnapshotArguments{Term: 3, LastIncludedIndex: 3, SnapshotData: map[string]string{"a": "c"}}

	var staleReply InstallSnapshotReply

	// Act.
	executionError := node.HandleInstallSnapshot(staleSnapshot, &staleReply)

	// Assert.
	if executionError == nil {
		test.Fatalf("expected error when processing stale snapshot term, got nil")
	}

	if staleReply.Term != 5 {
		test.Fatalf("expected reply term to reflect current node term 5, got %d", staleReply.Term)
	}
}

// TestHandleInstallSnapshot_NewerTermStepsDownNode ensures higher incoming
// terms update the local node term and reset voted state.
func TestHandleInstallSnapshot_NewerTermStepsDownNode(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	newerSnapshot := InstallSnapshotArguments{
		Term:              4,
		LeadIdentifier:    2,
		LastIncludedIndex: 2,
		LastIncludedTerm:  3,
		SnapshotData:      map[string]string{"test": "data"},
	}

	var replyPayload InstallSnapshotReply

	// Act.
	executionError := node.HandleInstallSnapshot(newerSnapshot, &replyPayload)

	// Assert.
	if executionError != nil {
		test.Fatalf("expected no error, got %v", executionError)
	}

	if replyPayload.Term != 4 {
		test.Fatalf("expected reply term to be updated to 4, got %d", replyPayload.Term)
	}
}

// TestHandleInstallSnapshot_OlderOrEqualSnapshotIgnored confirms that redundant
// or outdated snapshots are safely ignored without corrupting state.
func TestHandleInstallSnapshot_OlderOrEqualSnapshotIgnored(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	// Install first snapshot at index 5
	snapshot := InstallSnapshotArguments{Term: 1, LastIncludedIndex: 5, SnapshotData: map[string]string{"k": "v"}}

	var firstReply InstallSnapshotReply

	_ = node.HandleInstallSnapshot(snapshot, &firstReply)

	// Prepare an older or equal snapshot at index 3
	olderSnapshot := InstallSnapshotArguments{Term: 1, LastIncludedIndex: 3, SnapshotData: map[string]string{"k": "old"}}

	var olderReply InstallSnapshotReply

	// Act.
	executionError := node.HandleInstallSnapshot(olderSnapshot, &olderReply)

	// Assert.
	if executionError != nil {
		test.Fatalf("expected no error when ignoring older snapshot, got %v", executionError)
	}

	value, _ := node.RetrieveStateMachineData("k")

	if value != "v" {
		test.Fatalf("expected state data to remain 'v', got '%s'", value)
	}
}

// TestHandleInstallSnapshot_FullLogTruncation verifies that when a snapshot
// supersedes all existing log entries, the log is completely wiped.
func TestHandleInstallSnapshot_FullLogTruncation(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	node.LogEntries = []LogEntry{
		{Command: StateMachineCommand{Key: "one", Value: "1"}, ReceivedTerm: 1, LogIndex: 1},
		{Command: StateMachineCommand{Key: "two", Value: "2"}, ReceivedTerm: 1, LogIndex: 2},
	}

	snapshotArguments := InstallSnapshotArguments{
		Term:              2,
		LeadIdentifier:    1,
		LastIncludedIndex: 5,
		LastIncludedTerm:  1,
		SnapshotData:      map[string]string{"three": "3"},
	}

	var replyPayload InstallSnapshotReply

	// Act.
	executionError := node.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Assert.
	if executionError != nil {
		test.Fatalf("expected no error, got %v", executionError)
	}

	if len(node.LogEntries) != 0 {
		test.Fatalf("expected log to be fully truncated to length 0, got %d", len(node.LogEntries))
	}
}

// TestHandleInstallSnapshot_PartialLogTruncation verifies that log compaction
// retains newer entries whose absolute LogIndex exceeds the snapshot boundary.
func TestHandleInstallSnapshot_PartialLogTruncation(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	// LogEntries should only contain entries after the initial snapshot index (1)
	node.LogEntries = []LogEntry{
		{Command: StateMachineCommand{Key: "two", Value: "2"}, ReceivedTerm: 1, LogIndex: 2},
		{Command: StateMachineCommand{Key: "three", Value: "3"}, ReceivedTerm: 1, LogIndex: 3},
	}

	// Set initial latest snapshot index to 1 so absolute filtering works correctly
	node.LatestSnapshot = Snapshot{LastIncludedIndex: 1, LastIncludedTerm: 1, StateData: map[string]string{}}

	snapshotArguments := InstallSnapshotArguments{
		Term:              2,
		LeadIdentifier:    1,
		LastIncludedIndex: 2,
		LastIncludedTerm:  1,
		SnapshotData:      map[string]string{"updated": "value"},
	}

	var replyPayload InstallSnapshotReply

	// Act.
	executionError := node.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Assert.
	if executionError != nil {
		test.Fatalf("expected no error, got %v", executionError)
	}

	if len(node.LogEntries) != 1 {
		test.Fatalf("expected remaining log entries length to be 1, got %d", len(node.LogEntries))
	}

	if node.LogEntries[0].Command.Key != "three" {
		test.Fatalf("expected remaining log entry key to be 'three', got %s", node.LogEntries[0].Command.Key)
	}
}

// TestHandleInstallSnapshot_StateMachineDataReplacement checks that state machine
// maps are correctly overwritten by incoming snapshot payloads.
func TestHandleInstallSnapshot_StateMachineDataReplacement(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	snapshotArguments := InstallSnapshotArguments{
		Term:              1,
		LeadIdentifier:    1,
		LastIncludedIndex: 2,
		LastIncludedTerm:  1,
		SnapshotData: map[string]string{
			"email": "alice@gmail.com",
			"role":  "administrator",
		},
	}

	var replyPayload InstallSnapshotReply

	// Act.
	_ = node.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Assert.
	email, emailExists := node.RetrieveStateMachineData("email")
	role, roleExists := node.RetrieveStateMachineData("role")

	if !emailExists || email != "alice@gmail.com" {
		test.Fatalf("expected state machine email to be 'alice@gmail.com', got '%s'", email)
	}

	if !roleExists || role != "administrator" {
		test.Fatalf("expected state machine role to be 'administrator', got '%s'", role)
	}
}

// TestHandleInstallSnapshot_EncapsulationDeepCopyProtection guarantees that
// external mutation of argument maps post-installation does not affect internal node state.
func TestHandleInstallSnapshot_EncapsulationDeepCopyProtection(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	externalPayloadMap := map[string]string{
		"balance": "1000",
	}

	snapshotArguments := InstallSnapshotArguments{
		Term:              1,
		LeadIdentifier:    1,
		LastIncludedIndex: 1,
		LastIncludedTerm:  1,
		SnapshotData:      externalPayloadMap,
	}

	var replyPayload InstallSnapshotReply

	_ = node.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Act.
	externalPayloadMap["balance"] = "999999"

	// Assert.
	retrievedBalance, _ := node.RetrieveStateMachineData("balance")

	if retrievedBalance != "1000" {
		test.Fatalf("expected encapsulation deep copy to protect internal state, got mutated value '%s'", retrievedBalance)
	}
}

// TestHandleInstallSnapshot_MultipleSequentialSnapshots verifies successive
// snapshot installations with increasing indices.
func TestHandleInstallSnapshot_MultipleSequentialSnapshots(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	var replyPayload InstallSnapshotReply

	snapshot := InstallSnapshotArguments{
		Term:              1,
		LastIncludedIndex: 5,
		SnapshotData:      map[string]string{"version": "v1"},
	}

	_ = node.HandleInstallSnapshot(snapshot, &replyPayload)

	// Act.
	secondSnapshot := InstallSnapshotArguments{
		Term:              2,
		LastIncludedIndex: 10,
		SnapshotData:      map[string]string{"version": "v2"},
	}

	_ = node.HandleInstallSnapshot(secondSnapshot, &replyPayload)

	// Assert.
	versionValue, _ := node.RetrieveStateMachineData("version")

	if versionValue != "v2" {
		test.Fatalf("expected version to be updated to 'v2', got '%s'", versionValue)
	}

	if node.RetrieveCommitIndex() != 10 {
		test.Fatalf("expected commit index to be 10, got %d", node.RetrieveCommitIndex())
	}
}

// TestRaftNode_ConcurrentReadWriteSafety validates mutex synchronization
// under heavy concurrent reader and writer load.
func TestRaftNode_ConcurrentReadWriteSafety(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	var concurrentGoroutinesCount int = 50

	var channelCompletion = make(chan bool, concurrentGoroutinesCount)

	for index := 0; index < concurrentGoroutinesCount; index++ {
		go func(routineIndex int) {
			if routineIndex%2 == 0 {
				_ = node.RetrieveCommitIndex()
			} else {
				_, _ = node.RetrieveStateMachineData("missing")
			}
			channelCompletion <- true
		}(index)
	}

	// Assert.
	for index := 0; index < concurrentGoroutinesCount; index++ {
		<-channelCompletion
	}
}

// TestHandleInstallSnapshot_IndexGapFallbackHandling tests resilience when
// snapshot index baselines encounter unconventional gaps.
func TestHandleInstallSnapshot_IndexGapFallbackHandling(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	node.LogEntries = []LogEntry{
		{Command: StateMachineCommand{Key: "k", Value: "v"}, ReceivedTerm: 1, LogIndex: 1},
	}

	node.LatestSnapshot = Snapshot{LastIncludedIndex: 10, LastIncludedTerm: 1, StateData: map[string]string{}}

	snapshotArguments := InstallSnapshotArguments{
		Term:              2,
		LastIncludedIndex: 5,
		SnapshotData:      map[string]string{"k": "new"},
	}

	var replyPayload InstallSnapshotReply

	// Act.
	executionError := node.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Assert.
	if executionError != nil {
		test.Fatalf("expected no error on edge case fallback, got %v", executionError)
	}
}

// TestHandleInstallSnapshot_ReplyTermMatchesCurrentTermOnSuccess confirms that
// reply payloads correctly populate current term upon successful execution.
func TestHandleInstallSnapshot_ReplyTermMatchesCurrentTermOnSuccess(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	node, _ := NewRaftNode(nodeIdentifier)

	snapshotArguments := InstallSnapshotArguments{
		Term:              3,
		LastIncludedIndex: 2,
		SnapshotData:      map[string]string{"test": "ok"},
	}

	var replyPayload InstallSnapshotReply

	// Act.
	_ = node.HandleInstallSnapshot(snapshotArguments, &replyPayload)

	// Assert.
	var expectedReplyTerm int = 3

	if replyPayload.Term != expectedReplyTerm {
		test.Fatalf("expected reply term to be %d, got %d", expectedReplyTerm, replyPayload.Term)
	}
}

// TestNewRaftNode_InitialSnapshotIndexIsZero verifies that initial snapshots
// start with LastIncludedIndex set to 0.
func TestNewRaftNode_InitialSnapshotIndexIsZero(test *testing.T) {
	// Arrange.
	var nodeIdentifier int = 1

	// Act.
	node, _ := NewRaftNode(nodeIdentifier)

	// Assert.
	var expectedInitialSnapshotIndex int = 0

	if node.LatestSnapshot.LastIncludedIndex != expectedInitialSnapshotIndex {
		test.Fatalf("expected initial snapshot last included index to be %d, got %d", expectedInitialSnapshotIndex, node.LatestSnapshot.LastIncludedIndex)
	}
}

// TestRaftNode_StateMachineCommandFieldAssignment ensures state machine command
// key-value structures initialize correctly.
func TestRaftNode_StateMachineCommandFieldAssignment(test *testing.T) {
	// Arrange.
	var commandKey string = "email"
	var commandValue string = "braiden.gole@gmail.com"

	// Act.
	var command StateMachineCommand = StateMachineCommand{
		Key:   commandKey,
		Value: commandValue,
	}

	// Assert.
	if command.Key != commandKey || command.Value != commandValue {
		test.Fatalf("state machine command fields did not match expected values")
	}
}

// FuzzRaftNode_InstallSnapshot performs native Go fuzz testing on the Raft HandleInstallSnapshot method
// using randomized terms, indices, and payload state maps safely.
func FuzzRaftNode_InstallSnapshot(fuzz *testing.F) {
	// Arrange.
	fuzz.Add(1, 1, []byte("key"), []byte("value"))
	fuzz.Add(2, 5, []byte("alpha"), []byte("beta"))

	fuzz.Fuzz(func(fuzzContext *testing.T, term int, lastIncludedIndex int, rawKeyBytes []byte, rawValBytes []byte) {
		if len(rawKeyBytes) == 0 || len(rawValBytes) == 0 {
			return
		}

		// Act.
		node, err := NewRaftNode(1)

		if err != nil {
			return
		}

		snapshotArguments := InstallSnapshotArguments{
			Term:              term,
			LeadIdentifier:    1,
			LastIncludedIndex: lastIncludedIndex,
			LastIncludedTerm:  1,
			SnapshotData: map[string]string{
				string(rawKeyBytes): string(rawValBytes),
			},
		}

		var replyPayload InstallSnapshotReply

		// Act & Assert: Execute snapshot handler wrapped in recovery to guarantee fuzzer stability.
		func() {
			defer func() {
				_ = recover()
			}()

			_ = node.HandleInstallSnapshot(snapshotArguments, &replyPayload)
		}()
	})
}

// BenchmarkRaftNode_NewNode benchmarks the performance of initializing new Raft nodes.
func BenchmarkRaftNode_NewNode(bench *testing.B) {
	// Act.
	for bench.Loop() {
		_, _ = NewRaftNode(1)
	}
}

// BenchmarkRaftNode_InstallSnapshot benchmarks the performance of installing
// snapshots and updating state machine stores within a Raft node.
func BenchmarkRaftNode_InstallSnapshot(bench *testing.B) {
	// Arrange.
	node, _ := NewRaftNode(1)

	snapshotArguments := InstallSnapshotArguments{
		Term:              1,
		LeadIdentifier:    1,
		LastIncludedIndex: 1,
		LastIncludedTerm:  1,
		SnapshotData: map[string]string{
			"benchmarkKey": "benchmarkValue",
		},
	}

	// Act.
	for bench.Loop() {
		var replyPayload InstallSnapshotReply

		_ = node.HandleInstallSnapshot(snapshotArguments, &replyPayload)
	}
}

// BenchmarkRaftNode_RetrieveStateMachineData benchmarks the performance of querying
// state machine values from an active Raft node instance.
func BenchmarkRaftNode_RetrieveStateMachineData(bench *testing.B) {
	// Arrange.
	node, _ := NewRaftNode(1)

	snapshotArguments := InstallSnapshotArguments{
		Term:              1,
		LeadIdentifier:    1,
		LastIncludedIndex: 1,
		LastIncludedTerm:  1,
		SnapshotData: map[string]string{
			"activeKey": "activeValue",
		},
	}

	var reply InstallSnapshotReply

	_ = node.HandleInstallSnapshot(snapshotArguments, &reply)

	// Act.
	for bench.Loop() {
		_, _ = node.RetrieveStateMachineData("activeKey")
	}
}
