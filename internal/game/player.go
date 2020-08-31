package game

import (
	"fmt"
	"time"

	"github.com/cznic/mathutil"
	"github.com/n7down/ssh-chess/internal/logger"
	"github.com/notnil/chess"
)

type PlayerState int

const (
	SelectingPiece PlayerState = iota
	PlacingPiece
)

type KeyState int

const (
	KeyUp KeyState = iota
	KeyDown
	KeyLeft
	KeyRight
	KeyAction
	KeyNone
)

type Player struct {
	s                     *Session
	Name                  string
	CreatedAt             time.Time
	BoardPosition         *Position
	IsActive              bool
	PlayerColor           ChessPiecesColor
	PlayerState           PlayerState
	SelectedPiecePosition *Position
	currentKeyState       KeyState
	previousKeyState      KeyState
	TakenPiecesList       []string
	logger                logger.Logger
}

func NewPlayer(s *Session, worldWidth, worldHeight int, playerName string, logger logger.Logger) *Player {
	isActive := false

	player := &Player{
		s:                     s,
		Name:                  playerName,
		CreatedAt:             time.Now(),
		BoardPosition:         &Position{0, 0},
		IsActive:              isActive,
		PlayerColor:           White,
		PlayerState:           SelectingPiece,
		SelectedPiecePosition: &Position{-1, -1},
		currentKeyState:       KeyNone,
		previousKeyState:      KeyNone,
		TakenPiecesList:       []string{},
		logger:                logger,
	}

	// FIXME: add player to database if it doesnt exist
	return player
}

func (p *Player) SetIsActive(b bool) {
	p.IsActive = b
	if b {
		p.PlayerColor = Black
		p.BoardPosition = &Position{0, 7}
	} else {
		p.PlayerColor = White
		p.BoardPosition = &Position{0, 0}
	}
}

func (p *Player) positionInList(validPositions []Position) bool {
	for _, pp := range validPositions {
		if p.BoardPosition.x == pp.x && p.BoardPosition.y == pp.y {
			return true
		}
	}
	return false
}

func (p *Player) getVaildPositionsForSelectedPiece(validMoves []*chess.Move) []Position {
	positions := []Position{}

	// get the model position for the selected piece
	selectedPieceModel := p.SelectedPiecePosition.positionToModel()

	// add selected piece position to the valid positions
	positions = append(positions, *p.SelectedPiecePosition)

	for _, move := range validMoves {
		moveString := move.String()

		if moveString[0:2] == selectedPieceModel {

			// get the moveString[2:4] convert it to the position
			validPosition := modelToPosition(moveString[2:4])

			// append the new position to the positions
			positions = append(positions, validPosition)
		}
	}
	return positions
}

func (p *Player) HandleUp() {
	p.currentKeyState = KeyUp
}

func (p *Player) HandleLeft() {
	p.currentKeyState = KeyLeft
}

func (p *Player) HandleDown() {
	p.currentKeyState = KeyDown
}

func (p *Player) HandleRight() {
	p.currentKeyState = KeyRight
}

func (p *Player) HandleAction() {
	p.currentKeyState = KeyAction
}

func (p *Player) canTakePiece(pieceToTake string) bool {
	var canTakePiece bool = false
	switch p.PlayerColor {
	case White:
		if pieceToTake == blackKing || pieceToTake == blackQueen ||
			pieceToTake == blackRook || pieceToTake == blackBishop ||
			pieceToTake == blackKnight || pieceToTake == blackPawn {
			canTakePiece = true
		}
	case Black:
		if pieceToTake == whiteKing || pieceToTake == whiteQueen ||
			pieceToTake == whiteRook || pieceToTake == whiteBishop ||
			pieceToTake == whiteKnight || pieceToTake == whitePawn {
			canTakePiece = true
		}
	}
	return canTakePiece
}

func (p *Player) canMovePiece(pieceToMove string) bool {
	var canMovePiece bool = false
	switch p.PlayerColor {
	case Black:
		if pieceToMove == blackKing || pieceToMove == blackQueen ||
			pieceToMove == blackRook || pieceToMove == blackBishop ||
			pieceToMove == blackKnight || pieceToMove == blackPawn {
			canMovePiece = true
		}
	case White:
		if pieceToMove == whiteKing || pieceToMove == whiteQueen ||
			pieceToMove == whiteRook || pieceToMove == whiteBishop ||
			pieceToMove == whiteKnight || pieceToMove == whitePawn {
			canMovePiece = true
		}
	}
	return canMovePiece
}

