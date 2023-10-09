package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/inancgumus/screen"
	"github.com/notnil/chess"
)

func printMove(move *chess.Move, game *chess.Game) {
	screen.Clear()
	fmt.Println("Previous move : ", move)
	fmt.Println("======================================")
	fmt.Println(game.Position().Board().Draw())
}

func main() {
	game := chess.NewGame()
	// generate moves until game is over

	round := 0
	hasPrinted := false
	for game.Outcome() == chess.NoOutcome {
		round++
		if round < 10 {
			moves := game.ValidMoves()
			move := moves[rand.Intn(len(moves))]
			game.Move(move)
			printMove(move, game)
		} else {
			if !hasPrinted {
				fmt.Println("Starting to use alpha beta")
				hasPrinted = true
			}
			time.Sleep((time.Second * 1))
			move := alphaBetaSearch(game, 3)
			game.Move(move)
			printMove(move, game)
		}
	}
	// print outcome and game PGN
	fmt.Printf("Game completed. %s by %s.\n", game.Outcome(), game.Method())
	fmt.Println(game.String())
}

func alphaBetaSearch(game *chess.Game, depth int) *chess.Move {
	bestMove := game.ValidMoves()[0]
	bestScore := math.Inf(-1)
	alpha := math.Inf(-1)
	beta := math.Inf(1)

	for _, move := range game.ValidMoves() {
		clone := game.Clone()
		clone.Move(move)
		score := -alphaBeta(game, depth-1, -beta, -alpha)
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
		alpha = math.Max(alpha, bestScore)
		if alpha >= beta {
			break
		}
	}

	return bestMove
}

func alphaBeta(game *chess.Game, depth int, alpha float64, beta float64) float64 {
	if depth == 0 {
		return evaluatePosition(game.Position())
	}

	if game.Position().Turn() == chess.White {
		for _, move := range game.ValidMoves() {
			clone := game.Clone()
			clone.Move(move)
			alpha = math.Max(alpha, alphaBeta(game, depth-1, alpha, beta))
			if beta <= alpha {
				break
			}
		}
		return alpha
	} else {
		for _, move := range game.ValidMoves() {
			clone := game.Clone()
			clone.Move(move)
			beta = math.Min(beta, alphaBeta(game, depth-1, alpha, beta))
			if beta <= alpha {
				break
			}
		}
		return beta
	}
}

func evaluatePosition(position *chess.Position) float64 {
	pieceScores := map[chess.PieceType]float64{
		chess.Pawn:   1,
		chess.Knight: 3,
		chess.Bishop: 3,
		chess.Rook:   5,
		chess.Queen:  9,
		chess.King:   100,
	}

	score := 0.0

	boardStr := [64]string{
		"a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8",
		"b1", "b2", "b3", "b4", "b5", "b6", "b7", "b8",
		"c1", "c2", "c3", "c4", "c5", "c6", "c7", "c8",
		"d1", "d2", "d3", "d4", "d5", "d6", "d7", "d8",
		"e1", "e2", "e3", "e4", "e5", "e6", "e7", "e8",
		"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8",
		"g1", "g2", "g3", "g4", "g5", "g6", "g7", "g8",
		"h1", "h2", "h3", "h4", "h5", "h6", "h7", "h8",
	}

	// Evaluate material balance
	for i, _ := range boardStr {
		piece := position.Board().Piece(chess.Square(i))
		if piece.Color() == position.Turn() {
			score += pieceScores[piece.Type()]
		} else {
			score -= pieceScores[piece.Type()]
		}
	}

	return score
}
