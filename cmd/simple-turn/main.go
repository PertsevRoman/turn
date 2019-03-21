package main

import (
	"fmt"
	"github.com/pions/turn"
	"github.com/xo/dburl"
	"log"
	"os"
	"plugin"
	"strconv"
)

func getEnvConf() (port int, realm string) {
	log.Printf("Use environment variables")

	realm = os.Getenv("REALM")
	if realm == "" {
		log.Panic("REALM is a required environment variable")
	}

	udpPortStr := os.Getenv("PORT")
	if udpPortStr == "" {
		log.Panic("PORT is a required environment variable")
	}

	udpPort, err := strconv.Atoi(udpPortStr)
	if err != nil {
		log.Panic(err)
	}

	return udpPort, realm
}

func loadTurnServer() turn.Server {
	pluginPath := "./plugins/env.so"

	// TODO remove dburl dependency
	dsn := os.Getenv("DB_DSN")
	if dsn != "" {
		dburl.Register(dburl.Scheme{
			Driver:    "redis",
			Generator: dburl.GenScheme("redis"),
			Proto:     0,
			Opaque:    false,
			Aliases:   []string{},
			Override:  "",
		})

		dsnMap, err := dburl.Parse(dsn)

		if err != nil {
			log.Panic("Cannot parse DB dsn: ", err)
		}

		pluginPath = fmt.Sprintf("./plugins/%s.so", dsnMap.Driver)
	}

	plug, err := plugin.Open(pluginPath)

	symTurnServer, err := plug.Lookup("TurnServer")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("TurnServer symbol loaded")
	}

	var turnServer turn.Server
	turnServer, ok := symTurnServer.(turn.Server)

	if !ok {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("TurnServer instance loaded")
	}

	return turnServer
}

func main() {
	turnServer := loadTurnServer()
	port, realm := getEnvConf()

	log.Printf("Starting on port %d", port)

	turnServer.Init()
	turn.Start(turn.StartArguments{
		Server:  turnServer,
		Realm:   realm,
		UDPPort: port,
	})
}
