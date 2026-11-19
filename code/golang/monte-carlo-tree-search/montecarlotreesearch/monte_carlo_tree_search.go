// -----------------------------------------------------------------------------
// Package: montecarlotreesearch
// File: monte_carlo_tree_search.go
//
// Description:
//
//	This file implements a concurrent Monte Carlo Tree Search (MCTS) algorithm
//	in Go, tailored for game-tree search scenarios (demonstrated here via the
//	game of Nim). It features thread-safe tree expansion, selection, and traversal
//	using fine-grained node-level mutual exclusion locks (sync.Mutex) and parallel
//	worker pools coordinated via sync.WaitGroup.
//
//	The implementation incorporates the Upper Confidence Bound 1 (UCB1) formula
//	for balancing exploration and exploitation during tree policy selection,
//	randomized rollout simulations, and recursive backpropagation of scores.
//
// -----------------------------------------------------------------------------
package montecarlotreesearch

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// GameState encapsulates the state of the game being solved via MCTS
// (specifically the classic game of Nim). It tracks the remaining resource
// count and the identifier of the player whose turn it is to move.
type GameState struct {
	Stones     int
	PlayerTurn int
}

// Node represents a single node within the MCTS search tree. It maintains
// a reference to its game state, parent node, child nodes, unexpanded legal moves,
// visit counters, win tallies, and an internal mutex to ensure thread safety
// across concurrent worker goroutines.
type Node struct {
	state        GameState
	parent       *Node
	children     []*Node
	untriedMoves []int
	visits       int
	wins         float64
	mutex        sync.Mutex
}

// IsTerminal evaluates whether the current game state represents a terminal
// (game-over) condition where no further moves can be made.
func (state GameState) IsTerminal() bool {
	// Check if the remaining stone count has dropped to zero or below.
	// Return true if the game is over, otherwise return false.
	return state.Stones <= 0
}

// GetLegalMoves computes and returns all legal moves (e.g., taking 1, 2, or 3 stones)
// available from the current game state.
func (state GameState) GetLegalMoves() []int {
	// Initialize a slice to hold valid legal moves for the current state.
	var moves []int

	// Guard against terminal or negative stone states with a defensive check.
	if state.Stones <= 0 {
		return moves
	}

	// Iterate through possible actions of taking 1, 2, or 3 stones.
	for index := 1; index <= 3; index++ {
		// Verify that enough stones remain in the pile to execute the move.
		if state.Stones >= index {
			// Append the valid move to the list of legal actions.
			moves = append(moves, index)
		}
	}

	// Return the accumulated slice of legal moves.
	return moves
}

// MakeMove applies a specified move to the current state, returning a new
// resulting game state and alternating the active player turn.
func (state GameState) MakeMove(move int) GameState {
	// Validate that the move is within legal bounds and does not exceed available stones.
	if move < 1 || move > 3 || move > state.Stones {
		// Return the unmodified state as a safe defensive fallback for invalid moves.
		return state
	}

	// Construct and return a new GameState with the updated stone count and toggled player turn.
	return GameState{
		Stones:     state.Stones - move,
		PlayerTurn: 3 - state.PlayerTurn,
	}
}

// NewNode initializes and returns a new MCTS tree node wrapping the provided
// game state, linking it to its parent node, and pre-populating its untried legal moves.
func NewNode(state GameState, parent *Node) *Node {
	// Allocate and initialize a new search tree node with the given state and parent reference.
	// Populate untried moves by querying legal actions from the game state.
	return &Node{
		state:        state,
		parent:       parent,
		untriedMoves: state.GetLegalMoves(),
	}
}

// UCB1 calculates the Upper Confidence Bound 1 score for the node, balancing
// exploitation (win rate) with exploration (frequency of visits relative to the parent).
// Unvisited nodes return a maximum float value to guarantee immediate exploration priority.
func (node *Node) UCB1(explorationParam float64) float64 {
	// Guard against nil receiver references with a defensive check.
	if node == nil {
		return -math.MaxFloat64
	}

	// Handle unvisited nodes by returning the maximum possible float value to prioritize exploration.
	if node.visits == 0 {
		return math.MaxFloat64
	}

	// Guard against nil parent pointers (e.g., root node evaluation) to prevent runtime nil pointer panics.
	if node.parent == nil {
		return node.wins / float64(node.visits)
	}

	// Ensure the exploration parameter is non-negative.
	if explorationParam < 0 {
		explorationParam = 0
	}

	// Calculate the exploitation term using the average win rate of the node.
	// Calculate the exploration term using the natural logarithm of parent visits and node visits.
	// Combine both terms using the exploration parameter and return the final UCB1 score.
	return (node.wins / float64(node.visits)) + explorationParam*math.Sqrt(math.Log(float64(node.parent.visits))/float64(node.visits))
}

