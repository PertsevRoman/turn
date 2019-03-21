package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pions/stun"
	"github.com/xo/dburl"
	"log"
	"os"
	"strconv"
)

type turnServer struct {
	dsn *dburl.URL
}

func (m *turnServer) AuthenticateRequest(username string, srcAddr *stun.TransportAddr) (password string, ok bool) {
	port := m.dsn.Port()

	if port == "" {
		port = "6379"
	}

	addr := fmt.Sprintf("%s:%s", m.dsn.Host, port)

	db, err := strconv.Atoi(m.dsn.Path[1:])

	if err != nil {
		log.Panic("Redis DB scheme not parsed")
	}

	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})

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
	if dsn != "" {
		_, aliases := dburl.SchemeDriverAndAliases("redis")

		if aliases == nil {
			dburl.Register(dburl.Scheme{
				Driver:    "redis",
				Generator: dburl.GenScheme("redis"),
				Proto:     0,
				Opaque:    false,
				Aliases:   []string{},
				Override:  "",
			})
		}

		dsnMap, err := dburl.Parse(dsn)

		if err != nil {
			log.Panic("Cannot parse DB dsn: ", err)
		}

		m.dsn = dsnMap
	}
}

var TurnServer turnServer
