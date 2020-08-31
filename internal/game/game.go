package game

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/n7down/ssh-chess/internal/logger"

	aurora "github.com/logrusorgru/aurora"
	chess "github.com/notnil/chess"
	uuid "github.com/satori/go.uuid"
)

type BoardColor int

const (
	Red BoardColor = iota
	Green
	None
)

const (
	verticalWall   = '║'
	horizontalWall = '═'
	topLeft        = '╔'
	topRight       = '╗'
	bottomRight    = '╝'
	bottomLeft     = '╚'
	blank          = ' '
)

type Game struct {
	userCreatedGame bool
	Name            string
	Redraw          chan struct{}
	level           [][]string
	hub             Hub
	board           [][]string
	started         bool
	boardColors     map[Position]BoardColor
	mutex           sync.RWMutex
	Model           *chess.Game
	startTime       time.Time
	id              string
	logger          logger.Logger
}

func NewGame(worldWidth, worldHeight int, name string, logger logger.Logger) *Game {
	g := &Game{
		userCreatedGame: false,
		Name:            name,
		Redraw:          make(chan struct{}),
		hub:             NewHub(),
		Model:           chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{})),
		logger:          logger,
	}

	id, err := uuid.NewV4()
	if err != nil {
		g.logger.Debug(fmt.Sprintf("error generating uuid: %v", err.Error()))
	}
	g.id = id.String()

	g.mutex = sync.RWMutex{}

	g.started = false

	g.initializeColors()
	g.initializeBoard()
	g.initializeLevel(worldWidth, worldHeight)
	g.SetBoardColorsSelectingPiece(Position{0, 0}, White)
	g.drawBoard(worldWidth, worldHeight)

	return g
}

func NewUserCreatedGame(worldWidth, worldHeight int, name string) *Game {
	g := &Game{
		userCreatedGame: true,
		Name:            name,
		Redraw:          make(chan struct{}),
		hub:             NewHub(),
		Model:           chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{})),
	}

	id, err := uuid.NewV4()
	if err != nil {
		g.logger.Debug(fmt.Sprintf("error generating uuid: %v", err.Error()))
	}
	g.id = id.String()

	g.mutex = sync.RWMutex{}

	g.started = false

	g.initializeColors()
	g.initializeBoard()
	g.initializeLevel(worldWidth, worldHeight)
	g.SetBoardColorsSelectingPiece(Position{0, 0}, White)
	g.drawBoard(worldWidth, worldHeight)

	return g
}

func (g *Game) getColor(p Position, s string) string {
	g.mutex.RLock()
	c := g.boardColors[p]
	g.mutex.RUnlock()
	switch c {
	case Red:
		return aurora.Sprintf(aurora.Red(s))
	case Green:
		return aurora.Sprintf(aurora.Green(s))
	}
	return s
}

func (g *Game) SetBoardColorsSelectingPiece(playerPosition Position, chessPiecesColor ChessPiecesColor) {
	charOver := g.board[playerPosition.x][playerPosition.y]

	g.resetBoardColors()

	if chessPiecesColor == White {
		if charOver == whitePawn || charOver == whiteRook ||
			charOver == whiteKnight || charOver == whiteBishop ||
			charOver == whiteKing || charOver == whiteQueen {

			g.mutex.Lock()
			g.boardColors[Position{playerPosition.x, playerPosition.y}] = Green
			g.boardColors[Position{playerPosition.x + 1, playerPosition.y}] = Green
			g.boardColors[Position{playerPosition.x, playerPosition.y + 1}] = Green
			g.boardColors[Position{playerPosition.x + 1, playerPosition.y + 1}] = Green
			g.mutex.Unlock()
		} else {
			g.mutex.Lock()
			g.boardColors[Position{playerPosition.x, playerPosition.y}] = Red
			g.boardColors[Position{playerPosition.x + 1, playerPosition.y}] = Red
			g.boardColors[Position{playerPosition.x, playerPosition.y + 1}] = Red
			g.boardColors[Position{playerPosition.x + 1, playerPosition.y + 1}] = Red
			g.mutex.Unlock()
		}
	} else {
		if charOver == blackPawn || charOver == blackRook ||
			charOver == blackKnight || charOver == blackBishop ||
			charOver == blackKing || charOver == blackQueen {

			g.mutex.Lock()
			g.boardColors[Position{playerPosition.x, playerPosition.y}] = Green
			g.boardColors[Position{playerPosition.x + 1, playerPosition.y}] = Green
			g.boardColors[Position{playerPosition.x, playerPosition.y + 1}] = Green
			g.boardColors[Position{playerPosition.x + 1, playerPosition.y + 1}] = Green
			g.mutex.Unlock()
		} else {
			g.mutex.Lock()
			g.boardColors[Position{playerPosition.x, playerPosition.y}] = Red
			g.boardColors[Position{playerPosition.x + 1, playerPosition.y}] = Red
			g.boardColors[Position{playerPosition.x, playerPosition.y + 1}] = Red
			g.boardColors[Position{playerPosition.x + 1, playerPosition.y + 1}] = Red
			g.mutex.Unlock()
		}
	}
}

