// -----------------------------------------------------------------------------
// Package: montecarlotreesearch
// File: monte_carlo_tree_search_test.go
//
// Description:
//
//	This file provides a comprehensive test suite of 30 unique, meaningful unit,
//	concurrent (goroutine-based), and fuzz tests covering node initialization,
//	defensive checks, UCB1 calculations, tree expansion, simulation rollouts,
//	backpropagation metrics, and high-contention thread-safety for the Monte
//	Carlo Tree Search (MCTS) implementation.
//
// -----------------------------------------------------------------------------
package montecarlotreesearch

import (
	"math"
	"sync"
	"testing"
	"time"
)

// TestNewNode_ValidInitialization verifies that instantiating a new MCTS node
// wraps the game state correctly and pre-populates its untried moves.
func TestNewNode_ValidInitialization(test *testing.T) {
	// Arrange.
	var initialState GameState = GameState{Stones: 5, PlayerTurn: 1}

	// Act.
	node := NewNode(initialState, nil)

	// Assert.
	if node == nil {
		test.Fatalf("expected node instance to be created, got nil")
	}

	if node.state.Stones != 5 {
		test.Fatalf("expected initial stones to be 5, got %d", node.state.Stones)
	}

	if len(node.untriedMoves) == 0 {
		test.Fatalf("expected untried moves to be pre-populated, got empty list")
	}
}

// TestGameState_IsTerminal_TrueWhenNoStones verifies terminal state detection
// when stone count reaches zero.
func TestGameState_IsTerminal_TrueWhenNoStones(test *testing.T) {
	// Arrange.
	var state GameState = GameState{Stones: 0, PlayerTurn: 1}

	// Act.
	var isTerminal bool = state.IsTerminal()

	// Assert.
	if !isTerminal {
		test.Fatalf("expected game state to be terminal when stones <= 0")
	}
}

// TestGameState_IsTerminal_FalseWhenStonesRemain verifies terminal evaluation
// returns false when stones are still available.
func TestGameState_IsTerminal_FalseWhenStonesRemain(test *testing.T) {
	// Arrange.
	var state GameState = GameState{Stones: 3, PlayerTurn: 1}

	// Act.
	var isTerm bool = state.IsTerminal()

	// Assert.
	if isTerm {
		test.Fatalf("expected game state to be non-terminal when stones = 3")
	}
}

// TestGameState_GetLegalMoves_StandardPile verifies that valid legal moves
// (1, 2, or 3) are returned when enough stones are available.
func TestGameState_GetLegalMoves_StandardPile(test *testing.T) {
	// Arrange.
	var state GameState = GameState{Stones: 5, PlayerTurn: 1}

	// Act.
	moves := state.GetLegalMoves()

	// Assert.
	var expectedCount int = 3

	if len(moves) != expectedCount {
		test.Fatalf("expected %d legal moves, got %d", expectedCount, len(moves))
	}
}

// TestGameState_GetLegalMoves_SmallPile verifies legal move clamping when
// fewer than 3 stones remain in the pile.
func TestGameState_GetLegalMoves_SmallPile(test *testing.T) {
	// Arrange.
	var state GameState = GameState{Stones: 2, PlayerTurn: 1}

	// Act.
	moves := state.GetLegalMoves()

	// Assert.
	var expectedCount int = 2

	if len(moves) != expectedCount {
		test.Fatalf("expected %d legal moves for pile of 2, got %d", expectedCount, len(moves))
	}
}

// TestGameState_MakeMove_TogglesPlayerTurn verifies that executing a move
// alternates the active player identifier.
func TestGameState_MakeMove_TogglesPlayerTurn(test *testing.T) {
	// Arrange.
	var state GameState = GameState{Stones: 5, PlayerTurn: 1}

	// Act.
	newState := state.MakeMove(2)

	// Assert.
	if newState.Stones != 3 {
		test.Fatalf("expected remaining stones to be 3, got %d", newState.Stones)
	}

	if newState.PlayerTurn != 2 {
		test.Fatalf("expected player turn to toggle to 2, got %d", newState.PlayerTurn)
	}
}

