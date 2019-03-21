package main

import (
	"github.com/pions/stun"
	"net"
	"testing"
)

func TestTurnServer_Init(t *testing.T) {
	m := turnServer{}

	m.Init()

	if m.dsn == nil {
		t.Fail()
	}
}

func TestTurnServer_AuthenticateRequest(t *testing.T) {
	m := turnServer{}

	m.Init()

	ip := net.ParseIP("127.0.0.1")

	password, ok := m.AuthenticateRequest("tell", &stun.TransportAddr{
		Port: 3456,
		IP:   ip,
	})

	if !ok || password != "test" {
		t.Fail()
	}
}
