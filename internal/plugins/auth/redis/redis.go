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
	parts = turn.MakeDsnParts(url)

	if parts.Port == "" {
		parts.Port = "6379"
	}

	if parts.Db == "" {
		parts.Db = "0"
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

func (m *turnServer) Init() {
	dsn := os.Getenv("DB_DSN")
	m.dsn = GetDsnParts(dsn)
}

var TurnServer turnServer
