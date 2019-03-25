package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade: ", err)
		return
	}
	defer c.Close()

	for {

		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println(fmt.Sprintf("Read: %s", err.Error()))
			return
		}
		fmt.Printf(fmt.Sprintf("Receive: %s", message))

		g := GameData{}

		err = json.Unmarshal([]byte(message), &g)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error unmarshalling request: %v", err.Error()))
			return
		}

		if g.EndTime == "" {
			fmt.Println(fmt.Sprintf("updating data: %v", g))
			err = g.Update()
			if err != nil {
				fmt.Println(fmt.Sprintf("Error updating the game data: %v", err.Error()))
				return
			}
		} else {
			fmt.Println(fmt.Sprintf("completing game: %v", g))
			err = g.Completed()
			if err != nil {
				fmt.Println(fmt.Sprintf("Error saving the game data: %v", err.Error()))
			}

		}

		//err = c.WriteMessage(mt, message)
		//if err != nil {
		//log.Println("Write:", err)
		//return
		//}
	}

}
