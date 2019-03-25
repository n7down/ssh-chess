package main

import (
	data "github.com/n7down/ssh-chess/services/db"
	"golang.org/x/crypto/bcrypt"
)

type Player struct {
	UserName string `json:"username" binding:"required"`
	Secret   string `json:"secret" binding:"required"`
}

type CheckSecretRequest struct {
	UserName string `json:"username"`
	Secret   string `json:"secret"`
}

type Game struct {
	Name string `json:"name"`
}

func GetAllPlayers() []Player {
	return []Player{}
}

func GetAllActivePlayers() []Player {
	return []Player{}
}

func (p *Player) CreateNewPlayer() (Player, error) {
	db := data.GetDB()
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(p.Secret), bcrypt.MinCost)
	if err != nil {
		return Player{}, err
	}

	query := "INSERT INTO players (username, secret) VALUES ($1, $2)"
	_, err = db.Exec(query, p.UserName, string(hashedSecret))
	if err != nil {
		return Player{}, err
	}
	return Player{UserName: p.UserName, Secret: p.Secret}, nil
}

func (g *CheckSecretRequest) CheckHashedSecret() (bool, error) {
	var hashedSecret string
	db := data.GetDB()
	query := "SELECT secret FROM players WHERE username=$1"
	err := db.QueryRow(query, g.UserName).Scan(&hashedSecret)
	if err != nil {
		return false, err
	}

	// check the hashed secret and return true or false
	err = bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(g.Secret))
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetAllActiveGames() []Game {
	return []Game{}
}

func GetAllGames() []Game {
	return []Game{}
}