func (p *Player) Update(g *Game, delta float64) {
	if p.previousKeyState == p.currentKeyState {
		p.currentKeyState = KeyNone
	}

	switch p.currentKeyState {
	case KeyUp:
		if p.IsActive {
			p.BoardPosition.y--
			p.BoardPosition.x, p.BoardPosition.y = mathutil.Clamp(p.BoardPosition.x, 0, 7), mathutil.Clamp(p.BoardPosition.y, 0, 7)
			p.logger.Debug(fmt.Sprintf("x: %d y: %d", p.BoardPosition.x, p.BoardPosition.y))
		}

	case KeyDown:
		if p.IsActive {
			p.BoardPosition.y++
			p.BoardPosition.x, p.BoardPosition.y = mathutil.Clamp(p.BoardPosition.x, 0, 7), mathutil.Clamp(p.BoardPosition.y, 0, 7)
			p.logger.Debug(fmt.Sprintf("x: %d y: %d", p.BoardPosition.x, p.BoardPosition.y))
		}

	case KeyRight:
		if p.IsActive {
			p.BoardPosition.x++
			p.BoardPosition.x, p.BoardPosition.y = mathutil.Clamp(p.BoardPosition.x, 0, 7), mathutil.Clamp(p.BoardPosition.y, 0, 7)
			p.logger.Debug(fmt.Sprintf("x: %d y: %d", p.BoardPosition.x, p.BoardPosition.y))
		}

	case KeyLeft:
		if p.IsActive {
			p.BoardPosition.x--
			p.BoardPosition.x, p.BoardPosition.y = mathutil.Clamp(p.BoardPosition.x, 0, 7), mathutil.Clamp(p.BoardPosition.y, 0, 7)
			p.logger.Debug(fmt.Sprintf("x: %d y: %d", p.BoardPosition.x, p.BoardPosition.y))
		}

	case KeyAction:

		if p.IsActive && p.PlayerState == SelectingPiece {
			p.SelectedPiecePosition.x, p.SelectedPiecePosition.y = p.BoardPosition.x, p.BoardPosition.y
			p.logger.Debug(fmt.Sprintf("selected piece: %v  x: %d y: %d",
				g.board[p.SelectedPiecePosition.x][p.SelectedPiecePosition.y],
				p.SelectedPiecePosition.x,
				p.SelectedPiecePosition.y))

			// FIXME: check that a player can select the position on the board
			pieceToMove := g.board[p.SelectedPiecePosition.x][p.SelectedPiecePosition.y]
			if p.canMovePiece(pieceToMove) {
				p.PlayerState = PlacingPiece
				p.logger.Debug("piece selected - in placing piece state")

				// display the valid moves
				validMoves := g.Model.ValidMoves()
				validPositions := p.getVaildPositionsForSelectedPiece(validMoves)
				p.logger.Debug(fmt.Sprintf("valid positions: %v", validPositions))
			}

		} else if p.IsActive && p.SelectedPiecePosition.x == p.BoardPosition.x &&
			p.SelectedPiecePosition.y == p.BoardPosition.y &&
			p.PlayerState == PlacingPiece {

			p.SelectedPiecePosition = &Position{-1, -1}
			p.PlayerState = SelectingPiece
			p.logger.Debug("putting piece back - in selecting piece state")

		} else if p.IsActive && p.PlayerState == PlacingPiece {
			validMoves := g.Model.ValidMoves()
			validPositions := p.getVaildPositionsForSelectedPiece(validMoves)
			positionIsValid := p.positionInList(validPositions)
			if positionIsValid {

				p.logger.Debug(fmt.Sprintf("selected piece x: %d y: %d", p.SelectedPiecePosition.x, p.SelectedPiecePosition.y))

				p.logger.Debug(fmt.Sprintf("moving piece x: %d y: %d to x: %d y: %d - in selecting piece state",
					p.SelectedPiecePosition.x,
					p.SelectedPiecePosition.y,
					p.BoardPosition.x,
					p.BoardPosition.y))

				pieceToTake := g.board[p.BoardPosition.x][p.BoardPosition.y]
				if p.canTakePiece(pieceToTake) {
					p.logger.Debug(fmt.Sprintf("taking piece: %s", pieceToTake))
					p.TakenPiecesList = append(p.TakenPiecesList, pieceToTake)
					p.logger.Debug(fmt.Sprintf("taken list: %v", p.TakenPiecesList))
				}

				selectedPiece := g.board[p.SelectedPiecePosition.x][p.SelectedPiecePosition.y]
				g.board[p.BoardPosition.x][p.BoardPosition.y] = selectedPiece
				g.board[p.SelectedPiecePosition.x][p.SelectedPiecePosition.y] = " "
				p.PlayerState = SelectingPiece

				//  move pieces in the model
				selectedPieceModel := p.SelectedPiecePosition.positionToModel()
				boardPositionModel := p.BoardPosition.positionToModel()
				moveStr := selectedPieceModel + boardPositionModel
				p.logger.Debug(fmt.Sprintf("move string: " + moveStr))
				g.Model.MoveStr(moveStr)

				p.logger.Debug(g.Model.Position().Board().Draw())

				g.CheckGameState()
				g.SwitchPlayersIsActive()
			}
		}

	default:
	}

	if p.IsActive && p.PlayerState == SelectingPiece {
		g.SetBoardColorsSelectingPiece(Position{p.BoardPosition.x, p.BoardPosition.y}, p.PlayerColor)

	} else if p.IsActive && p.PlayerState == PlacingPiece {
		validMoves := g.Model.ValidMoves()
		validPositions := p.getVaildPositionsForSelectedPiece(validMoves)
		positionIsValid := p.positionInList(validPositions)

		if positionIsValid {
			g.SetPositionColor(*p.BoardPosition, Green)
		} else {
			g.SetPositionColor(*p.BoardPosition, Red)
		}
	}

	p.previousKeyState = p.currentKeyState
}