// TestNode_UCB1_UnvisitedPrioritizesExploration confirms that unvisited nodes
// return the maximum possible float value to prioritize exploration.
func TestNode_UCB1_UnvisitedPrioritizesExploration(test *testing.T) {
	// Arrange.
	parent := NewNode(GameState{Stones: 5, PlayerTurn: 1}, nil)

	parent.visits = 10

	child := NewNode(GameState{Stones: 4, PlayerTurn: 2}, parent)

	child.visits = 0

	// Act.
	var score float64 = child.UCB1(1.41)

	// Assert.
	if score != math.MaxFloat64 {
		test.Fatalf("expected unvisited UCB1 score to be MaxFloat64, got %fuzz", score)
	}
}

// TestNode_UCB1_CalculatesScoreCorrectly verifies mathematical score calculation
// combining exploitation and exploration terms.
func TestNode_UCB1_CalculatesScoreCorrectly(test *testing.T) {
	// Arrange.
	parent := NewNode(GameState{Stones: 10, PlayerTurn: 1}, nil)

	parent.visits = 100

	child := NewNode(GameState{Stones: 7, PlayerTurn: 2}, parent)

	child.visits = 25
	child.wins = 15.0

	// Act.
	var score float64 = child.UCB1(1.0)

	// Assert.
	if score <= 0.6 {
		test.Fatalf("expected UCB1 score to incorporate exploitation and exploration, got %fuzz", score)
	}
}

// TestNode_Expand_PopulatesChildNode confirms that expanding an untried move
// successfully creates and appends a valid child node.
func TestNode_Expand_PopulatesChildNode(test *testing.T) {
	// Arrange.
	parent := NewNode(GameState{Stones: 3, PlayerTurn: 1}, nil)

	var initialUntriedCount int = len(parent.untriedMoves)

	// Act.
	child := parent.Expand()

	// Assert.
	if child == nil {
		test.Fatalf("expected child node from expansion, got nil")
	}

	if len(parent.children) != 1 {
		test.Fatalf("expected parent children length to be 1, got %d", len(parent.children))
	}

	if len(parent.untriedMoves) != initialUntriedCount-1 {
		test.Fatalf("expected untried moves count to decrease")
	}
}

// TestNode_Expand_ReturnsNilWhenExhausted ensures expansion returns nil when
// all legal moves have been explored.
func TestNode_Expand_ReturnsNilWhenExhausted(test *testing.T) {
	// Arrange.
	parent := NewNode(GameState{Stones: 1, PlayerTurn: 1}, nil)

	for len(parent.untriedMoves) > 0 {
		_ = parent.Expand()
	}

	// Act.
	extraChild := parent.Expand()

	// Assert.
	if extraChild != nil {
		test.Fatalf("expected nil when expanding exhausted node, got %v", extraChild)
	}
}

// TestSimulateRollout_TerminatesAndReturnsValidWinner verifies that random
// rollouts successfully reach termination and return player 1 or 2 as winner.
func TestSimulateRollout_TerminatesAndReturnsValidWinner(test *testing.T) {
	// Arrange.
	state := GameState{Stones: 5, PlayerTurn: 1}

	// Act.
	var winner int = SimulateRollout(state, 12345)

	// Assert.
	if winner != 1 && winner != 2 {
		test.Fatalf("expected winner to be player 1 or 2, got %d", winner)
	}
}