func (g *Game) SetPositionColor(playerPosition Position, boardColor BoardColor) {
	g.resetBoardColors()

	g.mutex.Lock()
	g.boardColors[Position{playerPosition.x, playerPosition.y}] = boardColor
	g.boardColors[Position{playerPosition.x + 1, playerPosition.y}] = boardColor
	g.boardColors[Position{playerPosition.x, playerPosition.y + 1}] = boardColor
	g.boardColors[Position{playerPosition.x + 1, playerPosition.y + 1}] = boardColor
	g.mutex.Unlock()
}

func (g *Game) resetBoardColors() {
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			g.mutex.Lock()
			g.boardColors[Position{x, y}] = None
			g.mutex.Unlock()
		}
	}
}

func (g *Game) initializeColors() {
	g.boardColors = make(map[Position]BoardColor)
	g.resetBoardColors()
}

func (g *Game) initializeLevel(width int, height int) {
	g.level = make([][]string, width)
	for x := range g.level {
		g.level[x] = make([]string, height)
	}
}

func (g *Game) initializeBoard() {
	// this is the board that the players change
	g.board = [][]string{
		{whiteRook, whitePawn, " ", " ", " ", " ", blackPawn, blackRook},
		{whiteKnight, whitePawn, " ", " ", " ", " ", blackPawn, blackKnight},
		{whiteBishop, whitePawn, " ", " ", " ", " ", blackPawn, blackBishop},
		{whiteQueen, whitePawn, " ", " ", " ", " ", blackPawn, blackQueen},
		{whiteKing, whitePawn, " ", " ", " ", " ", blackPawn, blackKing},
		{whiteBishop, whitePawn, " ", " ", " ", " ", blackPawn, blackBishop},
		{whiteKnight, whitePawn, " ", " ", " ", " ", blackPawn, blackKnight},
		{whiteRook, whitePawn, " ", " ", " ", " ", blackPawn, blackRook},
	}
}

