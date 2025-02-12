package main

import (
	"log"

	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

type ConnLogger struct{}

func (cl *ConnLogger) Connected(net network.Network, conn network.Conn) {
	log.Printf("connected to: %s\n", conn.RemotePeer())
}

func (cl *ConnLogger) Disconnected(net network.Network, conn network.Conn) {
	log.Printf("disconnected from: %s\n", conn.RemotePeer())
}

func (cl *ConnLogger) Listen(net network.Network, addr ma.Multiaddr) {
	log.Printf("listening on: %s\n", addr)
}

func (cl *ConnLogger) ListenClose(net network.Network, addr ma.Multiaddr) {
	log.Printf("stopped listening on: %s\n", addr)
}
