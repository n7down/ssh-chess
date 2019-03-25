package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/n7down/ssh-chess/common"
	"github.com/n7down/ssh-chess/logger"
)

var (
	gameDataQueue chan GameData
)

type GameData struct {
	Uuid            string `json:"uuid"`
	Name            string `json:"name"`
	BlackPlayerName string `json:"blackplayername"`
	WhitePlayerName string `json:"whiteplayername"`
	StartTime       string `json:"starttime"`
	EndTime         string `json:"endtime"`
	Outcome         string `json:"outcome"`
	Pgn             string `json:"pgn"`
}

type UserAuth struct {
	UserName string `json:"username"`
	Secret   string `json:"secret"`
}

type CheckUser struct {
	UserName string `json:"username"`
}

type CreateUser struct {
	UserName string `json:"username"`
	Secret   string `json:"secret"`
}

func WsSendData(uuid string, name string, blackPlayerName string, whitePlayerName, startTime string, outcome string, pgn string) error {
	if !GameSettings.DisableGameRecorder {
		var g GameData
		if outcome == "" {
			g = GameData{
				Uuid:            uuid,
				Name:            name,
				BlackPlayerName: blackPlayerName,
				WhitePlayerName: whitePlayerName,
				StartTime:       startTime,
				Pgn:             pgn,
			}
		} else {
			g = GameData{
				Uuid:            uuid,
				Name:            name,
				BlackPlayerName: blackPlayerName,
				WhitePlayerName: whitePlayerName,
				StartTime:       startTime,
				EndTime:         time.Now().String(),
				Outcome:         outcome,
				Pgn:             pgn,
			}
		}

		logger.Debug(fmt.Sprintf("sending game data: %v", g))
		gameDataQueue <- g
		logger.Debug(fmt.Sprintf("finished sending game data: %v", g))
	}
	return nil
}

func WsRun() {
	if !GameSettings.DisableGameRecorder {
		var mutex sync.RWMutex
		gameDataQueue = make(chan GameData)
		done := make(chan struct{})
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		gameDataQueue = make(chan GameData)

		recorderHost := common.GetEnv("RECORDER_HOST", "localhost")
		recorderHost = recorderHost + ":8000"
		logger.Debug(fmt.Sprintf("connect to websocket: %v", recorderHost))

		u := url.URL{
			Scheme: "ws",
			Host:   recorderHost,
			Path:   "/ws/record",
		}

		logger.Debug(fmt.Sprintf("connecting to %s", u.String()))

		connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		defer connection.Close()

		if err != nil {
			fmt.Println(fmt.Sprintf("error with recorder connection: %v", err.Error()))
			return
		}
		logger.Debug("connected to recorder: " + u.String())

		go func() {
			defer close(done)
			for {
				_, _, err := connection.ReadMessage()
				if err != nil {
					logger.Debug(fmt.Sprintf("read error: ", err.Error()))
					return
				}
			}
		}()

		for {
			select {
			case <-done:
				return
			case g := <-gameDataQueue:
				mutex.Lock()
				err := connection.WriteJSON(g)
				mutex.Unlock()

				if err != nil {
					logger.Debug(fmt.Sprintf("write error:", err.Error()))
					return
				}
			case <-interrupt:
				logger.Debug("interrupt")

				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					logger.Debug(fmt.Sprintf("write close error: ", err.Error()))
					return
				}
				select {
				case <-done:
				}
				return
			}
		}
	}
}