// SelectChild traverses down the search tree by locking the node and evaluating
// child nodes using the UCB1 selection formula, returning the child with the highest score.
func (node *Node) SelectChild(explorationParam float64) *Node {
	// Guard against nil receiver references.
	if node == nil {
		return nil
	}

	// Acquire the mutex lock to ensure thread-safe read and write access to child collections.
	node.mutex.Lock()

	// Ensure the mutex lock is released when the function exits.
	defer node.mutex.Unlock()

	// Guard against empty children slices.
	if len(node.children) == 0 {
		return nil
	}

	// Declare a variable to track the best-scoring child node.
	var bestChild *Node

	// Initialize the best score tracker to the lowest possible float value.
	bestScore := -math.MaxFloat64

	// Iterate through each child node to evaluate its UCB1 score.
	for _, child := range node.children {
		// Guard against nil child pointers in the collection.
		if child == nil {
			continue
		}

		// Lock the individual child node to safely read its visit and win metrics.
		child.mutex.Lock()

		// Compute the UCB1 score for the child node.
		score := child.UCB1(explorationParam)

		// Release the child node mutex.
		child.mutex.Unlock()

		// Update the best score and best child reference if a higher score is found.
		if score > bestScore {
			bestScore = score
			bestChild = child
		}
	}

	// Return the child node with the highest UCB1 score.
	return bestChild
}

// Expand selects an untried move from the node's pool, creates a new child node
// corresponding to the resulting game state, appends it to the node's children list,
// and returns the newly instantiated child node. Protected by mutual exclusion locks.
func (node *Node) Expand() *Node {
	// Guard against nil receiver references.
	if node == nil {
		return nil
	}

	// Acquire the node mutex lock to protect untried moves and children slices from concurrent modification.
	node.mutex.Lock()

	// Ensure the mutex lock is released when the function exits.
	defer node.mutex.Unlock()

	// Check if there are any untried moves remaining to expand.
	if len(node.untriedMoves) == 0 {
		return nil
	}

	// Pop the last untried move from the slice.
	move := node.untriedMoves[len(node.untriedMoves)-1]

	node.untriedMoves = node.untriedMoves[:len(node.untriedMoves)-1]

	// Generate the next game state by applying the popped move.
	nextState := node.state.MakeMove(move)

	// Instantiate a new child node wrapping the new game state.
	childNode := NewNode(nextState, node)

	// Append the new child node to the node's children list.
	node.children = append(node.children, childNode)

	// Return the newly created child node.
	return childNode
}

// SimulateRollout executes a randomized game simulation (playout) from the given
// state until a terminal condition is reached, utilizing an isolated pseudo-random
// number generator seeded with the provided entropy source. Returns the winning player identifier.
func SimulateRollout(state GameState, seed int64) int {
	// Guard against already terminal states at the start of the rollout.
	if state.IsTerminal() {
		return 3 - state.PlayerTurn
	}

	// Initialize a local random number generator with a unique seed to prevent cross-goroutine contention.
	rng := rand.New(rand.NewSource(seed))

	// Create a mutable copy of the game state for the simulation loop.
	current := state

	// Loop until a terminal game state is reached.
	for !current.IsTerminal() {
		// Retrieve all legal moves available in the current simulated state.
		moves := current.GetLegalMoves()

		// Break out of the loop if no legal moves remain.
		if len(moves) == 0 {
			break
		}

		// Select a random move from the available legal choices.
		move := moves[rng.Intn(len(moves))]

		// Apply the move to advance the simulated game state.
		current = current.MakeMove(move)
	}

	// Return the winner by inverting the final player turn.
	return 3 - current.PlayerTurn
}

// Backpropagate traverses upward from the current node to the root, thread-safely
// incrementing visit counters and recording win tallies for nodes whose preceding
// player matches the winning player of the rollout simulation.
func (node *Node) Backpropagate(winner int) {
	// Guard against nil receiver references.
	if node == nil {
		return
	}

	// Validate the winner identifier to ensure it is either player 1 or 2.
	if winner != 1 && winner != 2 {
		return
	}

	// Initialize a pointer to traverse upward from the target node to the root.
	current := node

	// Traverse up the tree ancestry until the root's parent is reached.
	for current != nil {
		// Acquire the mutex lock for the current node during backpropagation updates.
		current.mutex.Lock()
		// Increment the total visit counter for the node.
		current.visits++

		// Determine whether this node counts as a win for the winner:
		// - For the root node (no parent), if the winner matches the root player's turn.
		// - For child nodes, if the winner matches the player who made the move to reach this state.
		var countsAsWin bool

		if current.parent == nil {
			countsAsWin = (current.state.PlayerTurn == winner)
		} else {
			movedPlayer := 3 - current.state.PlayerTurn
			countsAsWin = (movedPlayer == winner)
		}

		// Record a win if the condition is satisfied.
		if countsAsWin {
			current.wins++
		}

		// Release the mutex lock for the current node.
		current.mutex.Unlock()

		// Move up to the parent node in the tree hierarchy.
		current = current.parent
	}
}