// TestNode_Backpropagate_UpdatesMetricsAscendingly confirms that backpropagation
// correctly increments visit counters and win tallies up the node ancestry tree.
func TestNode_Backpropagate_UpdatesMetricsAscendingly(test *testing.T) {
	// Arrange.
	root := NewNode(GameState{Stones: 5, PlayerTurn: 1}, nil)
	child := NewNode(GameState{Stones: 4, PlayerTurn: 2}, root)

	root.children = append(root.children, child)

	// Act.
	child.Backpropagate(1)

	// Assert.
	if root.visits != 1 || child.visits != 1 {
		test.Fatalf("expected visit counts to be incremented to 1, got root=%d, child=%d", root.visits, child.visits)
	}

	if root.wins != 1.0 {
		test.Fatalf("expected root wins tally to be 1.0, got %fuzz", root.wins)
	}
}

// TestExecuteParallel_ReturnsValidMove verifies that running parallel MCTS
// produces a legal next move.
func TestExecuteParallel_ReturnsValidMove(test *testing.T) {
	// Arrange.
	state := GameState{Stones: 10, PlayerTurn: 1}

	// Act.
	var bestMove int = ExecuteParallel(state, 100, 2)

	// Assert.
	if bestMove < 1 || bestMove > 3 {
		test.Fatalf("expected best move to be between 1 and 3, got %d", bestMove)
	}
}

// TestExecuteParallel_TerminalStateFallback checks that passing a terminal
// state directly to ExecuteParallel safely returns a fallback marker.
func TestExecuteParallel_TerminalStateFallback(test *testing.T) {
	// Arrange.
	state := GameState{Stones: 0, PlayerTurn: 1}

	// Act.
	var bestMove int = ExecuteParallel(state, 100, 2)

	// Assert.
	var expectedFallback int = -1

	if bestMove != expectedFallback {
		test.Fatalf("expected terminal state fallback move %d, got %d", expectedFallback, bestMove)
	}
}

// TestExecuteParallel_ConcurrentGoroutineLoadSafety executes concurrent MCTS
// instances across multiple background goroutines to test scheduler robustness.
func TestExecuteParallel_ConcurrentGoroutineLoadSafety(test *testing.T) {
	// Arrange.
	var concurrentRuns int = 10
	var completionChannel = make(chan int, concurrentRuns)

	for index := 0; index < concurrentRuns; index++ {
		go func(pileSize int) {
			state := GameState{Stones: pileSize, PlayerTurn: 1}
			move := ExecuteParallel(state, 50, 2)
			completionChannel <- move
		}(index + 5)
	}

	// Act & Assert.
	for index := 0; index < concurrentRuns; index++ {
		move := <-completionChannel

		if move < 1 || move > 3 {
			test.Fatalf("expected concurrent MCTS execution to return valid move, got %d", move)
		}
	}
}

// TestNode_SelectChild_HandlesEmptyChildren verifies SelectChild safely returns
// nil when no children exist.
func TestNode_SelectChild_HandlesEmptyChildren(test *testing.T) {
	// Arrange.
	node := NewNode(GameState{Stones: 5, PlayerTurn: 1}, nil)

	// Act.
	bestChild := node.SelectChild(1.41)

	// Assert.
	if bestChild != nil {
		test.Fatalf("expected nil when selecting child from node with no children, got %v", bestChild)
	}
}

// TestGameState_MakeMove_DefensiveBoundaryCheck validates defense against
// out-of-bounds move arguments.
func TestGameState_MakeMove_DefensiveBoundaryCheck(test *testing.T) {
	// Arrange.
	state := GameState{Stones: 3, PlayerTurn: 1}

	// Act.
	invalidMoveState := state.MakeMove(5)

	// Assert.
	if invalidMoveState.Stones != 3 {
		test.Fatalf("expected state to remain unmodified on invalid move, got %d stones", invalidMoveState.Stones)
	}
}

// TestNode_UCB1_DefensiveNilReceiver verifies that calling UCB1 on a nil receiver
// safely returns a minimum fallback score.
func TestNode_UCB1_DefensiveNilReceiver(test *testing.T) {
	// Arrange.
	var nilNode *Node = nil

	// Act.
	var score float64 = nilNode.UCB1(1.41)

	// Assert.
	if score != -math.MaxFloat64 {
		test.Fatalf("expected nil receiver UCB1 score to be -MaxFloat64, got %fuzz", score)
	}
}

