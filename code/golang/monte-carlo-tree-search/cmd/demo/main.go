// -----------------------------------------------------------------------------
// Package: main
// File: main.go
//
// Description:
//
//	This file provides an executable demonstration of the concurrent Monte Carlo
//	Tree Search (MCTS) package by orchestrating a game of Nim between an AI player
//	powered by parallel goroutine workers and an automated opponent. It initializes
//	the game state, loops through turns until a terminal condition is met, and
//	outputs interactive turn updates and final game results.
//
// -----------------------------------------------------------------------------
package main

import (
	"fmt"

	"github.com/braiden-gole/monte-carlo-tree-search/montecarlotreesearch"
)

func main() {
	// Initialize the game state with a starting pile of 15 stones and Player 1's turn.
	game := montecarlotreesearch.GameState{Stones: 15, PlayerTurn: 1}

	// Print the welcome banner and introductory instructions for the Nim game.
	fmt.Println("=== Nim Game: AI (Player 1) vs Opponent (Player 2) ===")
	fmt.Println("Starting pile size: 15 stones. Take 1, 2, or 3 stones per turn. Last to take wins.")

	// Loop continuously until the game reaches a terminal state where no stones remain.
	for !game.IsTerminal() {
		// Display the current remaining stone count and the active player's turn.
		fmt.Printf("\nCurrent Stones: %d | Turn: Player %d\n", game.Stones, game.PlayerTurn)

		// Declare a variable to store the chosen move for the current turn.
		var move int

		// Check if it is the AI's turn (Player 1).
		if game.PlayerTurn == 1 {
			// Notify the user that the AI is calculating its move using parallel MCTS.
			fmt.Println("AI is thinking using Parallel MCTS across multiple goroutines...")

			// Execute parallel MCTS simulations using 30,000 iterations split across 8 workers.
			move = montecarlotreesearch.ExecuteParallel(game, 30000, 8)

			// Print the move selected by the AI agent.
			fmt.Printf("AI chooses to take: %d stones\n", move)
		} else {
			// Assign a simplified automated move of taking 1 stone for the opponent (Player 2).
			move = 1

			// Print the opponent's chosen move.
			fmt.Printf("Opponent takes: %d stones\n", move)
		}

		// Apply the chosen move to advance the game state to the next turn.
		game = game.MakeMove(move)
	}

	// Calculate the winner by inverting the final player turn once the game terminates.
	winner := 3 - game.PlayerTurn

	// Print the final game over message declaring the winning player.
	fmt.Printf("\nGame Over! Player %d wins!\n", winner)
}