// ExecuteParallel coordinates a concurrent Monte Carlo Tree Search across multiple
// worker goroutines using a sync.WaitGroup. Each worker independently performs
// Selection, Expansion, Rollout, and Backpropagation iterations against the shared
// search tree. Once all iterations complete, it selects and returns the optimal move
// based on the robust child rule (maximum visit count).
func ExecuteParallel(state GameState, numberOfIterations int, numberOfWorkers int) int {
	// Guard against invalid iteration counts by falling back to a safe default.
	if numberOfIterations <= 0 {
		numberOfIterations = 1000
	}

	// Guard against invalid worker counts by clamping or defaulting.
	if numberOfWorkers <= 0 {
		numberOfWorkers = 1
	}

	// Ensure worker count does not exceed iteration count.
	if numberOfWorkers > numberOfIterations {
		numberOfWorkers = numberOfIterations
	}

	// Guard against terminal initial game states.
	if state.IsTerminal() {
		return -1
	}

	// Create the root node of the search tree wrapping the initial game state.
	root := NewNode(state, nil)

	// Initialize a wait group to coordinate worker goroutine completion.
	var waitGroup sync.WaitGroup

	// Compute the number of simulation iterations assigned to each worker.
	iterationsPerWorker := numberOfIterations / numberOfWorkers

	// Track active worker goroutines using an atomic integer counter.
	var activeWorkers int64 = int64(numberOfWorkers)

	// Launch multiple worker goroutines to perform MCTS iterations concurrently.
	for workerIndex := 0; workerIndex < numberOfWorkers; workerIndex++ {
		// Increment the wait group counter for each spawned worker.
		waitGroup.Add(1)

		go func(workerID int) {
			// Ensure the wait group counter is decremented when the worker exits.
			defer waitGroup.Done()

			// Execute the assigned share of simulation iterations.
			for index := 0; index < iterationsPerWorker; index++ {
				// Start tree traversal from the root node.
				node := root

				// Selection Phase: Traverse down the tree using UCB1 scores.
				for {
					// Guard against nil node references during traversal.
					if node == nil {
						break
					}

					// Lock the node to safely read state attributes.
					node.mutex.Lock()

					hasUntried := len(node.untriedMoves) > 0
					isTerm := node.state.IsTerminal()

					node.mutex.Unlock()

					// Stop traversing if an expandable node or terminal state is reached.
					if hasUntried || isTerm {
						break
					}

					// Select the best child node using the UCB1 policy.
					best := node.SelectChild(1.41)

					// Break if no valid child can be selected.
					if best == nil {
						break
					}

					// Move down to the selected child node.
					node = best
				}

				// Guard against nil node references before expansion.
				if node == nil {
					continue
				}

				// Expansion Phase: Expand an untried move if the node is non-terminal.
				node.mutex.Lock()

				if !node.state.IsTerminal() && len(node.untriedMoves) > 0 {
					node.mutex.Unlock()
					node = node.Expand()
				} else {
					node.mutex.Unlock()
				}

				// Guard against nil node references before rollout.
				if node == nil {
					continue
				}

				// Simulation Phase: Perform a random rollout simulation from the current node state.
				winner := SimulateRollout(node.state, time.Now().UnixNano()+int64(index*workerID+1))

				// Backpropagation Phase: Propagate the simulation result upward through the tree.
				node.Backpropagate(winner)
			}
			// Atomically decrement the active worker counter upon completion.
			atomic.AddInt64(&activeWorkers, -1)
		}(workerIndex)
	}

	// Block execution until all worker goroutines have finished their work.
	waitGroup.Wait()

	// Guard against nil root references.
	if root == nil {
		return -1
	}

	// Lock the root node to safely analyze final visit counts.
	root.mutex.Lock()

	// Ensure the root node mutex lock is released upon exit.
	defer root.mutex.Unlock()

	// Guard against root nodes with no evaluated children.
	if len(root.children) == 0 {
		legalMoves := root.state.GetLegalMoves()

		if len(legalMoves) > 0 {
			return legalMoves[0]
		}

		return -1
	}

	// Initialize tracking variables for the best move and maximum visits.
	var bestMove int = -1

	maximumVisits := -1

	// Evaluate all children of the root node to find the most visited choice.
	for _, child := range root.children {
		// Guard against nil child pointers in the root's children slice.
		if child == nil {
			continue
		}

		// Lock the child node to safely read its total visit count.
		child.mutex.Lock()

		// Update the best move if this child has a strictly higher visit count.
		if child.visits > maximumVisits {
			maximumVisits = child.visits

			// Calculate the specific move value that resulted in this child state.
			bestMove = root.state.Stones - child.state.Stones
		}

		// Unlock the child node.
		child.mutex.Unlock()
	}

	// Return the optimal move selected by the robust child rule.
	return bestMove
}
