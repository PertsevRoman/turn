package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pions/stun"
	"github.com/pions/turn"
	"log"
	"os"
	"strconv"
)

type turnServer struct {
	dsn turn.DsnParts
}

func GetDsnParts(url string) (parts turn.DsnParts) {
	var port int
	var db string

	matches := turn.GetDnsMatches(url)

	proto := matches[0][1]

	if matches[0][5] == "" {
		port = 6379
	} else {
		port, _ = strconv.Atoi(matches[0][5])
	}

	if matches[0][6] == "" {
		db = "0"
	} else {
		db = matches[0][6]
	}

	parts = turn.DsnParts{
		Proto:    proto,
		Host:     matches[0][4],
		Username: matches[0][2],
		Password: matches[0][3],
		Port:     port,
		Db:       db,
	}

	return parts
}

func (m *turnServer) AuthenticateRequest(username string, srcAddr *stun.TransportAddr) (password string, ok bool) {
	port := m.dsn.Port

	addr := fmt.Sprintf("%s:%d", m.dsn.Host, port)

	dbPassword := m.dsn.Password

	db, err := strconv.Atoi(m.dsn.Db)

	if err != nil {
		log.Panic("Redis DB scheme not parsed")
	}

	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: dbPassword,
		DB:       db,
	})

	defer conn.Close()

	password, err = conn.Get(username).Result()

	if err == nil {
		return password, true
	}

	return "", false
}

func (m *turnServer) PrintUsers() {
}

func (m *turnServer) Init() {
	dsn := os.Getenv("DB_DSN")
	m.dsn = GetDsnParts(dsn)
}

var TurnServer turnServer
