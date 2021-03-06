package game

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/n7down/ssh-chess/internal/logger"
	"golang.org/x/crypto/ssh"

	randomData "github.com/Pallinder/go-randomdata"
)

const (
	gameWidth  = 78
	gameHeight = 22

	keyW = 'w'
	keyA = 'a'
	keyS = 's'
	keyD = 'd'

	keyH = 'h'
	keyJ = 'j'
	keyK = 'k'
	keyL = 'l'

	keyF = 'f'

	keyY = 'y'
	keyN = 'n'

	keyCtrlC = 3
)

type GameManager struct {
	UserCreatedGames map[string]*Game
	Games            map[string]*Game
	HandleChannel    chan ssh.Channel
	logger           logger.Logger
}

func NewGameManager(logger logger.Logger) *GameManager {
	return &GameManager{
		UserCreatedGames: map[string]*Game{},
		Games:            map[string]*Game{},
		HandleChannel:    make(chan ssh.Channel),
		logger:           logger,
	}
}

func (gm *GameManager) getAvailableGame() *Game {
	for _, game := range gm.Games {
		if game.SessionCount() == 1 {
			return game
		}
	}
	return nil
}

func (gm *GameManager) SessionCount() int {
	sum := 0
	for _, game := range gm.UserCreatedGames {
		sum += game.SessionCount()
	}
	for _, game := range gm.Games {
		sum += game.SessionCount()
	}
	return sum
}

func (gm *GameManager) GameCount() int {
	return len(gm.UserCreatedGames) + len(gm.Games)
}

/*
GameManager
- Games Game

Game
- hub

Hub
- Session

Session
- Player
*/

func (gm *GameManager) generateUserCreatedGame() *Game {
	g := NewUserCreatedGame(gameWidth, gameHeight, randomData.SillyName())
	gm.UserCreatedGames[g.Name] = g
	return g
}

func (gm *GameManager) getUserCreatedGame(gameName string) *Game {
	var g *Game

	// check if the UserGame already exists in the map
	if _, ok := gm.UserCreatedGames[gameName]; ok {
		if gm.UserCreatedGames[gameName].SessionCount() == 1 {
			g = gm.UserCreatedGames[gameName]
		} else {
			g = gm.generateUserCreatedGame()
		}
	}

	if g == nil {
		// create the game in UserGames
		g = NewUserCreatedGame(gameWidth, gameHeight, gameName)
		gm.UserCreatedGames[gameName] = g
	}
	return g
}

func (gm *GameManager) getPlayerAndGameName(username string) (string, string) {
	if strings.Contains(username, "#") {
		names := strings.Split(username, "#")
		playerName := names[0]
		gameName := names[1]
		return playerName, gameName
	}
	return username, ""
}

func (gm *GameManager) HandleNewChannel(c ssh.Channel, user string) {

	playerName, gameName := gm.getPlayerAndGameName(user)

	var g *Game
	if gameName != "" {
		gm.logger.Debug(fmt.Sprintf("user game name: %s", gameName))
		g = gm.getUserCreatedGame(gameName)
	}

	if g == nil {
		g = gm.getAvailableGame()
	}

	if g == nil {
		g = NewGame(gameWidth, gameHeight, randomData.SillyName(), gm.logger)
		gm.Games[g.Name] = g
	}
	go g.Run()

	session := NewSession(c, g.WorldWidth(), g.WorldHeight(), playerName, gm.logger)

	g.AddSession(session)

	gm.logger.Print(fmt.Sprintf("player connected: %v", playerName))
	gm.logger.Print(fmt.Sprintf("Player joined. Current stats: %d users, %d games", gm.SessionCount(), gm.GameCount()))

	go func() {
		reader := bufio.NewReader(c)
		for {
			r, _, err := reader.ReadRune()
			gm.logger.Debug(fmt.Sprintf("r: %d", r))
			if err != nil {
				gm.logger.Debug(err.Error())
				break
			}

			// FIXME: create check for arrow keys function
			if r != 0 && r != 27 && r != 91 && r != 65 && r != 66 && r != 67 && r != 68 {
				switch r {
				case keyW, keyK:
					session.Player.HandleUp()
				case keyA, keyH:
					session.Player.HandleLeft()
				case keyS, keyJ:
					session.Player.HandleDown()
				case keyD, keyL:
					session.Player.HandleRight()
				case keyF:
					session.Player.HandleAction()
				case keyCtrlC:
					if g.SessionCount() == 1 {
						if g.userCreatedGame {
							delete(gm.UserCreatedGames, g.Name)
						} else {
							delete(gm.Games, g.Name)
						}
					}

					g.RemoveSession(session, "a test message")
				}
			}
		}
	}()
}
