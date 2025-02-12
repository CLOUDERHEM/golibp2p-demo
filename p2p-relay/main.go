package main

import (
	"fmt"
	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/multiformats/go-multiaddr"
)

var (
	relayAddr = "/ip4/0.0.0.0/udp/64324/quic-v1"
)

func main() {
	golog.SetAllLoggers(golog.LevelDebug)

	node, err := createHost()
	if err != nil {
		panic(err)
	}
	fmt.Println("my peerId is: ", node.ID())
	fmt.Println("listen addr: ", node.Addrs())

	_, err = relay.New(node)
	if err != nil {
		panic(err)
	}

	connLogger := &ConnLogger{}
	node.Network().Notify(connLogger)

	ch := make(chan struct{})
	ch <- struct{}{}
}

func createHost() (host.Host, error) {
	addr, err := multiaddr.NewMultiaddr(relayAddr)
	if err != nil {
		return nil, err
	}
	key, err := loadOrCreatePrivateKey("private_key.pem")
	if err != nil {
		return nil, err
	}
	return libp2p.New(
		libp2p.ListenAddrs(addr),
		libp2p.Identity(key),
	)
}
