package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/n7down/ssh-chess/internal/game"
	"github.com/n7down/ssh-chess/internal/logger"
	"github.com/n7down/ssh-chess/internal/logger/logruslogger"
	"golang.org/x/crypto/ssh"
)

func handler(conn net.Conn, gm *game.GameManager, config *ssh.ServerConfig, logger logger.Logger) {
	// Before use, a handshake must be performed on the incoming
	// net.Conn.
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	//_, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		logger.Debug("Failed to handshake with new client")
		return
	}

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple
		// terminal interface.
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()

		if err != nil {
			logger.Debug("could not accept channel.")
			return
		}

		// FIXME: pull out the requests
		// find the exec request and get the payload
		// put it back into the requests and pass it to the go func

		// TODO: Remove this -- only temporary while we launch on HN
		//
		// To see how many concurrent users are online
		//fmt.Printf("Player joined. Current stats: %d users, %d games\n",
		//gm.SessionCount(), gm.GameCount())

		// Reject all out of band requests accept for the unix defaults, pty-req and
		// shell.
		go func(in <-chan *ssh.Request) {
			for req := range in {

				// FIXME: is it possible to the the data from an exec request
				// ssh test@localhost ls - can i pull out the 'ls' and do something with the exec command?
				logger.Print(fmt.Sprintf("req: %v payload: %v", req.Type, string(req.Payload)))
				switch req.Type {
				case "pty-req":
					req.Reply(true, nil)
					continue
				case "shell":
					req.Reply(true, nil)
					continue
				}
				req.Reply(false, nil)
			}
		}(requests)

		//gm.HandleNewChannel(channel, sshConn.User())
		gm.HandleNewChannel(channel, sshConn.User())
	}
}

func main() {
	port := os.Getenv("PORT")

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		panic("Failed to load private key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Failed to parse private key")
	}

	config.AddHostKey(private)

	// logger
	logger := logruslogger.NewLogrusLogger(true)

	// create the GameManager
	gm := game.NewGameManager(logger)

	fmt.Printf("Listening on port %s for SSH...\n", port)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		panic("failed to listen for connection")
	}

	for {
		nConn, err := listener.Accept()
		if err != nil {
			panic("failed to accept incoming connection")
		}

		go handler(nConn, gm, config, logger)
	}
}