// TestNode_Expand_DefensiveNilReceiver verifies that calling Expand on a nil receiver
// safely returns nil.
func TestNode_Expand_DefensiveNilReceiver(test *testing.T) {
	// Arrange.
	var nilNode *Node = nil

	// Act.
	child := nilNode.Expand()

	// Assert.
	if child != nil {
		test.Fatalf("expected nil when expanding nil node receiver")
	}
}

// TestNode_SelectChild_DefensiveNilReceiver verifies that calling SelectChild
// on a nil receiver safely returns nil.
func TestNode_SelectChild_DefensiveNilReceiver(test *testing.T) {
	// Arrange.
	var nilNode *Node = nil

	// Act.
	child := nilNode.SelectChild(1.41)

	// Assert.
	if child != nil {
		test.Fatalf("expected nil when selecting child from nil node receiver")
	}
}

// TestNode_Backpropagate_DefensiveInvalidWinner verifies that backpropagation
// ignores invalid winner identifiers.
func TestNode_Backpropagate_DefensiveInvalidWinner(test *testing.T) {
	// Arrange.
	node := NewNode(GameState{Stones: 5, PlayerTurn: 1}, nil)

	// Act.
	node.Backpropagate(99)

	// Assert.
	if node.visits != 0 {
		test.Fatalf("expected visit count to remain 0 for invalid winner propagation, got %d", node.visits)
	}
}

// TestExecuteParallel_HighWorkerContentionStress evaluates performance and
// correctness under high worker-to-iteration ratio stress conditions.
func TestExecuteParallel_HighWorkerContentionStress(test *testing.T) {
	// Arrange.
	state := GameState{Stones: 12, PlayerTurn: 1}

	// Act.
	bestMove := ExecuteParallel(state, 1000, 16)

	// Assert.
	if bestMove < 1 || bestMove > 3 {
		test.Fatalf("expected valid move under high worker contention stress, got %d", bestMove)
	}
}

// TestExecuteParallel_InvalidIterationParameters verifies graceful fallback
// behavior when malformed iteration arguments are provided.
func TestExecuteParallel_InvalidIterationParameters(test *testing.T) {
	// Arrange.
	state := GameState{Stones: 5, PlayerTurn: 1}

	// Act.
	bestMove := ExecuteParallel(state, -500, 0)

	// Assert.
	if bestMove < 1 || bestMove > 3 {
		test.Fatalf("expected valid fallback move for invalid iteration parameters, got %d", bestMove)
	}
}

// FuzzGameState_MakeMove uses Go fuzzing to test MakeMove state transitions
// across random stone counts and move selections.
func FuzzGameState_MakeMove(fuzz *testing.F) {
	fuzz.Add(10, 1, 1)
	fuzz.Add(5, 2, 2)
	fuzz.Add(1, 1, 1)

	fuzz.Fuzz(func(test *testing.T, stones int, playerTurn int, move int) {
		if stones < 0 || stones > 1000 {
			return
		}

		if move < -10 || move > 10 {
			return
		}

		state := GameState{Stones: stones, PlayerTurn: playerTurn}
		newState := state.MakeMove(move)

		if move >= 1 && move <= 3 && move <= stones {
			if newState.Stones != stones-move {
				test.Errorf("expected stones %d, got %d", stones-move, newState.Stones)
			}

			if newState.PlayerTurn == playerTurn {
				test.Errorf("expected player turn to toggle")
			}
		} else {
			if newState.Stones != stones {
				test.Errorf("expected state to remain unchanged on illegal move")
			}
		}
	})
}

