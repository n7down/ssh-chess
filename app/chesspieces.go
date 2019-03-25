package main

type ChessPiecesColor int

const (
	White ChessPiecesColor = iota
	Black
)

const (
	whitePawn   = "♟"
	whiteRook   = "♜"
	whiteKnight = "♞"
	whiteBishop = "♝"
	whiteKing   = "♚"
	whiteQueen  = "♛"

	blackPawn   = "♙"
	blackRook   = "♖"
	blackKnight = "♘"
	blackBishop = "♗"
	blackKing   = "♔"
	blackQueen  = "♕"
)

func (c ChessPiecesColor) String() string {
	return [...]string{whiteKing, blackKing}[c]
}
