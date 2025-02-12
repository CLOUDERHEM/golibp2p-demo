package main

import (
	"context"
	"fmt"
	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	ma "github.com/multiformats/go-multiaddr"
)

var (
	relayPeerIdString = "QmfNspDuxDv6bnFbzntc33YEgy3sGM9NTg9gCeMfUqcoeY"
	relayAddr         = "/ip4/127.0.0.1/udp/64324/quic-v1"
)

func main() {
	golog.SetAllLoggers(golog.LevelDebug)

	node, err := createHost()
	if err != nil {
		panic(err)
	}
	fmt.Println("my peerId is: ", node.ID())

	connectToRelay(context.Background(), node)

	node.SetStreamHandler("/test", func(stream network.Stream) {
		conn := stream.Conn()
		peerId := conn.RemotePeer()
		fmt.Printf("new peer connected, peerId: %v, addr: %v\n", peerId, conn.RemoteMultiaddr())
		bytes := make([]byte, 1024)
		for {
			n, err := stream.Read(bytes)
			if err != nil {
				return
			}
			fmt.Printf("read from: %v, data: %v\n", node.ID().String(), string(bytes[:n]))
		}
	})

	ch := make(chan struct{})
	ch <- struct{}{}
}

func createHost() (host.Host, error) {
	key, err := loadOrCreatePrivateKey("private_key.pem")
	if err != nil {
		return nil, err
	}

	node, err := libp2p.New(
		libp2p.Identity(key),
		libp2p.EnableRelay(),
		libp2p.NoListenAddrs,
	)
	if err != nil {
		return nil, err
	}

	connLogger := &ConnLogger{}
	node.Network().Notify(connLogger)

	return node, err
}

func connectToRelay(ctx context.Context, node host.Host) peer.AddrInfo {
	addr, err := ma.NewMultiaddr(relayAddr)
	if err != nil {
		panic(err)
	}
	peerId, err := peer.Decode(relayPeerIdString)
	if err != nil {
		panic(err)
	}
	addrInfo := peer.AddrInfo{
		ID:    peerId,
		Addrs: []ma.Multiaddr{addr},
	}
	err = node.Connect(ctx, addrInfo)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("connected to relay: ", relayPeerIdString)
	}
	_, err = client.Reserve(ctx, node, addrInfo)
	if err != nil {
		panic(err)
	}
	return addrInfo
}