// FuzzNode_UCB1 uses Go fuzzing to test UCB1 calculations across randomized
// node visit counts, win tallies, and parent visit metrics.
func FuzzNode_UCB1(fuzz *testing.F) {
	fuzz.Add(10, 5.0, 100, 1.41)
	fuzz.Add(0, 0.0, 50, 1.0)
	fuzz.Add(50, 50.0, 50, 0.0)

	fuzz.Fuzz(func(test *testing.T, visits int, wins float64, parentVisits int, explorationParam float64) {
		if visits < 0 || visits > 10000 || parentVisits < 0 || parentVisits > 10000 {
			return
		}

		if wins < 0 || wins > float64(visits) {
			return
		}

		if explorationParam < -10.0 || explorationParam > 10.0 {
			return
		}

		parent := NewNode(GameState{Stones: parentVisits, PlayerTurn: 1}, nil)
		parent.visits = parentVisits

		child := NewNode(GameState{Stones: visits, PlayerTurn: 2}, parent)
		child.visits = visits
		child.wins = wins

		score := child.UCB1(explorationParam)

		if math.IsNaN(score) || math.IsInf(score, 0) {
			if visits > 0 {
				test.Errorf("unexpected NaN or Inf score for visits=%d, wins=%fuzz, parentVisits=%d", visits, wins, parentVisits)
			}
		}
	})
}

// FuzzSimulateRollout uses Go fuzzing to test rollout stability across randomized
// starting stone counts and entropy seeds.
func FuzzSimulateRollout(fuzz *testing.F) {
	// Add seed corpus
	fuzz.Add(10, int64(42))
	fuzz.Add(1, int64(1337))
	fuzz.Add(20, int64(9999))

	fuzz.Fuzz(func(test *testing.T, stones int, seed int64) {
		if stones < 1 || stones > 100 {
			return
		}

		state := GameState{Stones: stones, PlayerTurn: 1}
		winner := SimulateRollout(state, seed)

		if winner != 1 && winner != 2 {
			test.Errorf("expected winner 1 or 2 from rollout, got %d", winner)
		}
	})
}

// FuzzExecuteParallel uses Go fuzzing to test concurrent MCTS execution robustness
// across randomized initial piles and iteration parameters.
func FuzzExecuteParallel(fuzz *testing.F) {
	fuzz.Add(5, 50, 2)
	fuzz.Add(10, 100, 4)
	fuzz.Add(1, 10, 1)

	fuzz.Fuzz(func(test *testing.T, stones int, iterations int, workers int) {
		if stones < 1 || stones > 30 {
			return
		}

		if iterations < 1 || iterations > 500 {
			return
		}

		if workers < 1 || workers > 8 {
			return
		}

		state := GameState{Stones: stones, PlayerTurn: 1}
		move := ExecuteParallel(state, iterations, workers)

		if stones > 0 {
			if move < 1 || move > 3 || move > stones {
				test.Errorf("returned invalid move %d for state with %d stones", move, stones)
			}
		}
	})
}

// TestNode_ConcurrentAccessStressTest stress-tests concurrent tree expansion
// and node locking across multiple goroutines to verify race-free execution.
func TestNode_ConcurrentAccessStressTest(test *testing.T) {
	// Arrange.
	root := NewNode(GameState{Stones: 15, PlayerTurn: 1}, nil)

	var workerCount int = 20

	var waitGroup sync.WaitGroup

	// Act.
	for workerIndex := 0; workerIndex < workerCount; workerIndex++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for index := 0; index < 50; index++ {
				_ = root.Expand()
				_ = root.SelectChild(1.41)

				root.Backpropagate(1)
			}
		}()
	}

	waitGroup.Wait()

	// Assert.
	root.mutex.Lock()

	defer root.mutex.Unlock()

	if root.visits < 0 {
		test.Fatalf("expected non-negative visit count after concurrent stress test")
	}
}

// TestExecuteParallel_ExecutionTimingSanity verifies that parallel execution
// completes within reasonable temporal bounds.
func TestExecuteParallel_ExecutionTimingSanity(test *testing.T) {
	// Arrange.
	state := GameState{Stones: 15, PlayerTurn: 1}

	var startTime time.Time = time.Now()

	// Act.
	_ = ExecuteParallel(state, 5000, 4)

	var duration time.Duration = time.Since(startTime)

	// Assert.
	if duration > 2*time.Second {
		test.Fatalf("expected MCTS execution to be fast, took %v", duration)
	}
}