func SendData(uuid string, name string, blackPlayerName string, whitePlayerName, startTime string, outcome string, pgn string) error {
	if !GameSettings.DisableGameRecorder {
		recorderHost := common.GetEnv("RECORDER_HOST", "localhost")
		recorderPort := common.GetEnv("RECORDER_PORT", "8000")
		recorderUrl := "http://" + recorderHost + ":" + recorderPort

		const recorderUpdate = "/update"
		const recorderCompleted = "/completed"

		var completeRecorderUrl string
		var g GameData

		if outcome == "" {
			g = GameData{
				Uuid:            uuid,
				Name:            name,
				BlackPlayerName: blackPlayerName,
				WhitePlayerName: whitePlayerName,
				StartTime:       startTime,
				Pgn:             pgn,
			}

			completeRecorderUrl = recorderUrl + recorderUpdate

		} else {
			g = GameData{
				Uuid:            uuid,
				Name:            name,
				BlackPlayerName: blackPlayerName,
				WhitePlayerName: whitePlayerName,
				StartTime:       startTime,
				EndTime:         time.Now().String(),
				Outcome:         outcome,
				Pgn:             pgn,
			}

			completeRecorderUrl = recorderUrl + recorderCompleted
		}

		gameDataToSend, err := json.Marshal(g)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, completeRecorderUrl, bytes.NewBuffer(gameDataToSend))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		r, err := client.Do(req)
		if err != nil {
			return err
		}

		defer r.Body.Close()

		resp := struct {
			Message string `json:"message"`
		}{}

		err = json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			return err
		}
		logger.Debug(fmt.Sprintf("sending data to recorder: %v", g))
	}
	return nil
}

func (u *UserAuth) CheckSecret() (bool, error) {
	recorderHost := common.GetEnv("RECORDER_HOST", "localhost")
	recorderPort := common.GetEnv("RECORDER_PORT", "8000")
	recorderUrl := "http://" + recorderHost + ":" + recorderPort
	endPoint := "/api/player/auth"
	completeRecorderUrl := recorderUrl + endPoint

	d, err := json.Marshal(u)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodGet, completeRecorderUrl, bytes.NewBuffer(d))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	r, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer r.Body.Close()

	resp := struct {
		//Message string `json:"message"`
		Message bool `json:"message"`
	}{}

	logger.Debug(fmt.Sprintf("Response: %v", resp))
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		logger.Debug(fmt.Sprintf("error in decoder: %v", err.Error()))
		return false, err
	}
	logger.Debug(fmt.Sprintf("Check secret: %v", resp.Message))
	return resp.Message, nil
}

func (a *UserAuth) CheckUserExists() (bool, error) {
	recorderHost := common.GetEnv("RECORDER_HOST", "localhost")
	recorderPort := common.GetEnv("RECORDER_PORT", "8000")
	recorderUrl := "http://" + recorderHost + ":" + recorderPort
	endPoint := "/api/player/checkuser"
	completeRecorderUrl := recorderUrl + endPoint

	c := CheckUser{
		UserName: a.UserName,
	}

	d, err := json.Marshal(c)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodGet, completeRecorderUrl, bytes.NewBuffer(d))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	r, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer r.Body.Close()

	resp := struct {
		Message bool `json:"message"`
	}{}

	logger.Debug(fmt.Sprintf("Response: %v", resp))
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		logger.Debug(fmt.Sprintf("error in decoder: %v", err.Error()))
		return false, err
	}
	logger.Debug(fmt.Sprintf("Check user exists: %v", resp.Message))
	return resp.Message, nil
}

func (a *UserAuth) CreateNewUser() (string, error) {
	recorderHost := common.GetEnv("RECORDER_HOST", "localhost")
	recorderPort := common.GetEnv("RECORDER_PORT", "8000")
	recorderUrl := "http://" + recorderHost + ":" + recorderPort
	endPoint := "/api/player/create"
	completeRecorderUrl := recorderUrl + endPoint

	d, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, completeRecorderUrl, bytes.NewBuffer(d))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	r, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer r.Body.Close()

	resp := struct {
		Message string `json:"message"`
	}{}

	logger.Debug(fmt.Sprintf("Response: %v", resp))
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		logger.Debug(fmt.Sprintf("error in decoder: %v", err.Error()))
		return "", err
	}
	logger.Debug(fmt.Sprintf("Created user: %v", resp.Message))
	return resp.Message, nil
}
