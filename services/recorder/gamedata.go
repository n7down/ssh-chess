package main

import (
	"fmt"

	data "github.com/n7down/ssh-chess/services/db"
)

type GameData struct {
	Uuid            string `json:"uuid" binding:"required"`
	Name            string `json:"name" binding:"required"`
	BlackPlayerName string `json:"blackplayername" binding:"required"`
	WhitePlayerName string `json:"whiteplayername" binding:"required"`
	StartTime       string `json:"starttime" binding:"required"`
	EndTime         string `json:"endtime"`
	Outcome         string `json:"outcome"`
	Pgn             string `json:"pgn" binding:"required"`
}

func (g *GameData) Completed() error {
	db := data.GetDB()

	query := `UPDATE games SET end_time=?, outcome=?, pgn=? WHERE uuid=?`

	_, err := db.Exec(query, g.EndTime, g.Outcome, g.Pgn, g.Uuid)
	if err != nil {
		return err
	}
	return nil
}

func (g *GameData) Update() error {
	db := data.GetDB()

	fmt.Println(fmt.Sprintf("uuid: %s name: %s black: %s white: %s start time: %s pgn: %s",
		g.Uuid, g.Name, g.BlackPlayerName, g.WhitePlayerName, g.StartTime, g.Pgn))

	fmt.Println(fmt.Sprintf("data in gamedata: %v", g))

	query := fmt.Sprintf("INSERT INTO games (id, name, black_player_id, white_player_id, start_time, pgn) VALUES ('%s', '%s', (SELECT id FROM players WHERE username='%s'), (SELECT id FROM players WHERE username='%s'), '%s', '%s') ON CONFLICT (id) DO UPDATE SET pgn=excluded.pgn;", g.Uuid, g.Name, g.BlackPlayerName, g.WhitePlayerName, g.StartTime, g.Pgn)

	//query := fmt.Sprintf("INSERT INTO games (id, name, black_player_name, white_player_name, start_time, pgn) VALUES ('%s', '%s', '%s', '%s', '%s', '%s') ON CONFLICT (id) DO UPDATE SET pgn=excluded.pgn;", g.Uuid, g.Name, g.BlackPlayerName, g.WhitePlayerName, g.StartTime, g.Pgn)
	fmt.Println(fmt.Sprintf("query: %v", query))

	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