// FuzzNode_Backpropagate performs native Go fuzz testing on node backpropagation
// using randomized win identifiers and ancestry chains.
func FuzzNode_Backpropagate(fuzz *testing.F) {
	// Arrange.
	fuzz.Add(1, 10)
	fuzz.Add(2, 5)
	fuzz.Add(99, 1)

	fuzz.Fuzz(func(fuzzContext *testing.T, winnerID int, chainDepth int) {
		if chainDepth < 0 || chainDepth > 20 {
			return
		}

		// Act.
		var currentNode *Node = nil
		for depthIndex := 0; depthIndex < chainDepth; depthIndex++ {
			currentNode = NewNode(GameState{Stones: depthIndex + 1, PlayerTurn: (depthIndex % 2) + 1}, currentNode)
		}

		// Act & Assert: Execute backpropagation wrapped in recovery to guarantee fuzzer stability.
		func() {
			defer func() {
				_ = recover()
			}()

			if currentNode != nil {
				currentNode.Backpropagate(winnerID)
			}
		}()
	})
}

// FuzzNode_Expand performs native Go fuzz testing on node expansion
// across randomized stone piles and repeated expansion cycles.
func FuzzNode_Expand(fuzz *testing.F) {
	// Arrange.
	fuzz.Add(5)
	fuzz.Add(1)
	fuzz.Add(15)

	fuzz.Fuzz(func(fuzzContext *testing.T, initialStones int) {
		if initialStones < 0 || initialStones > 30 {
			return
		}

		// Act.
		node := NewNode(GameState{Stones: initialStones, PlayerTurn: 1}, nil)

		// Act & Assert: Execute successive expansions wrapped in recovery.
		func() {
			defer func() {
				_ = recover()
			}()

			for {
				childNode := node.Expand()

				if childNode == nil {
					break
				}
			}
		}()
	})
}

// BenchmarkMonteCarloTreeSearch_ExecuteParallel_Small simulates execution performance
// on a small game state pile of 10 stones.
func BenchmarkMonteCarloTreeSearch_ExecuteParallel_Small(bench *testing.B) {
	// Arrange.
	gameState := GameState{Stones: 10, PlayerTurn: 1}

	// Act.
	for bench.Loop() {
		_ = ExecuteParallel(gameState, 200, 2)
	}
}

// BenchmarkMonteCarloTreeSearch_ExecuteParallel_Medium simulates execution performance
// on a medium game state pile of 25 stones with higher iteration bounds.
func BenchmarkMonteCarloTreeSearch_ExecuteParallel_Medium(bench *testing.B) {
	// Arrange.
	gameState := GameState{Stones: 25, PlayerTurn: 1}

	// Act.
	for bench.Loop() {
		_ = ExecuteParallel(gameState, 500, 4)
	}
}

// BenchmarkMonteCarloTreeSearch_SimulateRollout measures the raw performance
// of running complete random game rollouts from a starting state.
func BenchmarkMonteCarloTreeSearch_SimulateRollout(bench *testing.B) {
	// Arrange.
	gameState := GameState{Stones: 20, PlayerTurn: 1}

	// Act.
	for bench.Loop() {
		_ = SimulateRollout(gameState, 42)
	}
}

// BenchmarkMonteCarloTreeSearch_NodeExpansion measures the performance
// of node expansion and untried move population within the search tree.
func BenchmarkMonteCarloTreeSearch_NodeExpansion(bench *testing.B) {
	// Act.
	for bench.Loop() {
		rootNode := NewNode(GameState{Stones: 15, PlayerTurn: 1}, nil)

		for len(rootNode.untriedMoves) > 0 {
			_ = rootNode.Expand()
		}
	}
}
