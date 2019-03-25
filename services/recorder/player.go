package main

import (
	"errors"
	data "github.com/n7down/ssh-chess/services/db"
	"golang.org/x/crypto/bcrypt"
)

type AuthPlayerRequest struct {
	UserName string `json:"username"`
	Secret   string `json:"secret"`
}

type CreatePlayerRequest struct {
	UserName string `json:"username"`
	Secret   string `json:"secret"`
}

type CheckUserNameRequest struct {
	UserName string `json:"username"`
}

func (a *AuthPlayerRequest) CheckAuth() (bool, error) {
	var hashedSecret string
	db := data.GetDB()
	query := "SELECT secret FROM players WHERE username=$1"
	err := db.QueryRow(query, a.UserName).Scan(&hashedSecret)
	if err != nil {
		return false, err
	}

	// check the hashed secret and return true or false
	err = bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(a.Secret))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *AuthPlayerRequest) CheckUserNameExists() (bool, error) {
	var userNameExists *bool
	db := data.GetDB()
	query := "SELECT EXISTS(SELECT username FROM players WHERE username=$1)"
	row := db.QueryRow(query, a.UserName)
	err := row.Scan(&userNameExists)
	if err != nil {
		return false, err
	}

	if userNameExists == nil {
		return false, errors.New("username exists return nil")
	}

	return *userNameExists, nil
}

func (a *AuthPlayerRequest) CreateNewPlayer() error {
	db := data.GetDB()
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(a.Secret), bcrypt.MinCost)
	if err != nil {
		return err
	}

	query := "INSERT INTO players (username, secret) VALUES ($1, $2)"
	_, err = db.Exec(query, a.UserName, string(hashedSecret))
	if err != nil {
		return err
	}
	return nil
}

func (p *CreatePlayerRequest) CreateNewPlayer() (CreatePlayerRequest, error) {
	db := data.GetDB()
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(p.Secret), bcrypt.MinCost)
	if err != nil {
		return CreatePlayerRequest{}, err
	}

	query := "INSERT INTO players (username, secret) VALUES ($1, $2)"
	_, err = db.Exec(query, p.UserName, string(hashedSecret))
	if err != nil {
		return CreatePlayerRequest{}, err
	}
	return CreatePlayerRequest{UserName: p.UserName, Secret: p.Secret}, nil
}

func (c *CheckUserNameRequest) CheckUserNameExists() (bool, error) {
	var userNameExists *bool
	db := data.GetDB()
	query := "SELECT EXISTS(SELECT username FROM players WHERE username=$1)"
	row := db.QueryRow(query, c.UserName)
	err := row.Scan(&userNameExists)
	if err != nil {
		return false, err
	}

	if userNameExists == nil {
		return false, errors.New("username exists return nil")
	}

	return *userNameExists, nil
}
