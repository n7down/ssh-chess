package main

import (
	"fmt"
)

type Hub struct {
	Sessions   map[*Session]struct{}
	Redraw     chan struct{}
	Register   chan *Session
	Unregister chan UnregisterMessage
}

func NewHub() Hub {
	return Hub{
		Sessions:   make(map[*Session]struct{}),
		Redraw:     make(chan struct{}),
		Register:   make(chan *Session),
		Unregister: make(chan UnregisterMessage),
	}
}

func (h *Hub) Run(g *Game) {
	for {
		select {
		case <-h.Redraw:
			for s := range h.Sessions {
				go g.Render(s)
			}
		case s := <-h.Register:
			// Hide the cursor
			fmt.Fprint(s, "\033[?25l")

			h.Sessions[s] = struct{}{}
		case s := <-h.Unregister:
			if _, ok := h.Sessions[s.session]; ok {
				fmt.Fprint(s.session, s.message)

				// Unhide the cursor
				fmt.Fprint(s.session, "\033[?25h")

				delete(h.Sessions, s.session)
				s.session.c.Close()
			}
		}
	}
}
