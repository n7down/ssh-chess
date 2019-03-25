package main

import (
	"fmt"
	"github.com/notnil/chess"
)

func main() {

	// start a new game
	game := chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}))

	// check the valid moves
	moves := game.ValidMoves()
	fmt.Println(moves)

	// make a move
	game.MoveStr("b1a3")

	// switch turns
	//turn := game.Position().Turn()
	//fmt.Println(turn)

	// get the valid moves
	moves = game.ValidMoves()
	fmt.Println(moves)

	// make a move
	game.MoveStr("b8a6")

	moves = game.ValidMoves()
	fmt.Println(moves)
	game.MoveStr("d2d3")

	fmt.Println(game.Position().Board().Draw())
	fmt.Println(game)
}
