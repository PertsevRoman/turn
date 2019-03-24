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
	log.Printf("User auth: %s\n", username)

	port := m.dsn.Port

	addr := fmt.Sprintf("%s:%s", m.dsn.Host, port)

	dbPassword := m.dsn.Password

	db, err := strconv.Atoi(m.dsn.Db)

	if err != nil {
		log.Panic("Redis DB scheme not parsed")
	}

	log.Printf("Connecting to server...")
	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: dbPassword,
		DB:       db,
	})

	defer conn.Close()

	log.Printf("User search...\n")
	password, err = conn.Get(username).Result()

	if err == nil {
		log.Printf("User found: %s\n", username)
		return password, true
	}

	log.Printf("User not found: %s\n", username)
	return "", false
}

func (m *turnServer) Init() {
	dsn := os.Getenv("DB_DSN")
	m.dsn = GetDsnParts(dsn)

	log.Printf("Redis host: %s:%s", m.dsn.Host, m.dsn.Port)
	log.Printf("Redis DB: %s", m.dsn.Db)
}

var TurnServer turnServer
