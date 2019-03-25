package main

import (
	"fmt"
	"github.com/pions/turn"
	"log"
	"os"
	"path/filepath"
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
	pluginPath := "plugins/env.so"

	// TODO remove dburl dependency
	dsn := os.Getenv("DB_DSN")

	if dsn != "" {
		parts := turn.MakeDsnParts(dsn)

		pluginPath = fmt.Sprintf("plugins/%s.so", parts.Proto)

	}

	pluginFullPath, err := filepath.Abs(pluginPath)
	log.Printf("Load module: %s", pluginFullPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		fmt.Println("Plugin path is not exists")
		os.Exit(1)
	}

	plug, err := plugin.Open(pluginFullPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