func (g *Game) drawBoard(width, height int) {
	renderedBoard := [][]string{

		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", "8", " ", "7", " ", "6", " ", "5", " ", "4", " ", "3", " ", "2", " ", "1", " "},
		{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
		{
			" ",
			g.getColor(Position{0, 0}, "+"),
			"|",
			g.getColor(Position{0, 1}, "+"),
			"|",
			g.getColor(Position{0, 2}, "+"),
			"|",
			g.getColor(Position{0, 3}, "+"),
			"|",
			g.getColor(Position{0, 4}, "+"),
			"|",
			g.getColor(Position{0, 5}, "+"),
			"|",
			g.getColor(Position{0, 6}, "+"),
			"|",
			g.getColor(Position{0, 7}, "+"),
			"|",
			g.getColor(Position{0, 8}, "+"),
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			"A",
			"-",
			g.board[0][0],
			"-",
			g.board[0][1],
			"-",
			g.board[0][2],
			"-",
			g.board[0][3],
			"-",
			g.board[0][4],
			"-",
			g.board[0][5],
			"-",
			g.board[0][6],
			"-",
			g.board[0][7],
			"-",
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			" ",
			g.getColor(Position{1, 0}, "+"),
			"|",
			g.getColor(Position{1, 1}, "+"),
			"|",
			g.getColor(Position{1, 2}, "+"),
			"|",
			g.getColor(Position{1, 3}, "+"),
			"|",
			g.getColor(Position{1, 4}, "+"),
			"|",
			g.getColor(Position{1, 5}, "+"),
			"|",
			g.getColor(Position{1, 6}, "+"),
			"|",
			g.getColor(Position{1, 7}, "+"),
			"|",
			g.getColor(Position{1, 8}, "+"),
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			"B",
			"-",
			g.board[1][0],
			"-",
			g.board[1][1],
			"-",
			g.board[1][2],
			"-",
			g.board[1][3],
			"-",
			g.board[1][4],
			"-",
			g.board[1][5],
			"-",
			g.board[1][6],
			"-",
			g.board[1][7],
			"-",
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			" ",
			g.getColor(Position{2, 0}, "+"),
			"|",
			g.getColor(Position{2, 1}, "+"),
			"|",
			g.getColor(Position{2, 2}, "+"),
			"|",
			g.getColor(Position{2, 3}, "+"),
			"|",
			g.getColor(Position{2, 4}, "+"),
			"|",
			g.getColor(Position{2, 5}, "+"),
			"|",
			g.getColor(Position{2, 6}, "+"),
			"|",
			g.getColor(Position{2, 7}, "+"),
			"|",
			g.getColor(Position{2, 8}, "+"),
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			"C",
			"-",
			g.board[2][0],
			"-",
			g.board[2][1],
			"-",
			g.board[2][2],
			"-",
			g.board[2][3],
			"-",
			g.board[2][4],
			"-",
			g.board[2][5],
			"-",
			g.board[2][6],
			"-",
			g.board[2][7],
			"-",
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			" ",
			g.getColor(Position{3, 0}, "+"),
			"|",
			g.getColor(Position{3, 1}, "+"),
			"|",
			g.getColor(Position{3, 2}, "+"),
			"|",
			g.getColor(Position{3, 3}, "+"),
			"|",
			g.getColor(Position{3, 4}, "+"),
			"|",
			g.getColor(Position{3, 5}, "+"),
			"|",
			g.getColor(Position{3, 6}, "+"),
			"|",
			g.getColor(Position{3, 7}, "+"),
			"|",
			g.getColor(Position{3, 8}, "+"),
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			"D",
			"-",
			g.board[3][0],
			"-",
			g.board[3][1],
			"-",
			g.board[3][2],
			"-",
			g.board[3][3],
			"-",
			g.board[3][4],
			"-",
			g.board[3][5],
			"-",
			g.board[3][6],
			"-",
			g.board[3][7],
			"-",
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			" ",
			g.getColor(Position{4, 0}, "+"),
			"|",
			g.getColor(Position{4, 1}, "+"),
			"|",
			g.getColor(Position{4, 2}, "+"),
			"|",
			g.getColor(Position{4, 3}, "+"),
			"|",
			g.getColor(Position{4, 4}, "+"),
			"|",
			g.getColor(Position{4, 5}, "+"),
			"|",
			g.getColor(Position{4, 6}, "+"),
			"|",
			g.getColor(Position{4, 7}, "+"),
			"|",
			g.getColor(Position{4, 8}, "+"),
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			"E",
			"-",
			g.board[4][0],
			"-",
			g.board[4][1],
			"-",
			g.board[4][2],
			"-",
			g.board[4][3],
			"-",
			g.board[4][4],
			"-",
			g.board[4][5],
			"-",
			g.board[4][6],
			"-",
			g.board[4][7],
			"-",
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			" ",
			g.getColor(Position{5, 0}, "+"),
			"|",
			g.getColor(Position{5, 1}, "+"),
			"|",
			g.getColor(Position{5, 2}, "+"),
			"|",
			g.getColor(Position{5, 3}, "+"),
			"|",
			g.getColor(Position{5, 4}, "+"),
			"|",
			g.getColor(Position{5, 5}, "+"),
			"|",
			g.getColor(Position{5, 6}, "+"),
			"|",
			g.getColor(Position{5, 7}, "+"),
			"|",
			g.getColor(Position{5, 8}, "+"),
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			"F",
			"-",
			g.board[5][0],
			"-",
			g.board[5][1],
			"-",
			g.board[5][2],
			"-",
			g.board[5][3],
			"-",
			g.board[5][4],
			"-",
			g.board[5][5],
			"-",
			g.board[5][6],
			"-",
			g.board[5][7],
			"-",
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			" ",
			g.getColor(Position{6, 0}, "+"),
			"|",
			g.getColor(Position{6, 1}, "+"),
			"|",
			g.getColor(Position{6, 2}, "+"),
			"|",
			g.getColor(Position{6, 3}, "+"),
			"|",
			g.getColor(Position{6, 4}, "+"),
			"|",
			g.getColor(Position{6, 5}, "+"),
			"|",
			g.getColor(Position{6, 6}, "+"),
			"|",
			g.getColor(Position{6, 7}, "+"),
			"|",
			g.getColor(Position{6, 8}, "+"),
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			"G",
			"-",
			g.board[6][0],
			"-",
			g.board[6][1],
			"-",
			g.board[6][2],
			"-",
			g.board[6][3],
			"-",
			g.board[6][4],
			"-",
			g.board[6][5],
			"-",
			g.board[6][6],
			"-",
			g.board[6][7],
			"-",
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			" ",
			g.getColor(Position{7, 0}, "+"),
			"|",
			g.getColor(Position{7, 1}, "+"),
			"|",
			g.getColor(Position{7, 2}, "+"),
			"|",
			g.getColor(Position{7, 3}, "+"),
			"|",
			g.getColor(Position{7, 4}, "+"),
			"|",
			g.getColor(Position{7, 5}, "+"),
			"|",
			g.getColor(Position{7, 6}, "+"),
			"|",
			g.getColor(Position{7, 7}, "+"),
			"|",
			g.getColor(Position{7, 8}, "+"),
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			"H",
			"-",
			g.board[7][0],
			"-",
			g.board[7][1],
			"-",
			g.board[7][2],
			"-",
			g.board[7][3],
			"-",
			g.board[7][4],
			"-",
			g.board[7][5],
			"-",
			g.board[7][6],
			"-",
			g.board[7][7],
			"-",
		},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{" ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-", " ", "-"},
		{
			" ",
			g.getColor(Position{8, 0}, "+"),
			"|",
			g.getColor(Position{8, 1}, "+"),
			"|",
			g.getColor(Position{8, 2}, "+"),
			"|",
			g.getColor(Position{8, 3}, "+"),
			"|",
			g.getColor(Position{8, 4}, "+"),
			"|",
			g.getColor(Position{8, 5}, "+"),
			"|",
			g.getColor(Position{8, 6}, "+"),
			"|",
			g.getColor(Position{8, 7}, "+"),
			"|",
			g.getColor(Position{8, 8}, "+"),
		},
	}

	for x := range g.level {
		for y := range g.level[x] {
			lenX, lenY := len(renderedBoard), len(renderedBoard[0])
			if x < lenX && y < lenY {
				g.level[x][y] = renderedBoard[x][y]
			} else {
				g.level[x][y] = string(blank)
			}
		}
	}
}

func (g *Game) players() map[*Player]*Session {
	players := make(map[*Player]*Session)

	for session := range g.hub.Sessions {
		players[session.Player] = session
	}

	return players
}

func (g *Game) CheckGameState() {

	g.logger.Debug("checking game state")

	var err error
	if g.Model.Outcome() != chess.NoOutcome {

		outcome := g.Model.Outcome().String()
		gameMessage := fmt.Sprintf("game is over. %s by %s\ngame string: %s", outcome, g.Model.Method(), g.Model.String())
		for s := range g.hub.Sessions {
			unregisterMessage := UnregisterMessage{
				session: s,
				message: gameMessage,
			}
			g.hub.Unregister <- unregisterMessage
		}
	}

	if err != nil {
		g.logger.Debug(fmt.Sprintf("error sending data: %v", err.Error()))
	}
}

func (g *Game) SwitchPlayersIsActive() {
	if len(g.players()) > 1 {
		for player := range g.players() {
			isActive := player.IsActive
			if isActive {
				player.IsActive = false
			} else {
				player.IsActive = true
				activePlayerBoardPosition := player.BoardPosition
				g.SetBoardColorsSelectingPiece(
					Position{
						activePlayerBoardPosition.x,
						activePlayerBoardPosition.y,
					},
					player.PlayerColor)
			}
		}
	}
}

func (g *Game) roomString(s *Session) string {
	worldWidth := len(g.level)
	worldHeight := len(g.level[0])

	// Create two dimensional slice of strings to represent the world. It's two
	// characters larger in each direction to accomodate for walls.
	strWorld := make([][]string, worldWidth+2)
	for x := range strWorld {
		strWorld[x] = make([]string, worldHeight+2)
	}

	// Load the walls into the rune slice
	for x := 0; x < worldWidth+2; x++ {
		strWorld[x][0] = string(blank)
		strWorld[x][worldHeight+1] = string(blank)
	}
	for y := 0; y < worldHeight+2; y++ {
		strWorld[0][y] = string(blank)
		strWorld[worldWidth+1][y] = string(blank)
	}

	// Time for the edges!
	strWorld[0][0] = string(blank)
	strWorld[worldWidth+1][0] = string(blank)
	strWorld[worldWidth+1][worldHeight+1] = string(blank)
	strWorld[0][worldHeight+1] = string(blank)

	// draw the board
	g.drawBoard(worldWidth, worldHeight)

	// Load the level into the string slice
	for x := 0; x < worldWidth; x++ {
		for y := 0; y < worldHeight; y++ {
			tile := g.level[x][y]
			strWorld[x+1][y+1] = tile
		}
	}

	// draw players taken pieces
	playersTakenPieces := s.Player.TakenPiecesList
	defaultY := 5
	x := 10
	y := defaultY
	for i, p := range playersTakenPieces {
		if i%8 == 0 {
			x = x - 2
			y = defaultY
		}
		y = y + 1
		strWorld[x][y] = string(p)
	}

	// TODO: show if a piece is being placed
	// Draw the player's name
	playerChessPiecesColor := s.Player.PlayerColor
	playerIsActive := s.Player.IsActive
	playerState := s.Player.PlayerState
	playerName := s.Player.Name

	var playerNameToDisplay string
	if playerState == PlacingPiece {
		boardCoords := s.Player.SelectedPiecePosition.positionToModel()
		if playerIsActive {
			playerNameToDisplay = fmt.Sprintf(" [ %s%s %s ] ", playerChessPiecesColor, playerName, boardCoords)
		} else {
			playerNameToDisplay = fmt.Sprintf(" %s%s %s", playerChessPiecesColor, playerName, boardCoords)
		}
	} else {
		if playerIsActive {
			playerNameToDisplay = fmt.Sprintf(" [ %s%s ] ", playerChessPiecesColor, playerName)
		} else {
			playerNameToDisplay = fmt.Sprintf(" %s%s ", playerChessPiecesColor, playerName)
		}
	}

	for i, r := range playerNameToDisplay {
		strWorld[3+i][worldHeight+1] = string(r)
	}

	// Draw opponents name to the left of the players name
	if len(g.players()) > 1 {
		for player := range g.players() {
			if player == s.Player {
				continue
			}

			opponentName := player.Name
			opponentChessPiecesColor := player.PlayerColor
			opponentIsActive := player.IsActive
			opponentPlayerState := player.PlayerState

			var opponent string
			if opponentPlayerState == PlacingPiece {
				boardCoords := player.SelectedPiecePosition.positionToModel()
				if opponentIsActive {
					opponent = fmt.Sprintf(" [ %s%s %s ] ", opponentChessPiecesColor, opponentName, boardCoords)
				} else {
					opponent = fmt.Sprintf(" %s%s %s", opponentChessPiecesColor, opponentName, boardCoords)
				}
			} else {
				if opponentIsActive {
					opponent = fmt.Sprintf(" [ %s%s ] ", opponentChessPiecesColor, opponentName)
				} else {
					opponent = fmt.Sprintf(" %s%s ", opponentChessPiecesColor, opponentName)
				}
			}
			for i, r := range opponent {
				charsRemaining := len(opponent) - i
				strWorld[len(strWorld)-3-charsRemaining][len(strWorld[0])-1] = string(r)
			}
		}
	}

	// draw opponents taken pieces
	if len(g.players()) > 1 {
		for player := range g.players() {
			if player == s.Player {
				continue
			}
			opponentsTakenPieces := player.TakenPiecesList
			defaultY := 5
			x := 70
			y := defaultY
			for i, p := range opponentsTakenPieces {
				if i%8 == 0 {
					x = x - 2
					y = defaultY
				}
				y = y + 1
				strWorld[x][y] = string(p)
			}
		}
	}

	// Draw the game's name
	nameStr := fmt.Sprintf(" %s ", g.Name)
	for i, r := range nameStr {
		strWorld[3+i][0] = string(r)
	}

	// Convert the rune slice to a string
	buffer := bytes.NewBuffer(make([]byte, 0, worldWidth*worldHeight*2))
	for y := 0; y < len(strWorld[0]); y++ {
		for x := 0; x < len(strWorld); x++ {
			buffer.WriteString(strWorld[x][y])
		}

		// Don't add an extra newline if we're on the last iteration
		if y != len(strWorld[0])-1 {
			buffer.WriteString("\r\n")
		}
	}

	return buffer.String()

}

func (g *Game) WorldWidth() int {
	return len(g.level)
}

func (g *Game) WorldHeight() int {
	return len(g.level[0])
}

func (g *Game) SessionCount() int {
	return len(g.hub.Sessions)
}

func (g *Game) startGame() {
	rand.Seed(time.Now().UnixNano())
	var randomBool bool
	randomBool = rand.Float32() < 0.5

	for player := range g.players() {
		//fmt.Println(fmt.Sprintf("random bool: %v", randomBool))
		g.logger.Debug(fmt.Sprintf("random bool: %v", randomBool))
		player.SetIsActive(randomBool)
		randomBool = !randomBool
	}

	g.startTime = time.Now()
}

func (g *Game) Run() {

	// Proxy g.Redraw's channel to g.hub.Redraw
	go func() {
		for {
			g.hub.Redraw <- <-g.Redraw
		}
	}()

	// Run game loop
	go func() {
		var lastUpdate time.Time

		c := time.Tick(time.Second / 60)
		for now := range c {
			g.Update(float64(now.Sub(lastUpdate)) / float64(time.Millisecond))

			lastUpdate = now
		}
	}()

	// Redraw regularly.
	//
	// TODO: Implement diffing and only redraw when needed
	go func() {
		c := time.Tick(time.Second / 10)
		for range c {
			g.Redraw <- struct{}{}

			if g.started == false && len(g.players()) > 1 {
				g.logger.Debug("starting game")

				// start the game
				g.startGame()
				g.started = true
			}
		}
	}()

	g.hub.Run(g)
}

// Update is the main game logic loop. Delta is the time since the last update
// in milliseconds.
func (g *Game) Update(delta float64) {

	// Update player data
	for player, _ := range g.players() {
		player.Update(g, delta)
	}
}

func (g *Game) Render(s *Session) {
	worldStr := g.roomString(s)

	var b bytes.Buffer
	b.WriteString("\033[H\033[2J")
	b.WriteString(worldStr)

	// Send over the rendered world
	io.Copy(s, &b)
}

func (g *Game) AddSession(s *Session) {
	g.hub.Register <- s
}

func (g *Game) RemoveSession(s *Session, msg string) {
	message := "\r\n\r\n" + msg + "\r\n\r\n"
	u := UnregisterMessage{
		session: s,
		message: message,
	}
	g.hub.Unregister <- u
}
