package game

import (
	"time"

	"github.com/n7down/ssh-chess/internal/logger"
	"golang.org/x/crypto/ssh"
)

type Session struct {
	c ssh.Channel

	LastAction time.Time
	HighScore  int
	Player     *Player
	logger     logger.Logger
}

func NewSession(c ssh.Channel, worldWidth, worldHeight int, playerName string, logger logger.Logger) *Session {

	s := Session{
		c:          c,
		LastAction: time.Now(),
		logger:     logger,
	}
	s.newGame(worldWidth, worldHeight, playerName)

	return &s
}

func (s *Session) newGame(worldWidth, worldHeight int, playerName string) {
	s.Player = NewPlayer(s, worldWidth, worldHeight, playerName, s.logger)
}

func (s *Session) didAction() {
	s.LastAction = time.Now()
}

/*func (s *Session) StartOver(worldWidth, worldHeight int) {*/
//s.newGame(worldWidth, worldHeight, s.Player.Name)
/*}*/

func (s *Session) Read(p []byte) (int, error) {
	return s.c.Read(p)
}

func (s *Session) Write(p []byte) (int, error) {
	return s.c.Write(p)
}
