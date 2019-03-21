package main

import (
	"github.com/pions/stun"
	"log"
	"os"
	"regexp"
)

type turnServer struct {
	usersMap map[string]string
}

func (m *turnServer) AuthenticateRequest(username string, srcAddr *stun.TransportAddr) (password string, ok bool) {
	if password, ok := m.usersMap[username]; ok {
		return password, true
	}

	return "", false
}

func (m *turnServer) PrintUsers() {
	for key, val := range m.usersMap {
		log.Println(key, val)
	}
}

func (m *turnServer) Init() {
	m.usersMap = make(map[string]string)

	usersString := os.Getenv("USERS")
	if usersString == "" {
		log.Panic("USERS is a required environment variable")
	}
	for _, kv := range regexp.MustCompile(`(\w+)=(\w+)`).FindAllStringSubmatch(usersString, -1) {
		m.usersMap[kv[1]] = kv[2]
	}
}

var TurnServer turnServer
