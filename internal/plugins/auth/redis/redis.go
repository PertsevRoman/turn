package main

import (
	"github.com/pions/stun"
	"github.com/xo/dburl"
	"log"
	"os"
)

type turnServer struct {
	dsn *dburl.URL
}

func (m *turnServer) AuthenticateRequest(username string, srcAddr *stun.TransportAddr) (password string, ok bool) {
	return "", true
}

func (m *turnServer) PrintUsers() {
}

func (m *turnServer) Init() {
	dsn := os.Getenv("DB_DSN")
	if dsn != "" {
		dsnMap, err := dburl.Parse(dsn)

		if err != nil {
			log.Panic("Cannot parse DB dsn")
		}

		m.dsn = dsnMap
	}
}

var TurnServer turnServer
