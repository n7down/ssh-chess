package main

import (
	"time"

	"golang.org/x/crypto/ssh"
)

type Session struct {
	c ssh.Channel

	LastAction time.Time
	HighScore  int
	Player     *Player
}

func NewSession(c ssh.Channel, worldWidth, worldHeight int, playerName string) *Session {

	s := Session{c: c, LastAction: time.Now()}
	s.newGame(worldWidth, worldHeight, playerName)

	return &s
}

func (s *Session) newGame(worldWidth, worldHeight int, playerName string) {
	s.Player = NewPlayer(s, worldWidth, worldHeight, playerName)
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
